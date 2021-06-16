package users

import (
	"context"
	"fmt"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/persistence"
)

type Service struct {
	userPersistence persistence.Users
}

func (s Service) GetUser(ctx context.Context, userName string) (*persistence.User, error) {
	user, err := s.userPersistence.FindByName(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("could not find user with userName %s: %w", userName, api.ErrUserNotFound)
	}
	return user, nil
}
