package users_transport_http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
	core_http_request "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/request"
)

func TestDecodeAndValidateRequest(t *testing.T) {
	longName := strings.Repeat("a", 101)
	tests := []struct {
		name        string
		body        string
		dest        any
		wantErr     bool
		method      string
		userId      string
		errContains string
	}{
		{
			name:    "success decode",
			body:    `{"full_name":"John Doe","phone_number":"+123456789011"}`,
			dest:    &CreateUserRequest{},
			wantErr: false,
		},
		{
			name:        "error decode invalid json",
			body:        `"full_name":"John Doe","phone_number":"+123456789011`,
			dest:        &CreateUserRequest{},
			wantErr:     true,
			errContains: "decode json",
		},
		{
			name:        "error decode invalid dest",
			body:        `{"full_name":"John Doe","phone_number":"+123456789011"}`,
			dest:        nil,
			wantErr:     true,
			errContains: "decode json",
		},
		{
			name:        "error validate FullName short",
			body:        `{"full_name":"Jo","phone_number":"+123456789011"}`,
			dest:        &CreateUserRequest{},
			wantErr:     true,
			errContains: "request validation",
		},
		{
			name:        "error validate FullName long",
			body:        `{"full_name":"` + longName + `","phone_number":"+123456789011"}`,
			dest:        &CreateUserRequest{},
			wantErr:     true,
			errContains: "request validation",
		},
		{
			name:        "error validate PhoneNumber short",
			body:        `{"full_name":"John Doe","phone_number":"+123"}`,
			dest:        &CreateUserRequest{},
			wantErr:     true,
			errContains: "request validation",
		},
		{
			name:        "error validate PhoneNumber long",
			body:        `{"full_name":"John Doe","phone_number":"+123456789234653834756473243"}`,
			dest:        &CreateUserRequest{},
			wantErr:     true,
			errContains: "request validation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(tt.body)))
			err := core_http_request.DecodeAndValidateRequest(req, tt.dest)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.True(t, errors.Is(err, core_errors.ErrInvalidArgument))
			} else {
				require.NoError(t, err)
			}

		})
	}
}
