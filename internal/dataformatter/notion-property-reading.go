package dataformatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jomei/notionapi"
)

func readNotionPropertyStructure(key string, prop notionapi.PropertyConfig) ([]string, bool) {
	switch prop.GetType() {
	case "date":
		return []string{"start", "end"}, true
	case "select":
		return []string{"optionName", "optionId"}, true
	case "multi_select":
		return []string{"optionNames", "optionIds"}, true
	case "people":
		return []string{"names", "ids", "emails"}, true
	case "person":
		return []string{"name", "id", "email"}, true
	case "bot":
		return []string{"name", "id"}, true
	}
	return []string{}, false
}

func ReadNotionPropertyValue(property notionapi.Property) []string {
	switch p := property.(type) {
	default:
		fmt.Printf("unexpected type %T\n", p) // %T prints whatever type t has
		return []string{""}
	case *notionapi.TitleProperty:
		return []string{p.Title[0].PlainText}
	case *notionapi.RichTextProperty:
		return []string{p.RichText[0].PlainText}
	case *notionapi.CheckboxProperty:
		return []string{strconv.FormatBool(p.Checkbox)}
	case *notionapi.NumberProperty:
		return []string{strconv.FormatFloat(p.Number, 'f', -1, 64)}
	case *notionapi.PhoneNumberProperty:
		return []string{p.PhoneNumber}
	case *notionapi.EmailProperty:
		return []string{p.Email}
	case *notionapi.TextProperty:
		return []string{p.Text[0].PlainText}
	case *notionapi.CreatedTimeProperty:
		return []string{p.CreatedTime.String()}
	case *notionapi.LastEditedTimeProperty:
		return []string{p.LastEditedTime.String()}
	case *notionapi.DateProperty:
		return readDate(p.Date)
	case *notionapi.SelectProperty:
		return []string{p.Select.Name, p.Select.ID.String()}
	case *notionapi.MultiSelectProperty:
		return readOptions(p.MultiSelect)
	case *notionapi.PeopleProperty:
		return readUsers(p.People)
	case *notionapi.URLProperty:
		return []string{p.URL}
	case *notionapi.CreatedByProperty:
		return readUser(p.CreatedBy)
	case *notionapi.LastEditedByProperty:
		return readUser(p.LastEditedBy)
	case *notionapi.FormulaProperty:
		switch p.Formula.Type {
		case notionapi.FormulaTypeBoolean:
			return []string{strconv.FormatBool(p.Formula.Boolean)}
		case notionapi.FormulaTypeNumber:
			return []string{strconv.FormatFloat(p.Formula.Number, 'f', -1, 64)}
		case notionapi.FormulaTypeString:
			return []string{p.Formula.String}
		case notionapi.FormulaTypeDate:
			return []string{strings.Join(readDate(*p.Formula.Date), "-")}
		}
	case *notionapi.RelationProperty:
		return []string{readRelations(p.Relation)}
	case *notionapi.RollupProperty:
		switch p.Rollup.Type {
		case notionapi.RollupTypeNumber:
			return []string{strconv.FormatFloat(p.Rollup.Number, 'f', -1, 64)}
		case notionapi.RollupTypeDate:
			return []string{strings.Join(readDate(*p.Rollup.Date), "-")}
		case notionapi.RollupTypeArray:
			return []string{readRollupArray(p.Rollup.Array)}
		}
	case *notionapi.FilesProperty:
		return []string{readFiles(p.Files)}
	}
	return []string{""}
}

func readRollupArray(array notionapi.PropertyArray) string {
	data := []string{}
	for _, property := range array {
		data = append(data, strings.Join(ReadNotionPropertyValue(property), "||"))
	}
	return strings.Join(data, "&&")
}

func readFiles(files []notionapi.File) string {
	data := []string{}
	for _, file := range files {
		data = append(data, file.Name)
	}
	return strings.Join(data, "&&")
}

func readDate(date notionapi.DateObject) []string {
	//OPTIMIZE: Better text description
	dates := []string{}
	if date.Start != nil {
		dates = append(dates, date.Start.String())
	} else {
		dates = append(dates, "")
	}
	if date.End != nil {
		dates = append(dates, date.End.String())
	} else {
		dates = append(dates, "")
	}
	return dates
}

func readRelations(relations []notionapi.Relation) string {
	data := []string{}
	for _, relation := range relations {
		data = append(data, relation.ID.String())
	}
	return strings.Join(data, "&&")
}

func readOptions(options []notionapi.Option) []string {
	optionsData := [][]string{{}, {}}
	for _, option := range options {
		optionsData[0] = append(optionsData[0], option.Name)
		optionsData[1] = append(optionsData[1], option.ID.String())
	}
	return []string{
		strings.Join(optionsData[0], ","),
		strings.Join(optionsData[1], "&&"),
	}
}

func readUsers(users []notionapi.User) []string {
	usersData := [][]string{{}, {}, {}}
	for _, user := range users {
		data := readUser(user)
		usersData[0] = append(usersData[0], data[0])
		usersData[1] = append(usersData[1], data[1])
		usersData[2] = append(usersData[2], data[2])
	}
	return []string{
		strings.Join(usersData[0], ","),
		strings.Join(usersData[1], "&&"),
		strings.Join(usersData[2], ","),
	}
}

func readUser(user notionapi.User) []string {
	return []string{user.Name, user.ID.String(), user.Person.Email}
}
