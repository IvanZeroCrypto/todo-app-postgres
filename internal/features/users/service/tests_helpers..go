package users_service

import (
	"errors"
	"strings"
	"testing"

	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	"gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service/mocks"
	"go.uber.org/mock/gomock"
)

var (
	emptyFullName   = ""
	fullName        = "Ivan Ivanov"
	shortFullName   = "Iv"
	longFullName    = strings.Repeat("a", 101)
	fullNameNumbers = "463473563433"
	fullNameSymbols = "%^$&#($&%^##*$)"
	fullNameArray   = "[Ivan Ivanov]"
	fullNameObject  = "{name:Ivan}"

	phoneNumber         = "+77777777777"
	shortPhoneNumber    = "+78987"
	longPhoneNumber     = "+7896786786786786786"
	symbolsPhoneNumber  = "+%^$*!@#"
	lettersPhoneNumber  = "+werfdsazxcf"
	arrayPhoneNumber    = "[+789789789789]"
	objectPhoneNumber   = "{phoneNumber:+789789789789}"
	noPrefixPhoneNumber = "77777777777"
	dbError             = errors.New("db error")
	user                = domain.User{
		ID: 1, Version: 1, FullName: fullName, PhoneNumber: &phoneNumber,
	}
	patchFullName    = "Ivan Puzov"
	patchPhoneNumber = "+789789789789"
	patchUser        = domain.User{
		ID: 1, Version: 2, FullName: patchFullName, PhoneNumber: &patchPhoneNumber,
	}
)

type testMocks struct {
	Ctrl       *gomock.Controller
	Repository *mocks.MockUsersRepository
}

func newTestMocks(t *testing.T) *testMocks {
	t.Helper()
	crtl := gomock.NewController(t)
	return &testMocks{
		Ctrl:       crtl,
		Repository: mocks.NewMockUsersRepository(crtl),
	}
}
func (t *testMocks) Finish() {
	t.Ctrl.Finish()
}
