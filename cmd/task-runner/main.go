package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	log.Println("Welcome")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Seconds().Do(cronTask)
	scheduler.StartBlocking()
	log.Println("main did end")
}

func cronTask() {
	log.Println("-> next")
	database.queryAutomations()
}
