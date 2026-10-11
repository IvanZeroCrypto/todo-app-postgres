package users_service

import (
	"context"
	"fmt"

	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("id cannot be less than zero or negative:%w", core_errors.ErrInvalidArgument)
	}
	if err := s.usersRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user:%w", err)
	}
	return nil
}
