package repository

import (
	"fmt"
	"math"

	"github.com/kenta-kenta/diary-music/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IDiaryRepository interface {
	GetAllDiaries(query *model.PaginationQuery, userId uint) (*model.PaginationResponse, error)
	GetDiaryById(diary *model.Diary, userId uint, diaryId uint) error
	CreateDiary(diary *model.Diary) error
	UpdateDiary(diary *model.Diary, userId uint, diaryId uint) error
	DeleteDiary(userId uint, diaryId uint) error
	GetDiaryDates(userId uint, year, month int) ([]model.DiaryDateCount, error)
}

type diaryRepository struct {
	db *gorm.DB
}

func NewDiaryRepository(db *gorm.DB) IDiaryRepository {
	return &diaryRepository{db}
}

func (dr *diaryRepository) GetAllDiaries(query *model.PaginationQuery, userId uint) (*model.PaginationResponse, error) {
	var diaries []model.Diary
	var total int64
	// ページ番号からオフセットを計算
	offset := (query.Page - 1) * query.PageSize
	// Countメソッドを使ってデータの総数を取得
	if err := dr.db.Model(&model.Diary{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, err
	}
	// Whereメソッドを使ってデータを取得
	if err := dr.db.Where("user_id = ?", userId).Offset(offset).Limit(query.PageSize).Order("created_at DESC").Find(&diaries).Error; err != nil {
		return nil, err
	}
	// 総ページ数を計算
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &model.PaginationResponse{
		Data:       diaries,
		TotalItems: total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (dr *diaryRepository) GetDiaryById(diary *model.Diary, userId uint, diaryId uint) error {
	// Joinメソッドを使ってUserテーブルと結合
	if err := dr.db.Joins("User").Where("user_id = ?", userId).First(diary, diaryId).Error; err != nil {
		return err
	}
	return nil
}

func (dr *diaryRepository) GetDiaryDates(userId uint, year, month int) ([]model.DiaryDateCount, error) {
	var results []model.DiaryDateCount

	err := dr.db.Model(&model.Diary{}).
		Select("TO_CHAR(DATE(created_at), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("user_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND EXTRACT(MONTH FROM created_at) = ?",
			userId, year, month).
		Group("DATE(created_at)").
		Order("date").
		Scan(&results).
		Error

	if err != nil {
		return nil, fmt.Errorf("日記の日付と件数の取得に失敗: %w", err)
	}

	return results, nil
}

func (dr *diaryRepository) CreateDiary(diary *model.Diary) error {
	// Createメソッドを使ってデータを作成
	if err := dr.db.Create(diary).Error; err != nil {
		return err
	}
	return nil
}

func (dr *diaryRepository) UpdateDiary(diary *model.Diary, userId uint, diaryId uint) error {
	// Returningメソッドを使って更新後のデータを取得
	result := dr.db.Model(diary).Clauses(clause.Returning{}).Where("user_id = ? AND id = ?", userId, diaryId).Update("content", diary.Content)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (dr *diaryRepository) DeleteDiary(userId uint, diaryId uint) error {
	// Deleteメソッドを使ってデータを削除
	result := dr.db.Where("user_id = ? AND id = ?", userId, diaryId).Delete(&model.Diary{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
