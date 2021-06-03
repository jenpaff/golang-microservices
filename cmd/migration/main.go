package main

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jenpaff/golang-microservices/config"
	"os"
	"strconv"
)

func main() {
	initLogging()

	if len(os.Args) != 7 {
		log.Fatalf("the migrations command expects 6 command line parameters: dbHost dbPort dbName sslEnabled dbUsername dbPassword")
	}

	port, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("the migrations command expects the port to be a number: %s", os.Args[2])
	}

	pgConfig := config.Postgres{
		Host:       os.Args[1],
		Port:       int(port),
		DBName:     os.Args[3],
		SSLEnabled: os.Args[4] == "true",
		UserName:   os.Args[5],
		Password:   os.Args[6],
	}
	log.Infof("config is %v", pgConfig)

	pgOptions := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", pgConfig.Host, pgConfig.Port, pgConfig.UserName, pgConfig.DBName, pgConfig.Password)
	if !pgConfig.SSLEnabled {
		pgOptions = pgOptions + " sslmode=disable"
	}

	db, err := sql.Open("postgres", pgOptions)
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}

	log.Infof("PostgreSQL storage: connected to host %s:%d database %s with user %s", pgConfig.Host, pgConfig.Port, pgConfig.DBName, pgConfig.UserName)

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: pgConfig.DBName,
	})
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		pgConfig.DBName, driver)
	if err != nil {
		log.Fatalf("could not init DB migrations: %s", err.Error())
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("could not migrate the DB: %s", err.Error())
		} else {
			log.Info("No database changes necessary")
		}
	}
	log.Info("PostgreSQL storage: migrations finished")
}

func initLogging() {
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
}
