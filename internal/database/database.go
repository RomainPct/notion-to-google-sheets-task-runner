package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func query(request string) *sql.Rows {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_IDENTIFIER")+":"+os.Getenv("DATABASE_PASSWORD")+"@tcp("+os.Getenv("DATABASE_HOST")+":"+os.Getenv("DATABASE_PORT")+")/"+os.Getenv("DATABASE_NAME"))

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fetch, err := db.Query(request)

	if err != nil {
		panic(err.Error())
	}

	return fetch
}

func QueryAutomations() []Automation {
	fetch := query(`
				SELECT ntg_automations.id id, ntg_automations.notion_database, ntg_automations.google_sheet, ntg_automations.google_sheet_tab, ntg_users.notion_token, ntg_google_connections.google_refresh_token
				FROM ntg_automations
				INNER JOIN ntg_users ON ntg_users.id = ntg_automations.user_id
				INNER JOIN ntg_google_connections ON ntg_google_connections.id = ntg_automations.google_connection_id
				WHERE last_run IS NULL OR DATE_ADD(last_run, INTERVAL sync_recurrence MINUTE) <= NOW()
				ORDER BY last_run`)

	automations := []Automation{}

	for fetch.Next() {
		var automation Automation
		fetch.Scan(&automation.Id, &automation.Notion_database, &automation.google_sheet, &automation.google_sheet_tab, &automation.notion_token, &automation.google_refresh_token)
		automations = append(automations, automation)
	}

	defer fetch.Close()

	return automations
}
