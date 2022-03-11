package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RomainPct/notion-to-google-sheets-task-runner/internal/database"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvError := godotenv.Load(".env")
	if loadEnvError != nil {
		log.Fatalf("Error loading .env file")
	}
	log.Println("Welcome")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Seconds().Do(cronTask)
	scheduler.StartBlocking()
	log.Println("main did end")
}

func cronTask() {
	log.Println("-> next")
	automations := database.QueryAutomations()
	for _, automation := range automations {
		fmt.Println(automation.Id)
		go runAutomation(automation)
	}
}

func runAutomation(automation database.Automation) {
	// time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Println(automation)
}
