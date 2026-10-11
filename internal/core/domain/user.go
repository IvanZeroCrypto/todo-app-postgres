package domain

import (
	"fmt"
	"regexp"
	"unicode"

	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUnitialized(fullname string, phoneNumber *string) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullname,
		phoneNumber,
	)
}

func (u *User) Validate() error {

	fullNameLen := len([]rune(u.FullName))

	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf("invalid `Fullname` len:%d:%w",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	for _, s := range u.FullName {
		if unicode.IsSpace(s) || s == '-' {
			continue
		}
		if !unicode.IsLetter(s) {
			return fmt.Errorf("FullName can only be letters")
		}
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len:%d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
func IsOnlyLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

type UserPath struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPath {
	return UserPath{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (p *UserPath) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("fullname cant be patched to NULL:%w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPath) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch:%w", err)
	}
	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user:%w", err)
	}
	*u = tmp
	return nil
}
