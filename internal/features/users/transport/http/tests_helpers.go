package users_transport_http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_logger "gitverse.ru/IvanVerse/go-todoApp/internal/core/logger"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http/mocks"
	"go.uber.org/mock/gomock"
)

func CreateTestLogger(t *testing.T) *core_logger.Logger {
	tempDir := t.TempDir()
	config := core_logger.Config{
		Level:  "debug",
		Folder: tempDir,
	}
	logger, err := core_logger.NewLogger(config)
	require.NoError(t, err)

	t.Cleanup(func() {
		logger.Close()
	})
	return logger
}
func newTestRequest(t *testing.T, method string, path string, body interface{}) *http.Request {
	t.Helper()

	var bodyReader *bytes.Reader

	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(data)
	} else {
		bodyReader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	testLogger := CreateTestLogger(t)
	ctx := core_logger.ToContext(context.Background(), testLogger)
	return req.WithContext(ctx)
}

func decodeResponse(t *testing.T, rw *httptest.ResponseRecorder, dest any) {
	err := json.NewDecoder(rw.Body).Decode(dest)
	require.NoError(t, err)
}

type testMocks struct {
	Ctrl    *gomock.Controller
	Service *mocks.MockUsersService
}

func newTestMocks(t *testing.T) *testMocks {
	t.Helper()
	crtl := gomock.NewController(t)
	return &testMocks{
		Ctrl:    crtl,
		Service: mocks.NewMockUsersService(crtl),
	}
}
func (m *testMocks) Finish() {
	m.Ctrl.Finish()
}

func generateTestUserDTOs(limit, offset *int) []domain.User {

	testUsers := []domain.User{
		{ID: 1, Version: 1, FullName: "User1", PhoneNumber: nil},
		{ID: 2, Version: 1, FullName: "User2", PhoneNumber: nil},
		{ID: 3, Version: 2, FullName: "User3", PhoneNumber: nil},
		{ID: 4, Version: 1, FullName: "User4", PhoneNumber: nil},
		{ID: 5, Version: 2, FullName: "User5", PhoneNumber: nil},
		{ID: 6, Version: 2, FullName: "User6", PhoneNumber: nil},
		{ID: 7, Version: 2, FullName: "User7", PhoneNumber: nil},
		{ID: 8, Version: 2, FullName: "User8", PhoneNumber: nil},
		{ID: 9, Version: 2, FullName: "User9", PhoneNumber: nil},
		{ID: 10, Version: 2, FullName: "User10", PhoneNumber: nil},
	}

	start := 0
	end := len(testUsers)

	if limit == nil && offset == nil {
		return testUsers
	}

	if offset != nil && *offset >= 0 {
		if *offset < len(testUsers) {
			start = *offset
		} else {
			return []domain.User{}
		}
	}
	if limit != nil && *limit > 0 {
		potentialEnd := *limit + start
		if potentialEnd <= len(testUsers) {
			end = *limit + start
		} else {
			end = len(testUsers)
		}
	}
	users := testUsers[start:end]
	return users
}

type errResp struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func getErrorMessage(t *testing.T, rw *httptest.ResponseRecorder) *errResp {
	t.Helper()
	var respErr errResp

	err := json.NewDecoder(rw.Body).Decode(&respErr)
	require.NoError(t, err)

	return &errResp{
		Error:   respErr.Error,
		Message: respErr.Message,
	}

}
