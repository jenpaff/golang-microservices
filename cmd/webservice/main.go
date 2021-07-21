package main

import (
	"fmt"
	"github.com/go-playground/errors"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/jenpaff/golang-microservices/app"
	"os"
	"os/signal"
	"syscall"
)

const (
	injectedSecretsEnv = "ENV_SECRETS"
)

func main() {
	initLogging()
	log.Info("Starting app...")

	configPath, secretsPath, secretsEnv, err := getStartArgs()
	if err != nil {
		logErrorAndExit(fmt.Errorf("could not fetch necessary paths: %s", err))
	}

	app, err := app.NewApp("12345", configPath, secretsPath, secretsEnv)
	if err != nil {
		logErrorAndExit(fmt.Errorf("could not initialise app: %s", err.Error()))
	}

	err = app.Start()
	if err != nil {
		logErrorAndExit(err)
	}

	// essential to make the main wait for a stop signal
	// otherwise we exit the main function and also the go-routine in app.Start() will be stopped
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	app.Stop()
}

func initLogging() {
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
}

func logErrorAndExit(err error) {
	cLog := console.New(true)

	log.AddHandler(cLog, log.AllLevels...)
	log.Fatalf(err.Error())
	os.Exit(1)
}

func getStartArgs() (string, string, string, error) {
	if len(os.Args) < 2 {
		return "", "", "", errors.New("you must supply a config file path")
	}

	if len(os.Args) < 3 {
		return "", "", "", errors.New("you must supply a secrets directory path")
	}

	envConfigFilePath := os.Args[1]
	secretsDirectoryPath := os.Args[2]
	injectedSecrets := os.Getenv(injectedSecretsEnv)
	return envConfigFilePath, secretsDirectoryPath, injectedSecrets, nil
}
