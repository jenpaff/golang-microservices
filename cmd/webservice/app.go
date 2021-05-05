package main

import (
	"context"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	server *http.Server
	port   string
}

func NewApp(port string) *App {
	controller := api.NewController()
	router := api.NewRouter(controller)
	server := &http.Server{Addr: ":" + port, Handler: router}
	return &App{server: server, port: port}
}

func (a *App) Start() {
	var err error
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Errorf("Could not listen on port %s", a.port)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("Shutting down server")

	if serverErr := a.server.Shutdown(context.Background()); serverErr != nil {
		log.WithError(err).Warn("Couldn't shutdown server")
		os.Exit(1)
	}

	log.Info("Shut down successful")
}
