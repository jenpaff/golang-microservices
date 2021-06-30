package users

import (
	"context"
	"fmt"
	"github.com/jenpaff/golang-microservices/common"
)

type Service interface {
	GetUser(ctx context.Context, userName string) (*common.User, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return service{storage: storage}
}

func (s service) GetUser(ctx context.Context, userName string) (*common.User, error) {
	user, err := s.storage.FindByName(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("could not find user with userName %s: %w", userName, err)
	}
	return user, nil
}
