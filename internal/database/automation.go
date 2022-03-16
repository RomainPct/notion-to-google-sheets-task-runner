package database

type Automation struct {
	Id                   uint64
	Notion_database      string
	Google_sheet         string
	Google_sheet_tab     string
	Notion_token         string
	Google_refresh_token string
}
