package users_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestServicePatchUser(t *testing.T) {
	tests := []struct {
		name               string
		id                 int
		patch              domain.UserPath
		mockSetupGetUser   func(m *mocks.MockUsersRepository, id int)
		mockSetupPatchUser func(m *mocks.MockUsersRepository)
		wantError          bool
		errorMsg           string
	}{
		{
			name: "Success patch user update FullName and PhoneNumber",
			id:   1,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Set: true, Value: &patchFullName},
				domain.Nullable[string]{Set: true, Value: &patchPhoneNumber},
			),
			mockSetupGetUser: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return user, nil
					}).Times(1)
			},
			mockSetupPatchUser: func(m *mocks.MockUsersRepository) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, user domain.User) (domain.User, error) {
						assert.Equal(t, 1, user.ID)
						assert.Equal(t, 1, user.Version)
						assert.Equal(t, patchFullName, user.FullName)
						assert.Equal(t, patchPhoneNumber, *user.PhoneNumber)
						return patchUser, nil
					}).Times(1)
			},

			wantError: false,
		},
		{
			name: "Success patch user update FullName ",
			id:   1,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Set: true, Value: &patchFullName},
				domain.Nullable[string]{Set: false, Value: nil},
			),
			mockSetupGetUser: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return user, nil
					}).Times(1)
			},
			mockSetupPatchUser: func(m *mocks.MockUsersRepository) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, user domain.User) (domain.User, error) {
						assert.Equal(t, 1, user.ID)
						assert.Equal(t, 1, user.Version)
						assert.Equal(t, patchFullName, user.FullName)
						assert.Equal(t, phoneNumber, *user.PhoneNumber)
						return patchUser, nil
					}).Times(1)
			},

			wantError: false,
		},
		{
			name: "Success patch user update PhoneNumber ",
			id:   1,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Set: false, Value: nil},
				domain.Nullable[string]{Set: true, Value: &patchPhoneNumber},
			),
			mockSetupGetUser: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return user, nil
					}).Times(1)
			},
			mockSetupPatchUser: func(m *mocks.MockUsersRepository) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, user domain.User) (domain.User, error) {
						assert.Equal(t, 1, user.ID)
						assert.Equal(t, 1, user.Version)
						assert.Equal(t, fullName, user.FullName)
						assert.Equal(t, patchPhoneNumber, *user.PhoneNumber)
						return patchUser, nil
					}).Times(1)
			},

			wantError: false,
		},
		{
			name: "Fail update user with id not found ",
			id:   1,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Set: true, Value: &patchFullName},
				domain.Nullable[string]{Set: true, Value: &patchPhoneNumber},
			),
			mockSetupGetUser: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return domain.User{}, core_errors.ErrNotFound
					}).Times(1)
			},
			mockSetupPatchUser: func(m *mocks.MockUsersRepository) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			wantError: true,
			errorMsg:  "get user:not found",
		},
		{
			name: "Fail update user error db",
			id:   1,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{Set: true, Value: &patchFullName},
				domain.Nullable[string]{Set: true, Value: &patchPhoneNumber},
			),
			mockSetupGetUser: func(m *mocks.MockUsersRepository, id int) {
				m.EXPECT().GetUser(gomock.Any(), id).
					DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
						return user, nil
					}).Times(1)
			},
			mockSetupPatchUser: func(m *mocks.MockUsersRepository) {
				m.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int, user domain.User) (domain.User, error) {
						return domain.User{}, dbError
					}).Times(1)
			},
			wantError: true,
			errorMsg:  "patch user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			tt.mockSetupGetUser(mock.Repository, tt.id)
			tt.mockSetupPatchUser(mock.Repository)

			service := NewUsersService(mock.Repository)
			result, err := service.PatchUser(context.Background(), tt.id, tt.patch)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.id)
				assert.Equal(t, 2, result.Version)
				assert.Equal(t, patchFullName, result.FullName)
				assert.Equal(t, patchPhoneNumber, *result.PhoneNumber)
			}

		})
	}
}

func TestServicePatchUserValidate(t *testing.T) {
	tests := []struct {
		name  string
		patch domain.UserPath
	}{
		{"FullName empty", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &emptyFullName},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName short", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &shortFullName},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName long", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &longFullName},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName numbers", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &fullNameNumbers},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName array", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &fullNameArray},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName object", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &fullNameObject},
			domain.Nullable[string]{Set: false, Value: nil},
		)},
		{"FullName incorrect symbols", domain.NewUserPatch(
			domain.Nullable[string]{Set: true, Value: &fullNameSymbols},
			domain.Nullable[string]{Set: false, Value: &patchPhoneNumber},
		)},

		{"PhoneNumber short", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &shortPhoneNumber},
		)},
		{"PhoneNumber long", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &longPhoneNumber},
		)},
		{"PhoneNumber symbols", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &symbolsPhoneNumber},
		)},
		{"PhoneNumber array", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &arrayPhoneNumber},
		)},
		{"PhoneNumber object", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &objectPhoneNumber},
		)},
		{"PhoneNumber letters", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &lettersPhoneNumber},
		)},
		{"PhoneNumber no prefix", domain.NewUserPatch(
			domain.Nullable[string]{Set: false, Value: nil},
			domain.Nullable[string]{Set: true, Value: &noPrefixPhoneNumber},
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newTestMocks(t)
			defer mock.Finish()

			mock.Repository.EXPECT().GetUser(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id int) (domain.User, error) {
					return user, nil
				}).Times(1)
			mock.Repository.EXPECT().PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			service := NewUsersService(mock.Repository)
			result, err := service.PatchUser(context.Background(), 1, tt.patch)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "apply user patch")
			assert.Empty(t, result)
		})
	}
}
