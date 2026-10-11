package users_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceGetUser(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		mockSetup func(m *mocks.MockUsersRepository, id int)
		wantError bool
		errorMsg  string
	}{
		{
			name: "Success get user ",
			id:   1,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return user, nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name: "Fail get user if id <= 0",
			id:   0,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).Times(0)
			},
			wantError: true,
			errorMsg:  "id cannot be negative and is equal to zero",
		},
		{
			name: "Fail get user error db",
			id:   1,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return domain.User{}, dbError
					}).Times(1)
			},
			wantError: true,
			errorMsg:  "get user from repository",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			tt.mockSetup(mock.Repository, tt.id)

			service := NewUsersService(mock.Repository)
			result, err := service.GetUser(context.Background(), tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.id)
				assert.Equal(t, 1, result.Version)
				assert.Equal(t, user.FullName, result.FullName)
				assert.Equal(t, user.PhoneNumber, result.PhoneNumber)
			}

		})
	}
}
