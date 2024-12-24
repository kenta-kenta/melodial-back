package usecase

import (
	"strconv"

	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/repository"
	"github.com/kenta-kenta/diary-music/validator"
)

type IDiaryUsecase interface {
	GetAllDiaries(userId uint, page int, pageSize int) (*model.PaginationResponse, error)
	GetDiaryById(userId uint, diaryId uint) (model.DiaryResponse, error)
	CreateDiary(diary model.Diary) (model.DiaryResponse, error)
	UpdateDiary(userId uint, diaryId uint, diary model.Diary) (model.DiaryResponse, error)
	DeleteDiary(userId uint, diaryId uint) error
	GetDiaryDates(userId uint, year, month string) (*model.DiaryDateCountResponse, error)
}

type diaryUsecase struct {
	dr repository.IDiaryRepository
	dv validator.IDiaryValidator
}

func NewDiaryUsecase(dr repository.IDiaryRepository, dv validator.IDiaryValidator) IDiaryUsecase {
	return &diaryUsecase{dr, dv}
}

func (du *diaryUsecase) GetAllDiaries(userId uint, page, pageSize int) (*model.PaginationResponse, error) {
	// PaginationQueryを作成
	query := &model.PaginationQuery{
		Page:     page,
		PageSize: pageSize,
	}

	return du.dr.GetAllDiaries(query, userId)
}

func (du *diaryUsecase) GetDiaryById(userId uint, diaryId uint) (model.DiaryResponse, error) {
	diary := model.Diary{}
	if err := du.dr.GetDiaryById(&diary, userId, diaryId); err != nil {
		return model.DiaryResponse{}, err
	}
	resDiary := model.DiaryResponse{
		ID:        diary.ID,
		Content:   diary.Content,
		CreatedAt: diary.CreatedAt,
		UpdatedAt: diary.UpdatedAt,
	}
	return resDiary, nil
}

func (du *diaryUsecase) GetDiaryDates(userId uint, year, month string) (*model.DiaryDateCountResponse, error) {
	y, _ := strconv.Atoi(year)
	m, _ := strconv.Atoi(month)

	data, err := du.dr.GetDiaryDates(userId, y, m)
	if err != nil {
		return nil, err
	}

	var diaryDates []model.DiaryDateCount
	for _, d := range data {
		diaryDates = append(diaryDates, model.DiaryDateCount{
			Date:  d.Date,
			Count: d.Count,
		})
	}

	return &model.DiaryDateCountResponse{
		Dates: diaryDates,
	}, nil
}

func (du *diaryUsecase) CreateDiary(diary model.Diary) (model.DiaryResponse, error) {
	// Validate the diary
	if err := du.dv.DiaryValidate(diary); err != nil {
		return model.DiaryResponse{}, err
	}
	// Create the diary
	if err := du.dr.CreateDiary(&diary); err != nil {
		return model.DiaryResponse{}, err
	}
	resDiary := model.DiaryResponse{
		ID:        diary.ID,
		Content:   diary.Content,
		CreatedAt: diary.CreatedAt,
		UpdatedAt: diary.UpdatedAt,
	}
	return resDiary, nil
}

func (du *diaryUsecase) UpdateDiary(userId uint, diaryId uint, diary model.Diary) (model.DiaryResponse, error) {
	if err := du.dv.DiaryValidate(diary); err != nil {
		return model.DiaryResponse{}, err
	}

	if err := du.dr.UpdateDiary(&diary, userId, diaryId); err != nil {
		return model.DiaryResponse{}, err
	}
	resDiary := model.DiaryResponse{
		ID:        diary.ID,
		Content:   diary.Content,
		CreatedAt: diary.CreatedAt,
		UpdatedAt: diary.UpdatedAt,
	}
	return resDiary, nil
}

func (dr *diaryUsecase) DeleteDiary(userId uint, diaryId uint) error {
	if err := dr.dr.DeleteDiary(userId, diaryId); err != nil {
		return err
	}
	return nil
}
