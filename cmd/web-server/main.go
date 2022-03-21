package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/automationrunner"
	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvError := godotenv.Load("secret/.env")
	if loadEnvError != nil {
		log.Fatalf("Error loading .env file")
	}
	// Web server
	log.Println("--- NTGS web server is running ---")
	http.HandleFunc("/trigger", trigger)
	http.ListenAndServe(":80", nil)
}

func trigger(w http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["id"]
	if !ok || len(keys) < 1 {
		return
	}
	id := keys[0]
	automation, err := database.QueryAutomationWithID(id)
	if err != nil {
		automationrunner.SaveResult(automation, err, "trigger_querying")
		fmt.Println("Fail querying automation : ", err.Error())
		return
	}
	automationrunner.RunAutomation(automation)
	fmt.Println(id)
}
