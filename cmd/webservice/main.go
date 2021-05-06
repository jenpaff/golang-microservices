package main

import (
	"github.com/go-playground/log"
	. "github.com/jenpaff/golang-microservices/app"
)

func main() {
	log.Info("Starting app...")
	app := NewApp("12345")
	app.Start()
}
