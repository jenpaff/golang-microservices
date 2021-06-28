package main

import (
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	. "github.com/jenpaff/golang-microservices/app"
	"os"
)

func main() {
	initLogging()
	log.Info("Starting app...")
	app := NewApp("12345")
	err := app.Start()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}

func initLogging() {
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
}
