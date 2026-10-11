package users_service

import (
	"context"
	"fmt"

	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

func (s *UserService) GetUser(ctx context.Context, id int) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, fmt.Errorf("id cannot be negative and is equal to zero:%w", core_errors.ErrInvalidArgument)
	}
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository:%w", err)
	}
	return user, nil
}
