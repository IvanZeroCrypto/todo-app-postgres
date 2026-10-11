package users_service

import (
	"context"
	"fmt"

	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
)

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain:%w", err)
	}

	usernew, err := s.usersRepository.UserCreate(ctx, user)

	if err != nil {
		return domain.User{}, fmt.Errorf("create user:%w", err)
	}
	return usernew, nil
}
