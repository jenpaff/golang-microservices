package persistence

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jenpaff/golang-microservices/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres(config config.PersistenceConfig) (*sql.DB, error) {
	pgOptions := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", config.DbHost, config.DbPort, config.DbUsername, config.DbName, config.DbPassword)
	if !config.SslEnabled {
		pgOptions = pgOptions + " sslmode=disable"
	}

	db, err := sql.Open("postgres", pgOptions)
	if err != nil {
		return nil, err
	}

	log.Infof("PostgreSQL storage: connected to host %s:%d database %s with user %s", config.DbHost, config.DbPort, config.DbName, config.DbUsername)
	return db, nil
}
