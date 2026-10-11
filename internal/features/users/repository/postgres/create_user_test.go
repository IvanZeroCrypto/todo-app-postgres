package users_postgres_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
)

func TestRepositoryUserCreate(t *testing.T) {
	tests := []struct {
		name          string
		userInput     domain.User
		checkResponse func(*testing.T, domain.User, *TestPool)
		setup         func(*testing.T, *UsersRepository)
		wantError     bool
		errorMsg      string
		cancelContext bool
	}{
		{
			name: "Success create user in repository with FullName and PhoneNumber ",
			userInput: domain.User{
				FullName:    fullName,
				PhoneNumber: &phoneNumber,
			},
			checkResponse: func(t *testing.T, u domain.User, pool *TestPool) {
				assert.NotZero(t, u.ID)
				assert.Equal(t, 1, u.Version)
				assert.Equal(t, fullName, u.FullName)
				assert.Equal(t, phoneNumber, *u.PhoneNumber)

				var count int
				query := `SELECT COUNT(*) FROM todoapp.users WHERE phone_number=$1`
				errQuery := pool.QueryRow(context.Background(), query, phoneNumber).Scan(&count)
				require.NoError(t, errQuery)
				assert.Equal(t, 1, count)
			},
			wantError: false,
		},
		{
			name: "Success create user in repository only FullName ",
			userInput: domain.User{
				FullName: fullName,
			},
			checkResponse: func(t *testing.T, u domain.User, pool *TestPool) {
				assert.NotZero(t, u.ID)
				assert.Equal(t, 1, u.Version)
				assert.Equal(t, fullName, u.FullName)
				assert.Nil(t, u.PhoneNumber)

				var count int
				query := `SELECT COUNT(*) FROM todoapp.users WHERE full_name=$1`
				errQuery := pool.QueryRow(context.Background(), query, fullName).Scan(&count)
				require.NoError(t, errQuery)
				assert.Equal(t, 1, count)
			},
			wantError: false,
		},
		{
			name: "Fail create user in repository duplicate PhoneNumber",
			userInput: domain.User{
				FullName:    fullName,
				PhoneNumber: &phoneNumber,
			},
			setup: func(t *testing.T, repo *UsersRepository) {
				existingUser := domain.User{
					FullName:    fullName,
					PhoneNumber: &phoneNumber,
				}
				_, err := repo.UserCreate(context.Background(), existingUser)
				require.NoError(t, err)
			},
			wantError: true,
			errorMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "Fail create user in repository context cancelled",
			userInput: domain.User{
				FullName:    fullName,
				PhoneNumber: &phoneNumber,
			},
			wantError:     true,
			cancelContext: true,
			errorMsg:      "scan error:context canceled",
		},
		{
			name: "Fail create user in repository FullName null",
			userInput: domain.User{
				FullName:    emptyFullName,
				PhoneNumber: &phoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
		{
			name: "Fail create user in repository FullName long",
			userInput: domain.User{
				FullName:    longFullName,
				PhoneNumber: &phoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
		{
			name: "Fail create user in repository FullName short",
			userInput: domain.User{
				FullName:    shortFullName,
				PhoneNumber: &phoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
		{
			name: "Fail create user in repository PhoneNumber short",
			userInput: domain.User{
				FullName:    shortFullName,
				PhoneNumber: &shortPhoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
		{
			name: "Fail create user in repository PhoneNumber long",
			userInput: domain.User{
				FullName:    shortFullName,
				PhoneNumber: &longPhoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
		{
			name: "Fail create user in repository PhoneNumber no prefix",
			userInput: domain.User{
				FullName:    shortFullName,
				PhoneNumber: &noPrefixPhoneNumber,
			},
			wantError: true,
			errorMsg:  "scan error:ERROR:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewTestPool(t)

			repo := NewUsersRepository(pool)
			if tt.setup != nil {
				tt.setup(t, repo)
			}
			ctx := context.Background()

			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			result, err := repo.UserCreate(ctx, tt.userInput)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, result, pool)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
