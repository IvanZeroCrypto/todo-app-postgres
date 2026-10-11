package users_transport_http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandlerPatchUser(t *testing.T) {
	oldPhoneNumber := "+7999999999"
	oldFullName := "Ivan Puzov"
	oldUser := domain.User{
		ID:          1,
		Version:     1,
		FullName:    oldFullName,
		PhoneNumber: &oldPhoneNumber,
	}
	tests := []struct {
		name          string
		request       map[string]interface{}
		mockSetup     func(*mocks.MockUsersService)
		checkResponse func(t *testing.T, response PatchUserResponse)
		userID        string
		wantStatus    int
		wantError     bool
		checkError    func(t *testing.T, body *errResp)
	}{
		{
			name: "Successful user update when FullName and PhoneNumber were sent in the request",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						assert.Equal(t, *patch.FullName.Value, "Ivan Ivanov")
						assert.Equal(t, *patch.PhoneNumber.Value, "+78888888888")
						assert.Equal(t, 1, id)
						return domain.User{
							ID:          oldUser.ID,
							Version:     oldUser.Version + 1,
							FullName:    *patch.FullName.Value,
							PhoneNumber: patch.PhoneNumber.Value,
						}, nil
					}).Times(1)
			},
			checkResponse: func(t *testing.T, response PatchUserResponse) {
				assert.Equal(t, response.ID, oldUser.ID)
				assert.Equal(t, response.Version, 2)
				assert.Equal(t, response.FullName, "Ivan Ivanov")
				assert.Equal(t, *response.PhoneNumber, "+78888888888")
			},
			wantError:  false,
			userID:     "1",
			wantStatus: http.StatusOK,
		},

		{
			name: "Successful user update when FullName but PhoneNumber = null were sent in the request",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": nil,
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						assert.Equal(t, *patch.FullName.Value, "Ivan Ivanov")
						assert.Nil(t, patch.PhoneNumber.Value)
						assert.Equal(t, 1, id)
						return domain.User{
							ID:          oldUser.ID,
							Version:     oldUser.Version + 1,
							FullName:    *patch.FullName.Value,
							PhoneNumber: nil,
						}, nil
					}).Times(1)
			},
			checkResponse: func(t *testing.T, response PatchUserResponse) {
				assert.Equal(t, response.ID, oldUser.ID)
				assert.Equal(t, response.Version, 2)
				assert.Equal(t, response.FullName, "Ivan Ivanov")
				assert.Nil(t, response.PhoneNumber)
			},
			wantError:  false,
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name: "Successful user update when the request sent FullName but the PhoneNumber field is missing",
			request: map[string]interface{}{
				"full_name": "Ivan Ivanov",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						assert.Equal(t, *patch.FullName.Value, "Ivan Ivanov")
						assert.Empty(t, patch.PhoneNumber.Value)
						assert.Equal(t, 1, id)
						return domain.User{
							ID:          oldUser.ID,
							Version:     oldUser.Version + 1,
							FullName:    *patch.FullName.Value,
							PhoneNumber: oldUser.PhoneNumber,
						}, nil
					}).Times(1)
			},
			checkResponse: func(t *testing.T, response PatchUserResponse) {
				assert.Equal(t, response.ID, oldUser.ID)
				assert.Equal(t, response.Version, 2)
				assert.Equal(t, response.FullName, "Ivan Ivanov")
				assert.Equal(t, *oldUser.PhoneNumber, *response.PhoneNumber)
			},
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name: "Successful user update when the request sent a PhoneNumber that is null but the FullName field is missing",
			request: map[string]interface{}{
				"phone_number": "+123123123123",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						phoneNumber := "+123123123123"
						assert.Equal(t, phoneNumber, *patch.PhoneNumber.Value)
						assert.Equal(t, 1, id)
						return domain.User{
							ID:          oldUser.ID,
							Version:     oldUser.Version + 1,
							FullName:    oldUser.FullName,
							PhoneNumber: &phoneNumber,
						}, nil
					}).Times(1)
			},
			checkResponse: func(t *testing.T, response PatchUserResponse) {
				assert.Equal(t, response.ID, oldUser.ID)
				assert.Equal(t, response.Version, 2)
				assert.Equal(t, response.FullName, oldFullName)
				assert.Equal(t, "+123123123123", *response.PhoneNumber)
			},
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		//	errors tests
		{
			name: "error updating user by id not found",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+123123123123",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						phoneNumber := "+123123123123"
						assert.Equal(t, "Ivan Ivanov", *patch.FullName.Value)
						assert.Equal(t, phoneNumber, *patch.PhoneNumber.Value)
						assert.Equal(t, 1, id)
						return domain.User{}, core_errors.ErrNotFound
					}).Times(1)
			},
			userID:     "1",
			wantStatus: http.StatusNotFound,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "not found")
				assert.Contains(t, body.Message, "failed to patch user")
			},
		},
		{
			name: "user update error apply patch",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+123123123123",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						phoneNumber := "+123123123123"
						assert.Equal(t, "Ivan Ivanov", *patch.FullName.Value)
						assert.Equal(t, phoneNumber, *patch.PhoneNumber.Value)
						assert.Equal(t, 1, id)
						return domain.User{}, fmt.Errorf("apply user patch")
					}).Times(1)
			},
			userID:     "1",
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "apply user patch")
				assert.Contains(t, body.Message, "failed to patch user")
			},
		},

		{
			name: "Error updating user with empty  FullName in request",
			request: map[string]interface{}{
				"full_name":    "",
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "FullName must be between 3 and 100 symbols:invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "Error updating user with FullName = null in request",
			request: map[string]interface{}{
				"full_name":    nil,
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "'FullName' cant be NULL:invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "Error updating user with short FullName",
			request: map[string]interface{}{
				"full_name":    "Iv",
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "FullName must be between 3 and 100 symbols:invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "Error updating user with long FullName",
			request: map[string]interface{}{
				"full_name":    strings.Repeat("a", 101),
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "FullName must be between 3 and 100 symbols:invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "Error updating user with symbols FullName",
			request: map[string]interface{}{
				"full_name":    "Iv%7#$@()_+=",
				"phone_number": "+78888888888",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "FullName can only be letters:invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "update error where PhoneNumber is letters",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+hfgdyrddbhs",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), 1, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, patch domain.UserPath) (domain.User, error) {
						return domain.User{}, core_errors.ErrInvalidArgument
					}).Times(1)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "invalid argument")
				assert.Contains(t, body.Message, "failed to patch user")
			},
		},
		{
			name: "update error where PhoneNumber short",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+9999999",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "update error where PhoneNumber long",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+999999999999999999999",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "invalid argument")
				assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
			},
		},
		{
			name: "long id Ddos",
			request: map[string]interface{}{
				"full_name":    "Ivan Ivanov",
				"phone_number": "+99999999999",
			},
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:     strings.Repeat("1", 50),
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Error, "value out of range:invalid argument")
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

			req := newTestRequest(t, http.MethodPatch, "/users/", tt.request)
			req.SetPathValue("id", tt.userID)
			rw := httptest.NewRecorder()

			handler.PatchUser(rw, req)
			assert.Equal(t, tt.wantStatus, rw.Code)

			if tt.wantError {
				body := getErrorMessage(t, rw)
				tt.checkError(t, body)
			} else {
				var response PatchUserResponse
				decodeResponse(t, rw, &response)
				tt.checkResponse(t, response)
			}
		})
	}
}
