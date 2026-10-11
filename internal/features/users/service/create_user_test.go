package users_service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceCreateUser(t *testing.T) {

	tests := []struct {
		name      string
		inputUser domain.User
		mockSetup func(m *mocks.MockUsersRepository, user domain.User)
		wantError bool
		errorMsg  string
	}{
		{
			name: "Success create user with FullName and PhoneNumber",
			inputUser: domain.User{
				ID:          -1,
				Version:     -1,
				FullName:    fullName,
				PhoneNumber: &phoneNumber,
			},
			mockSetup: func(m *mocks.MockUsersRepository, user domain.User) {
				m.EXPECT().UserCreate(gomock.Any(), user).
					DoAndReturn(func(ctx context.Context, u domain.User) (domain.User, error) {
						assert.Equal(t, -1, u.ID)
						assert.Equal(t, -1, u.Version)
						assert.Equal(t, fullName, u.FullName)
						assert.Equal(t, phoneNumber, *u.PhoneNumber)

						return domain.User{
							ID:          1,
							Version:     1,
							FullName:    u.FullName,
							PhoneNumber: u.PhoneNumber,
						}, nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name: "Success create user with FullName and PhoneNumber = nil",
			inputUser: domain.User{
				ID:          -1,
				Version:     -1,
				FullName:    fullName,
				PhoneNumber: nil,
			},
			mockSetup: func(m *mocks.MockUsersRepository, user domain.User) {

				m.EXPECT().UserCreate(gomock.Any(), user).
					DoAndReturn(func(ctx context.Context, u domain.User) (domain.User, error) {
						assert.Equal(t, -1, u.ID)
						assert.Equal(t, -1, u.Version)
						assert.Equal(t, fullName, u.FullName)
						assert.Nil(t, u.PhoneNumber)
						return domain.User{
							ID:          1,
							Version:     1,
							FullName:    user.FullName,
							PhoneNumber: nil,
						}, nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name: "Erorr create user db error",
			inputUser: domain.User{
				ID:          -1,
				Version:     -1,
				FullName:    fullName,
				PhoneNumber: &phoneNumber,
			},
			mockSetup: func(m *mocks.MockUsersRepository, user domain.User) {
				assert.Equal(t, -1, user.ID)
				assert.Equal(t, -1, user.Version)
				assert.Equal(t, fullName, user.FullName)
				assert.Equal(t, phoneNumber, *user.PhoneNumber)

				m.EXPECT().UserCreate(gomock.Any(), user).
					DoAndReturn(func(ctx context.Context, user domain.User) (domain.User, error) {
						dbError := errors.New("db error")
						return domain.User{}, dbError
					}).Times(1)
			},
			wantError: true,
			errorMsg:  "create user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			mock.Finish()

			tt.mockSetup(mock.Repository, tt.inputUser)

			service := NewUsersService(mock.Repository)
			result, err := service.CreateUser(context.Background(), tt.inputUser)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.inputUser.ID, result.ID)
				assert.Equal(t, tt.inputUser.Version, result.Version)
				assert.Equal(t, tt.inputUser.FullName, result.FullName)

				if tt.inputUser.PhoneNumber != nil {
					assert.Equal(t, tt.inputUser.PhoneNumber, result.PhoneNumber)
				} else {
					assert.Nil(t, result.PhoneNumber)
				}

			}

		})
	}
}

func TestServiceCreateUserValidation(t *testing.T) {
	tests := []struct {
		name        string
		fullName    string
		phoneNumber *string
	}{
		{"FullName empty", "", &phoneNumber},
		{"FullName short", shortFullName, &phoneNumber},
		{"FullName long", longFullName, &phoneNumber},
		{"FullName numbers", fullNameNumbers, &phoneNumber},
		{"FullName array", fullNameArray, &phoneNumber},
		{"FullName object", fullNameObject, &phoneNumber},
		{"FullName incorrect symbols", fullNameSymbols, &phoneNumber},

		{"PhoneNumber short", fullName, &shortPhoneNumber},
		{"PhoneNumber long", fullName, &longPhoneNumber},
		{"PhoneNumber symbols", fullName, &symbolsPhoneNumber},
		{"PhoneNumber array", fullName, &arrayPhoneNumber},
		{"PhoneNumber object", fullName, &objectPhoneNumber},
		{"PhoneNumber letters", fullName, &lettersPhoneNumber},
		{"PhoneNumber no prefix", fullName, &noPrefixPhoneNumber},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			mock.Repository.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)

			service := NewUsersService(mock.Repository)
			user := domain.User{
				ID:          -1,
				Version:     -1,
				FullName:    tt.fullName,
				PhoneNumber: tt.phoneNumber,
			}
			result, err := service.CreateUser(context.Background(), user)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "validate user domain")
			assert.Empty(t, result)
		})
	}
}
