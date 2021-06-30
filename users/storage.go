//go:generate mockgen -destination=storage_mock.go -package=users github.com/jenpaff/golang-microservices/users Storage

package users

import (
	"context"
	"database/sql"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/persistence/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Storage interface {
	Add(ctx context.Context, user common.User) error
	FindByName(ctx context.Context, userName string) (*common.User, error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{db: db}
}

func (p *storage) Add(ctx context.Context, user common.User) error {
	newUser := &models.User{
		Username:    user.UserName,
		Email:       null.StringFrom(user.Email),
		PhoneNumber: null.StringFrom(user.PhoneNumber),
	}
	return newUser.Insert(ctx, p.db, boil.Infer())
}

func (p *storage) FindByName(ctx context.Context, userName string) (*common.User, error) {
	returnedUser, err := models.Users(models.UserWhere.Username.EQ(userName)).One(ctx, p.db)
	if err != nil {
		return nil, err
	}
	return &common.User{
		UserName:    returnedUser.Username,
		PhoneNumber: returnedUser.PhoneNumber.String,
		Email:       returnedUser.Email.String,
	}, nil
}
