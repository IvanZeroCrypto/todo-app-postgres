package users_transport_http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandlerDeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(*mocks.MockUsersService)
		userID     string
		wantStatus int
		wantError  bool
		checkError func(t *testing.T, body *errResp)
	}{
		{
			name: "Success delete user",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), 1).Return(nil).Times(1)
			},
			userID:     "1",
			wantStatus: http.StatusNoContent,
		},
		{
			name: "User not found",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), 20).Return(core_errors.ErrNotFound).Times(1)
			},
			userID:     "20",
			wantStatus: http.StatusNotFound,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "not found")
				assert.Contains(t, body.Message, "failed to delete user")
			},
		},
		{
			name: "Error InternalServer 500 error db",
			mockSetup: func(m *mocks.MockUsersService) {
				errorDb := errors.New("error db")
				m.EXPECT().DeleteUser(gomock.Any(), 1).Return(errorDb).Times(1)
			},
			userID:     "1",
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Message, "failed to delete user")
			},
		},

		{
			name: "id = 0",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), 0).Return(core_errors.ErrInvalidArgument).Times(1)
			},
			userID:     "0",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Message, "failed to delete user")
			},
		},
		{
			name: "id with spaces",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(core_errors.ErrInvalidArgument).Times(0)
			},
			userID:     " 10 ",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid syntax:invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "very large id (DoS protection)",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(core_errors.ErrInvalidArgument).Times(0)
			},
			userID:     strings.Repeat("9", 50),
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "negative id",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), -3).Return(core_errors.ErrInvalidArgument).Times(1)
			},
			userID:     "-3",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid argument")
				assert.Contains(t, body.Message, "failed to delete user")
			},
		},
		{
			name: "float id",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1.1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid syntax:invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "id = false",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "false",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid syntax:invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "empty id ",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "array id ",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "[1]",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid syntax:invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
		{
			name: "object id ",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "{id:1}",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.NotEmpty(t, body.Error)
				assert.Contains(t, body.Error, "invalid syntax:invalid argument")
				assert.Contains(t, body.Message, "failed to get userID path value")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := newTestMocks(t)
			defer mocks.Finish()
			tt.mockSetup(mocks.Service)

			handler := NewUsersHTTPHandler(mocks.Service)

			req := newTestRequest(t, http.MethodDelete, "/users/", nil)
			req.SetPathValue("id", tt.userID)
			rw := httptest.NewRecorder()

			handler.DeleteUser(rw, req)
			assert.Equal(t, tt.wantStatus, rw.Code)

			if tt.wantError {
				body := getErrorMessage(t, rw)
				tt.checkError(t, body)
			}
		})
	}
}
