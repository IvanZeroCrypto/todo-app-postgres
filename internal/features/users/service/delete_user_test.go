package users_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceDeleteUser(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		mockSetup func(m *mocks.MockUsersRepository, id int)
		wantError bool
		errorMsg  string
	}{
		{
			name: "Success delete user ",
			id:   1,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().DeleteUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) error {
						return nil
					}).Times(1)
			},
			wantError: false,
		},
		{
			name: "Fail get user if id <= 0",
			id:   0,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().DeleteUser(gomock.Any(), id).Times(0)
			},
			wantError: true,
			errorMsg:  "id cannot be less than zero or negative",
		},
		{
			name: "Fail delete user error db",
			id:   1,
			mockSetup: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().DeleteUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) error {
						return dbError
					}).Times(1)
			},
			wantError: true,
			errorMsg:  "failed to delete user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			tt.mockSetup(mock.Repository, tt.id)

			service := NewUsersService(mock.Repository)
			err := service.DeleteUser(context.Background(), tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.id)
			}

		})
	}
}
