package users

import (
	"context"
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
		return nil, err
	}
	return user, nil
}
