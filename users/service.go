//go:generate mockgen -destination=service_mock.go -package=users github.com/jenpaff/golang-microservices/users Service

package users

import (
	"context"
	"github.com/jenpaff/golang-microservices/common"
)

type Service interface {
	GetUser(ctx context.Context, userName string) (*common.User, error)
	CreateUser(ctx context.Context, userName, email, phoneNumber string) (*common.User, error)
	CreateUserWithNewFeature(ctx context.Context, userName, email, phoneNumber string) (*common.User, error)
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

func (s service) CreateUser(ctx context.Context, userName, email, phoneNumber string) (*common.User, error) {
	// TODO: set uuid here
	user, err := s.storage.Create(ctx, userName, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) CreateUserWithNewFeature(ctx context.Context, userName, email, phoneNumber string) (*common.User, error) {
	// TODO: set uuid here
	user, err := s.storage.Create(ctx, userName, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	return user, nil
}
