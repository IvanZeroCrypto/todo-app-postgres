package users_postgres_repository

import (
	"context"
	"fmt"

	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	SELECT id,version,full_name,phone_number FROM todoapp.users
	ORDER BY id ASC 
	LIMIT $1
	OFFSET $2;
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users:%w", err)
	}
	defer rows.Close()

	var usersModel []UserModel

	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users:%w", err)

		}
		usersModel = append(usersModel, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows getusers:%w", err)
	}
	usersDomains := userDomainsFromModels(usersModel)

	return usersDomains, nil
}
