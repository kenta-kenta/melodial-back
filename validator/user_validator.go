package validator

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kenta-kenta/diary-music/model"
)

type IUserValidator interface {
	UserValidate(user model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

// UserValidate validates the user
func (uv *userValidator) UserValidate(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(
			&user.Email,
			validation.Required.Error("Email is required"),                              // Error message when email is empty
			validation.Length(1, 30).Error("Email must be between 1 and 30 characters"), // Error message when email is not between 1 and 30 characters
			is.Email.Error("Email is invalid"),                                          // Error message when email is invalid
		),
		validation.Field(
			&user.Password,
			validation.Required.Error("Password is required"),                              // Error message when password is empty
			validation.Length(6, 30).Error("Password must be between 6 and 30 characters"), // Error message when password is not between 6 and 30 characters
		),
	)
}
