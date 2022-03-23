package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/automationrunner"
	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/database"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvError := godotenv.Load("secret/.env")
	if loadEnvError != nil {
		log.Fatalf("Error loading .env file")
	}
	// Cron
	log.Println("--- NTGS task runner is running ---")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(10).Seconds().Do(cronTask)
	scheduler.StartBlocking()
}

func cronTask() {
	log.Println("-> Task time")
	automations, err := database.QueryWaitingAutomations()
	if err != nil {
		fmt.Println("Fail querying waiting automations : ", err.Error())
	}
	for _, automation := range automations {
		fmt.Println("-> Run automation ", automation.Id)
		go automationrunner.RunAutomation(automation, nil)
	}
}
