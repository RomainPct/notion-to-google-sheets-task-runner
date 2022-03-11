package database

type Automation struct {
	Id                   uint8
	Notion_database      string
	google_sheet         string
	google_sheet_tab     string
	notion_token         string
	google_refresh_token string
}
