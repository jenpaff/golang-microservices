package containertests

import (
	"database/sql"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/persistence"
)

func createConnection() (db *sql.DB, err error) {
	cfg, err := config.BuildConfig("../config/test.json", "", "")
	db, _ = persistence.ConnectPostgres(cfg.Persistence)
	if err != nil {
		return nil, err
	}

	// Activate in case you want to debug the SQL statements
	// db.LogMode(true)
	return db, nil
}

func CleanUpDatabase() error {
	db, err := createConnection()
	if err != nil {
		return err
	}

	_, err = db.Exec("delete from users")
	if err != nil {
		return err
	}

	return nil
}
