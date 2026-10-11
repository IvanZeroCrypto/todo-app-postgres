package users_transport_http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandlerGetUsersSuccessANDErrors(t *testing.T) {
	tests := []struct {
		name              string
		mockSetup         func(*mocks.MockUsersService)
		path              string
		wantStatus        int
		wantError         bool
		wantResponseCheck func(*testing.T, GetUsersResponse)
		checkError        func(t *testing.T, body *errResp)
	}{
		{
			name: "Success - no limit and offset with context check",
			mockSetup: func(m *mocks.MockUsersService) {
				users := generateTestUserDTOs(nil, nil)
				m.EXPECT().GetUsers(gomock.Any(), nil, nil).
					DoAndReturn(func(ctx context.Context, limit, offset *int) ([]domain.User, error) {
						assert.NotNil(t, ctx)
						assert.NoError(t, ctx.Err())
						return users, nil
					}).Times(1)
			},
			path:       "",
			wantStatus: http.StatusOK,
			wantResponseCheck: func(t *testing.T, response GetUsersResponse) {
				assert.Len(t, response, 10)
				for i := 1; i < len(response); i++ {
					assert.Equal(t, i, response[i-1].ID)
				}
			},
		},
		{
			name: "Success - only limit ",
			mockSetup: func(m *mocks.MockUsersService) {
				limit := 5
				users := generateTestUserDTOs(&limit, nil)
				m.EXPECT().GetUsers(gomock.Any(), &limit, nil).
					Return(users, nil).
					Times(1)
			},
			path:       "?limit=5",
			wantStatus: http.StatusOK,
			wantResponseCheck: func(t *testing.T, response GetUsersResponse) {
				assert.GreaterOrEqual(t, len(response), 0)
				assert.LessOrEqual(t, len(response), 5)
				for i := 1; i < len(response); i++ {
					assert.Equal(t, i, response[i-1].ID)
				}

			},
		},

		{
			name: "Success - only offset",
			mockSetup: func(m *mocks.MockUsersService) {
				offset := 2
				users := generateTestUserDTOs(nil, &offset)
				m.EXPECT().GetUsers(gomock.Any(), nil, &offset).
					Return(users, nil).
					Times(1)
			},
			path:       "?offset=2",
			wantStatus: http.StatusOK,
			wantResponseCheck: func(t *testing.T, response GetUsersResponse) {
				assert.GreaterOrEqual(t, len(response), 0)
				assert.LessOrEqual(t, len(response), 8)
				for i := 3; i < len(response); i++ {
					assert.Equal(t, i, response[i-3].ID)
				}
			},
		},
		{
			name: "Success - with limit and offset ",
			mockSetup: func(m *mocks.MockUsersService) {
				limit := 5
				offset := 3
				users := generateTestUserDTOs(&limit, &offset)
				m.EXPECT().GetUsers(gomock.Any(), &limit, &offset).
					Return(users, nil).
					Times(1)
			},
			path:       "?limit=5&offset=3",
			wantStatus: http.StatusOK,
			wantResponseCheck: func(t *testing.T, response GetUsersResponse) {
				assert.GreaterOrEqual(t, len(response), 0)
				assert.LessOrEqual(t, len(response), 5)

				for i := 4; i < len(response); i++ {
					assert.Equal(t, i, response[i-4].ID)
				}
			},
		},
		{
			name: "Should return 400 when limit/offset are invalid - service not called",
			mockSetup: func(m *mocks.MockUsersService) {
				m.EXPECT().GetUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			path:       "?limit=-5&offset=-3",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body, "failed to get 'limit/offset' query ")
			},
		},
		{
			name: "Should return 500 with correct message when database error occurs",
			mockSetup: func(m *mocks.MockUsersService) {
				limit := 5
				offset := 3
				dbError := errors.New("db error")
				m.EXPECT().GetUsers(gomock.Any(), &limit, &offset).
					Return(nil, dbError).
					Times(1)
			},
			path:       "?limit=5&offset=3",
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			checkError: func(t *testing.T, body *errResp) {
				assert.Contains(t, body.Message, "failed to get users")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			mocks := newTestMocks(t)
			defer mocks.Finish()
			tt.mockSetup(mocks.Service)
			handler := NewUsersHTTPHandler(mocks.Service)

			req := newTestRequest(t, http.MethodGet, "/users"+tt.path, nil)
			rw := httptest.NewRecorder()

			handler.GetUsers(rw, req)

			assert.Equal(t, tt.wantStatus, rw.Code)

			if tt.wantResponseCheck != nil {
				var response GetUsersResponse
				decodeResponse(t, rw, &response)
				tt.wantResponseCheck(t, response)
			}

			if tt.wantError {
				body := getErrorMessage(t, rw)
				tt.checkError(t, body)
			}

		})
	}
}

func TestHandlerGetUsersValidadtion(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"limit negative", "?limit=-10"},
		{"limit bool", "?limit=true"},
		{"limit string", "?limit=limit"},
		{"limit float", "?limit=12.45"},
		{"limit array", "?limit=[]string{1}"},
		{"limit object", "?limit={limit:43}"},
		{"offset negative", "?offset=-20"},
		{"offset bool", "?offset=false"},
		{"offset float", "?offset=203.45"},
		{"offset string", "?offset=offset"},
		{"offset array", "?offset=[]string{20}"},
		{"offset object", "?offset={offset:59}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := newTestMocks(t)
			defer mocks.Finish()

			handler := NewUsersHTTPHandler(mocks.Service)

			req := newTestRequest(t, http.MethodGet, "/users"+tt.path, nil)
			rw := httptest.NewRecorder()

			mocks.Service.EXPECT().GetUsers(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			handler.GetUsers(rw, req)

			assert.Equal(t, http.StatusBadRequest, rw.Code)

		})
	}
}
