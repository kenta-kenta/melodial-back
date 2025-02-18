package repository

import (
	"fmt"
	"math"

	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/service"
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
	CreateDiaryWithMusic(diary *model.Diary, musicReq *model.MusicRequest) (*model.DiaryResponse, error)
}

type diaryRepository struct {
	db           *gorm.DB
	MusicService service.IMusicService
}

func NewDiaryRepository(db *gorm.DB) IDiaryRepository {
	return &diaryRepository{db, service.NewMusicService()}
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
	if err := dr.db.Preload("Music").Where("user_id = ?", userId).
		Offset(offset).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&diaries).Error; err != nil {
		return nil, err
	}
	// 総ページ数を計算
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	// DiaryResponseのスライスを作成
	var diaryResponses []model.DiaryResponse
	for _, diary := range diaries {
		var musicData []model.MusicData
		for _, music := range diary.Music {
			musicData = append(musicData, model.MusicData{
				AudioFile: music.AudioFile,
				ImageFile: music.ImageFile,
				ItemUUID:  music.ItemUUID,
				Title:     music.Title,
				Lyric:     music.Lyrics,
				Tags:      music.Tags,
			})
		}
		diaryResponses = append(diaryResponses, model.DiaryResponse{
			ID:        diary.ID,
			Content:   diary.Content,
			MusicData: musicData,
			CreatedAt: diary.CreatedAt,
			UpdatedAt: diary.UpdatedAt,
		})
	}

	return &model.PaginationResponse{
		DiaryResponse: diaryResponses,
		TotalItems:    total,
		Page:          query.Page,
		PageSize:      query.PageSize,
		TotalPages:    totalPages,
	}, nil
}

func (dr *diaryRepository) GetDiaryById(diary *model.Diary, userId uint, diaryId uint) error {
	// Joinメソッドを使ってUserテーブルと結合し、Preloadメソッドを使ってMusicデータを事前にロード
	if err := dr.db.Joins("JOIN users ON users.id = diaries.user_id").
		Preload("Music").
		Where("diaries.user_id = ? AND diaries.id = ?", userId, diaryId).
		First(diary).Error; err != nil {
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

func (dr *diaryRepository) CreateDiaryWithMusic(diary *model.Diary, musicReq *model.MusicRequest) (*model.DiaryResponse, error) {
	var diaryRes *model.DiaryResponse
	err := dr.db.Transaction(func(tx *gorm.DB) error {
		// 1. 日記を保存
		if err := tx.Create(diary).Error; err != nil {
			return err
		}

		// 2. 音楽を生成・保存
		musicReq.Prompt = diary.Content
		music, err := dr.MusicService.CreateMusic(musicReq.Prompt)
		if err != nil {
			return err
		}

		music.DiaryID = diary.ID
		music.UserID = diary.UserId
		if err := tx.Create(music).Error; err != nil {
			return err
		}

		diaryRes = &model.DiaryResponse{
			ID:      diary.ID,
			Content: diary.Content,
			MusicData: []model.MusicData{
				{
					AudioFile: music.AudioFile,
					ImageFile: music.ImageFile,
					ItemUUID:  music.ItemUUID,
					Title:     music.Title,
					Lyric:     music.Lyrics,
					Tags:      music.Tags,
				},
			},
			CreatedAt: diary.CreatedAt,
			UpdatedAt: diary.UpdatedAt,
		}

		return nil
	})

	return diaryRes, err
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
	return dr.db.Transaction(func(tx *gorm.DB) error {
		// 1. まず関連するMusicレコードを削除
		if err := tx.Where("diary_id = ?", diaryId).Delete(&model.Music{}).Error; err != nil {
			return err
		}

		// 2. 次にDiaryレコードを削除
		result := tx.Where("user_id = ? AND id = ?", userId, diaryId).Delete(&model.Diary{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected < 1 {
			return fmt.Errorf("日記が見つかりません")
		}

		return nil
	})
}
