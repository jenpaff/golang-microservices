package main

import (
	"github.com/go-playground/log"
)

func main() {
	log.Info("Starting app...")
	app := NewApp("12345")
	app.Start()
}
