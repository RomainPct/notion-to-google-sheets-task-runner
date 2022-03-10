package database

import (
	"database/sql"
	"fmt"
)

func query(request string) {

}

func queryAutomations() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/notion-to-google-sheets")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fetch, err := db.Query("SELECT id, notion_database FROM ntg_automations")

	if err != nil {
		panic(err.Error())
	}

	automations := []Automation{}

	for fetch.Next() {
		var automation Automation
		fetch.Scan(&automation.id, &automation.notion_database)
		automations = append(automations, automation)
	}
	fmt.Println(automations)

	defer fetch.Close()
}
