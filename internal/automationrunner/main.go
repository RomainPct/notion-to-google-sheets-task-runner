package automationrunner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/database"
	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/dataformatter"
	"github.com/jomei/notionapi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func RunAutomation(automation database.Automation) bool {
	// Set automation last run date
	database.SetAutomationLastRun(automation)
	// Set notion and gsheets api
	notion := notionapi.NewClient(notionapi.Token(automation.Notion_token))
	notionDatabaseId := notionapi.DatabaseID(automation.Notion_database)
	sheetsService, err := getSheetService(automation.Google_refresh_token)
	if err != nil {
		return SaveResult(automation, err, "google_configuration")
	}
	// Create gsheet tab if needed
	spreadhsheet, err := sheetsService.Spreadsheets.Get(automation.Google_sheet).Do()
	if err != nil {
		return SaveResult(automation, err, "google_read_sheet")
	}
	tabExists := dataformatter.TabExists(automation.Google_sheet_tab, spreadhsheet.Sheets)
	if !tabExists {
		_, err := sheetsService.Spreadsheets.BatchUpdate(automation.Google_sheet, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{Title: automation.Google_sheet_tab},
				},
			}},
		}).Do()
		if err != nil {
			return SaveResult(automation, err, "google_add_tab")
		}
	}
	// Get notion fields
	database, err := notion.Database.Get(context.Background(), notionDatabaseId)
	if err != nil {
		return SaveResult(automation, err, "notion_read_fields")
	}
	properties, fields := dataformatter.GenerateNotionFields(database.Properties)
	notionHeaders := append([]string{"id", "created_time", "last_edited_time"}, fields...)
	// Get gsheet headers
	headers, err := sheetsService.Spreadsheets.Values.Get(automation.Google_sheet, automation.Google_sheet_tab+"!A1:ZZ1").Do()
	if err != nil {
		return SaveResult(automation, err, "google_read_headers")
	}
	var needRebuild bool
	if len(headers.Values) > 0 {
		needRebuild = !dataformatter.Equal(headers.Values[0], notionHeaders)
	} else {
		needRebuild = true
	}
	// Get and format notion data
	notionRows, err := getNotionData(notion, notionDatabaseId, properties, needRebuild)
	if err != nil {
		return SaveResult(automation, err, "notion_get_data")
	}
	// Read existing ids in sheet
	existingIds, err := sheetsService.Spreadsheets.Values.Get(automation.Google_sheet, automation.Google_sheet_tab+"!A1:A").Do()
	if err != nil {
		return SaveResult(automation, err, "google_read_ids")
	}
	// Organize between new data and data to update
	existingRows, newRows := dataformatter.SplitNotionData(notionRows, existingIds.Values)
	// Update headers
	if needRebuild {
		_, err = sheetsService.Spreadsheets.Values.Update(
			automation.Google_sheet,
			automation.Google_sheet_tab+"!A1:ZZ1",
			&sheets.ValueRange{Values: [][]interface{}{dataformatter.FormatToRowValueRange(notionHeaders)}},
		).ValueInputOption("RAW").Do()
		if err != nil {
			return SaveResult(automation, err, "google_update_headers")
		}
	}
	// Add new data
	newRowsCount := len(existingIds.Values)
	if newRowsCount == 0 {
		newRowsCount = 1
	}
	_, err = sheetsService.Spreadsheets.Values.Update(
		automation.Google_sheet,
		automation.Google_sheet_tab+"!A"+strconv.Itoa(1+newRowsCount)+":ZZ"+strconv.Itoa(1+newRowsCount+len(newRows)),
		&sheets.ValueRange{Values: dataformatter.FormatToValueRange(newRows)},
	).ValueInputOption("RAW").Do()
	if err != nil {
		return SaveResult(automation, err, "google_insert_data")
	}
	// Update existing data
	_, err = sheetsService.Spreadsheets.Values.BatchUpdate(automation.Google_sheet, &sheets.BatchUpdateValuesRequest{
		Data:             dataformatter.FormatValueRangeBatch(automation.Google_sheet_tab, existingRows),
		ValueInputOption: "RAW",
	}).Do()
	if err != nil {
		return SaveResult(automation, err, "google_update_data")
	}
	return SaveResult(automation, nil)
}

func SaveResult(automation database.Automation, automationErr error, label ...string) bool {
	result := automationErr == nil
	executionId, err := database.SetAutomationExecution(automation, result, label...)
	if err != nil {
		panic(err.Error())
	}
	stringExecutionId := strconv.Itoa(int(executionId))
	if automationErr != nil {
		os.WriteFile("./logs/error-"+stringExecutionId+".txt", []byte(automationErr.Error()), 0444)
		fmt.Println("Automation ", automation.Id, " did fail (Check error-"+stringExecutionId+".txt for more details)")
	} else {
		fmt.Println("Automation ", automation.Id, " did run successfully")
	}
	return result
}

func getNotionData(notion *notionapi.Client, id notionapi.DatabaseID, properties []string, rebuild bool) ([][]string, error) {
	data := []notionapi.Page{}
	var startCursor *notionapi.Cursor = nil
	for {
		request := notionapi.DatabaseQueryRequest{
			PageSize: 100,
			Sorts:    []notionapi.SortObject{{Timestamp: "last_edited_time", Direction: "descending"}},
		}
		if startCursor != nil {
			request.StartCursor = notionapi.Cursor(*startCursor)
		}
		req, err := notion.Database.Query(context.Background(), id, &request)
		if err != nil {
			return nil, err
		}
		data = append(data, req.Results...)
		if req.HasMore {
			startCursor = &req.NextCursor
		} else {
			startCursor = nil
		}
		if !rebuild || startCursor == nil {
			break
		}
	}
	rows := make([][]string, len(data))
	for index, row := range data {
		columns := []string{
			row.ID.String(),
			row.CreatedTime.String(),
			row.LastEditedTime.String(),
		}
		for _, property := range properties {
			val := row.Properties[property]
			if val != nil {
				columns = append(columns, dataformatter.ReadNotionPropertyValue(val)...)
			} else {
				columns = append(columns, "")
			}
		}
		rows[index] = columns
	}
	return rows, nil
}

func getSheetService(refreshToken string) (*sheets.Service, error) {
	googleCtx := context.Background()
	b, err := ioutil.ReadFile("secret/credentials.json")
	if err != nil {
		return nil, err
	}
	googleConfig, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/drive.metadata.readonly")
	if err != nil {
		return nil, err
	}
	googleToken := oauth2.Token{RefreshToken: refreshToken}
	sheetsService, err := sheets.NewService(googleCtx, option.WithTokenSource(googleConfig.TokenSource(googleCtx, &googleToken)))
	if err != nil {
		return nil, err
	}
	return sheetsService, nil
}
