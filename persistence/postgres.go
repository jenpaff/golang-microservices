package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/persistence/models"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Add(ctx context.Context, user User) error {
	newUser := &models.User{
		Username:    user.UserName,
		Email:       null.StringFrom(user.Email),
		PhoneNumber: null.StringFrom(user.Telephone),
	}
	return newUser.Insert(ctx, p.db, boil.Infer())
}

func (p *Postgres) FindByName(ctx context.Context, userName string) (*User, error) {
	returnedUser, err := models.Users(models.UserWhere.Username.EQ(userName)).One(ctx, p.db)
	if err != nil {
		return nil, err
	}
	return &User{
		UserName:  returnedUser.Username,
		Telephone: returnedUser.PhoneNumber.String,
		Email:     returnedUser.Email.String,
	}, nil
}
