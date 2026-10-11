package users_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceGetUsers(t *testing.T) {
	expectedUsers := []domain.User{
		{ID: 1, Version: 1, FullName: "John1", PhoneNumber: nil},
		{ID: 2, Version: 1, FullName: "John2", PhoneNumber: nil},
		{ID: 3, Version: 1, FullName: "John3", PhoneNumber: nil},
		{ID: 4, Version: 1, FullName: "John4", PhoneNumber: nil},
		{ID: 5, Version: 1, FullName: "John5", PhoneNumber: nil},
		{ID: 6, Version: 1, FullName: "John6", PhoneNumber: nil},
		{ID: 7, Version: 1, FullName: "John7", PhoneNumber: nil},
		{ID: 8, Version: 1, FullName: "John8", PhoneNumber: nil},
	}
	tests := []struct {
		name      string
		limit     int
		offset    int
		mockSetup func(m *mocks.MockUsersRepository, limit, offset *int)
		wantError bool
		errorMsg  string
	}{
		{
			name:   "Success get users with limit = 10 offset = 0",
			limit:  10,
			offset: 0,
			mockSetup: func(m *mocks.MockUsersRepository, limit, offset *int) {
				m.EXPECT().GetUsers(gomock.Any(), limit, offset).
					DoAndReturn(func(ctx context.Context, l, o *int) ([]domain.User, error) {
						assert.Equal(t, 10, *l)
						assert.Equal(t, 0, *o)
						return expectedUsers, nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name:   "Success get users with limit = 5 offset = 3",
			limit:  5,
			offset: 3,
			mockSetup: func(m *mocks.MockUsersRepository, limit, offset *int) {
				m.EXPECT().GetUsers(gomock.Any(), limit, offset).
					DoAndReturn(func(ctx context.Context, l, o *int) ([]domain.User, error) {
						assert.Equal(t, 5, *l)
						assert.Equal(t, 3, *o)
						sliceUsers := expectedUsers[*o : *l+*o]
						return sliceUsers, nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name:   "Fail get users db error",
			limit:  10,
			offset: 3,
			mockSetup: func(m *mocks.MockUsersRepository, limit, offset *int) {
				m.EXPECT().GetUsers(gomock.Any(), limit, offset).
					DoAndReturn(func(ctx context.Context, l, o *int) ([]domain.User, error) {
						assert.Equal(t, 10, *l)
						assert.Equal(t, 3, *o)
						return []domain.User{}, dbError
					}).Times(1)
			},
			wantError: true,
			errorMsg:  "get users from repository",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			tt.mockSetup(mock.Repository, &tt.limit, &tt.offset)

			service := NewUsersService(mock.Repository)
			result, err := service.GetUsers(context.Background(), &tt.limit, &tt.offset)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(result), 0)
				assert.LessOrEqual(t, len(result), tt.limit)
			}

		})
	}
}
