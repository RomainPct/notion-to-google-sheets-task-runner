package database

import (
	"database/sql"
	"os"
	"strconv"
)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_IDENTIFIER")+":"+os.Getenv("DATABASE_PASSWORD")+"@tcp("+os.Getenv("DATABASE_HOST")+":"+os.Getenv("DATABASE_PORT")+")/"+os.Getenv("DATABASE_NAME"))
	return db, err
}

func exec(request string, args ...interface{}) (sql.Result, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	exec, err := db.Exec(request, args...)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return exec, nil
}

func query(request string, args ...interface{}) (*sql.Rows, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	fetch, err := db.Query(request, args...)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return fetch, nil
}

func SetAutomationExecution(automation Automation, success bool, errorLabel string) (int64, error) {
	exec, err := exec(`
				INSERT INTO ntg_automation_executions (automation_id, success, error_label)
				VALUES (?, ?, ?)
				`, automation.Id, success, errorLabel)
	if err != nil {
		return 0, err
	}
	lid, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lid, nil
}

func SetAutomationLastRun(automation Automation) error {
	fetch, err := query(`
				UPDATE ntg_automations
				SET last_run = NOW()
				WHERE id = ?
				`, automation.Id)
	defer fetch.Close()
	return err
}

func QueryWaitingAutomations() ([]Automation, error) {
	fetch, err := query(`
				SELECT ntg_automations.id id, ntg_automations.notion_database, ntg_automations.google_sheet, ntg_automations.google_sheet_tab, ntg_users.notion_token, ntg_google_connections.google_refresh_token
				FROM ntg_automations
				INNER JOIN ntg_users ON ntg_users.id = ntg_automations.user_id
				INNER JOIN ntg_google_connections ON ntg_google_connections.id = ntg_automations.google_connection_id
				WHERE last_run IS NULL OR DATE_ADD(last_run, INTERVAL sync_recurrence MINUTE) <= NOW()
				ORDER BY last_run
				`)
	if err != nil {
		return nil, err
	}
	automations := []Automation{}
	for fetch.Next() {
		var automation Automation
		fetch.Scan(&automation.Id, &automation.Notion_database, &automation.Google_sheet, &automation.Google_sheet_tab, &automation.Notion_token, &automation.Google_refresh_token)
		automations = append(automations, automation)
	}
	defer fetch.Close()
	return automations, err
}

func QueryAutomationWithID(id string) (Automation, error) {
	idInt, _ := strconv.ParseInt(id, 10, 64)
	automation := Automation{
		Id: uint64(idInt),
	}
	fetch, err := query(`
				SELECT ntg_automations.id id, ntg_automations.notion_database, ntg_automations.google_sheet, ntg_automations.google_sheet_tab, ntg_users.notion_token, ntg_google_connections.google_refresh_token
				FROM ntg_automations
				INNER JOIN ntg_users ON ntg_users.id = ntg_automations.user_id
				INNER JOIN ntg_google_connections ON ntg_google_connections.id = ntg_automations.google_connection_id
				WHERE ntg_automations.id = ?
				LIMIT 1
				`, id)
	if err != nil {
		return automation, err
	}
	fetch.Next()
	err = fetch.Scan(&automation.Id, &automation.Notion_database, &automation.Google_sheet, &automation.Google_sheet_tab, &automation.Notion_token, &automation.Google_refresh_token)
	if err != nil {
		return automation, err
	}
	defer fetch.Close()
	return automation, nil
}
