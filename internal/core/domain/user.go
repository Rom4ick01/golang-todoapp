package domain

import (
	"fmt"
	"regexp"

	core_errors "todoapp/internal/core/errors"

	"github.com/google/uuid"
)

// User — пользователь системы.
type User struct {
	ID      uuid.UUID
	Version int

	FullName    string
	PhoneNumber *string
}

// NewUser восстанавливает пользователя из готовых данных.
func NewUser(
	id uuid.UUID,
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

// CreateUser создаёт нового пользователя с новым ID и версией 1.
func CreateUser(
	fullName string,
	phoneNumber *string,
) User {
	var (
		id      = uuid.New()
		version = 1
	)

	return NewUser(
		id,
		version,
		fullName,
		phoneNumber,
	)
}

// Validate проверяет инварианты пользователя.
// Телефон: «+» и цифры, 10–15 символов, например +79001234567.
func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// UserPatch — набор изменений для PATCH-запроса.
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// Validate запрещает обнулять обязательное поле FullName.
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
