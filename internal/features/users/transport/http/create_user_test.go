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

	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http/mocks"
	"go.uber.org/mock/gomock"
)

var phoneNumber = "+78889999999"
var full_name = "Ivan Ivanov"

func TestHandlerCreateUser(t *testing.T) {
	t.Run("success test create user and error", func(t *testing.T) {
		tests := []struct {
			name              string
			request           map[string]interface{}
			setupMock         func(*mocks.MockUsersService)
			wantStatus        int
			wantResponseCheck func(*testing.T, CreateUserResponse)
		}{

			{
				name: "Successful creation of user with full_name and phonenumber",
				request: map[string]interface{}{
					"full_name":    full_name,
					"phone_number": &phoneNumber,
				},
				setupMock: func(m *mocks.MockUsersService) {
					m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user domain.User) (domain.User, error) {
							assert.Equal(t, -1, user.ID)
							assert.Equal(t, -1, user.Version)
							assert.Equal(t, full_name, user.FullName)
							assert.Equal(t, phoneNumber, *user.PhoneNumber)
							return domain.User{
								ID:          1,
								Version:     1,
								FullName:    full_name,
								PhoneNumber: &phoneNumber,
							}, nil
						}).Times(1)
				},
				wantStatus: http.StatusCreated,
				wantResponseCheck: func(t *testing.T, resp CreateUserResponse) {
					assert.Equal(t, 1, resp.ID)
					assert.Equal(t, 1, resp.Version)
					assert.Equal(t, full_name, resp.FullName)
					assert.Equal(t, phoneNumber, *resp.PhoneNumber)
				},
			},
			{
				name: "Successful creation of user only full_name",
				request: map[string]interface{}{
					"full_name": full_name,
				},
				setupMock: func(m *mocks.MockUsersService) {
					m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user domain.User) (domain.User, error) {
							assert.Equal(t, -1, user.ID)
							assert.Equal(t, -1, user.Version)
							assert.Equal(t, full_name, user.FullName)
							assert.Nil(t, user.PhoneNumber)
							return domain.User{
								ID:          1,
								Version:     1,
								FullName:    full_name,
								PhoneNumber: nil,
							}, nil
						}).Times(1)
				},
				wantStatus: http.StatusCreated,
				wantResponseCheck: func(t *testing.T, resp CreateUserResponse) {
					assert.Equal(t, 1, resp.ID)
					assert.Equal(t, 1, resp.Version)
					assert.Equal(t, full_name, resp.FullName)
					assert.Nil(t, resp.PhoneNumber)
				},
			},
			{
				name: "Successful creation of user with full_name and phone_number = nil",
				request: map[string]interface{}{
					"full_name":    full_name,
					"phone_number": nil,
				},
				setupMock: func(m *mocks.MockUsersService) {
					m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user domain.User) (domain.User, error) {
							assert.Equal(t, -1, user.ID)
							assert.Equal(t, -1, user.Version)
							assert.Equal(t, full_name, user.FullName)
							assert.Nil(t, user.PhoneNumber)
							return domain.User{
								ID:          1,
								Version:     1,
								FullName:    full_name,
								PhoneNumber: nil,
							}, nil
						}).Times(1)
				},
				wantStatus: http.StatusCreated,
				wantResponseCheck: func(t *testing.T, resp CreateUserResponse) {
					assert.Equal(t, 1, resp.ID)
					assert.Equal(t, 1, resp.Version)
					assert.Equal(t, full_name, resp.FullName)
					assert.Nil(t, resp.PhoneNumber)
				},
			},
			{
				name: "Fail service in return InternalServerError",
				request: map[string]interface{}{
					"full_name":    full_name,
					"phone_number": &phoneNumber,
				},
				setupMock: func(m *mocks.MockUsersService) {
					m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, user domain.User) (domain.User, error) {
							assert.Equal(t, -1, user.ID)
							assert.Equal(t, -1, user.Version)
							assert.Equal(t, full_name, user.FullName)
							assert.Equal(t, phoneNumber, *user.PhoneNumber)
							return domain.User{}, fmt.Errorf("failed to create user")
						}).Times(1)
				},
				wantStatus:        http.StatusInternalServerError,
				wantResponseCheck: nil,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mocks := newTestMocks(t)
				tt.setupMock(mocks.Service)

				handler := NewUsersHTTPHandler(mocks.Service)

				req := newTestRequest(t, http.MethodPost, "/users", tt.request)

				rw := httptest.NewRecorder()

				handler.CreateUser(rw, req)

				assert.Equal(t, tt.wantStatus, rw.Code)
				if tt.wantResponseCheck != nil {
					var response CreateUserResponse

					if rw.Code == http.StatusCreated {
						// err := json.NewDecoder(rw.Body).Decode(&response)
						// require.NoError(t, err)
						decodeResponse(t, rw, &response)
						tt.wantResponseCheck(t, response)
					}
				}

			})
		}
	})
}

func TestHandlerCreateUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		request map[string]interface{}
	}{
		{
			name: "FullName empty",
			request: map[string]interface{}{
				"full_name":    "",
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName bool",
			request: map[string]interface{}{
				"full_name":    true,
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName int",
			request: map[string]interface{}{
				"full_name":    1234,
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName float",
			request: map[string]interface{}{
				"full_name":    1234.456,
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName null",
			request: map[string]interface{}{
				"full_name":    nil,
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName array",
			request: map[string]interface{}{
				"full_name":    []string{"Ivan Ivanov"},
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName object",
			request: map[string]interface{}{
				"full_name":    map[string]string{"name": "Ivan Ivanov"},
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName short",
			request: map[string]interface{}{
				"full_name":    "Iv",
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName long",
			request: map[string]interface{}{
				"full_name":    strings.Repeat("a", 101),
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "FullName missing",
			request: map[string]interface{}{
				"phone_number": &phoneNumber,
			},
		},
		{
			name: "PhoneNumber has no prefix +",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": "78889999999",
			},
		},
		{
			name: "PhoneNumber long",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": "+788888888888888888",
			},
		},
		{
			name: "PhoneNumber string",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": "+number",
			},
		},
		{
			name: "PhoneNumber int",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": +78998773761,
			},
		},
		{
			name: "PhoneNumber bool",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": true,
			},
		},
		{
			name: "PhoneNumber float",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": +45.65,
			},
		},
		{
			name: "PhoneNumber array",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": []string{"+78998773761"},
			},
		},
		{
			name: "PhoneNumber object",
			request: map[string]interface{}{
				"full_name":    full_name,
				"phone_number": map[string]string{"phone": "+78998773761"},
			},
		},
		{
			name:    "Empty valid JSON",
			request: map[string]interface{}{},
		},
		{
			name:    "Invalid JSON",
			request: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := newTestMocks(t)
			defer mocks.Finish()

			handler := NewUsersHTTPHandler(mocks.Service)

			req := newTestRequest(t, http.MethodPost, "/users", tt.request)
			rw := httptest.NewRecorder()
			t.Logf("Response body: %q", rw.Body.String())
			t.Logf("Response body length: %d", rw.Body.Len())

			mocks.Service.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).Times(0)

			handler.CreateUser(rw, req)

			assert.Equal(t, http.StatusBadRequest, rw.Code)

			body := getErrorMessage(t, rw)
			assert.NotEmpty(t, body.Error)
			assert.Contains(t, body.Message, "failed to decode and validate HTTP request")
		})
	}

}
