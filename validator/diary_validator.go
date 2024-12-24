package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kenta-kenta/diary-music/model"
)

type IDiaryValidator interface {
	DiaryValidate(diary model.Diary) error
}

type diaryValidator struct{}

func NewDiaryValidator() IDiaryValidator {
	return &diaryValidator{}
}

func (dv *diaryValidator) DiaryValidate(diary model.Diary) error {
	return validation.ValidateStruct(&diary,
		validation.Field(
			&diary.Content,
			validation.Required.Error("Content is required"),
			validation.RuneLength(1, 1000).Error("Content must be between 1 and 1000 characters"),
		),
	)
}
