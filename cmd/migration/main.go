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

	dbConfig := config.PersistenceConfig{
		DbHost:     os.Args[1],
		DbPort:     int(port),
		DbName:     os.Args[3],
		SslEnabled: os.Args[4] == "true",
		DbUsername: os.Args[5],
		DbPassword: os.Args[6],
	}

	pgOptions := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbUsername, dbConfig.DbName, dbConfig.DbPassword)
	if !dbConfig.SslEnabled {
		pgOptions = pgOptions + " sslmode=disable"
	}

	db, err := sql.Open("postgres", pgOptions)
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}

	log.Infof("PostgreSQL storage: connected to host %s:%d database %s with user %s", dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbName, dbConfig.DbUsername)

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: dbConfig.DbName,
	})
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbConfig.DbName, driver)
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
