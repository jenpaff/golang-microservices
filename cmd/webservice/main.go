package main

import (
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	. "github.com/jenpaff/golang-microservices/app"
)

func main() {
	initLogging()
	log.Info("Starting app...")
	app := NewApp("12345")
	app.Start()
}

func initLogging() {
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
}
