package persistence

import (
	"context"
)

type Users interface {
	Add(ctx context.Context, user User) error
	FindByName(ctx context.Context, userName string) (*User, error)
}

type User struct {
	UserName  string
	Telephone string
	Email     string
}
