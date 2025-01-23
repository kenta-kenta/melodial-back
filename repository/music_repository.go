package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kenta-kenta/diary-music/model"
	"gorm.io/gorm"
)

type IMusicRepository interface {
	CreateMusic(req *model.MusicRequest, diaryID uint) (*model.MusicResponse, error)
	SaveMusic(music *model.Music) error
	CreateMusicWithDiary(diary *model.Diary, musicReq *model.MusicRequest) (*model.MusicResponse, error)
	GetMusicsList(page int, limit int) ([]model.Music, error)
}

type musicRepository struct {
	baseURL    string
	httpClient *http.Client
	db         *gorm.DB
}

func NewMusicRepository(db *gorm.DB) IMusicRepository {
	return &musicRepository{
		baseURL:    "https://api.topmediai.com/v1",
		httpClient: &http.Client{},
		db:         db,
	}
}

func (mr *musicRepository) CreateMusic(req *model.MusicRequest, diaryID uint) (*model.MusicResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", mr.baseURL+"/music", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", os.Getenv("API_KEY"))

	response, err := mr.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result model.MusicResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	music := model.Music{
		DiaryID:      diaryID,
		Title:        result.Data[0].Title,
		AudioFile:    result.Data[0].AudioFile,
		ImageFile:    result.Data[0].ImageFile,
		ItemUUID:     result.Data[0].ItemUUID,
		Lyrics:       result.Data[0].Lyric,
		IsAuto:       req.IsAuto,
		Prompt:       req.Prompt,
		Instrumental: req.Instrumental,
	}

	if err := mr.SaveMusic(&music); err != nil {
		return nil, fmt.Errorf("failed to save music: %w", err)
	}

	return &result, nil
}

func (mr *musicRepository) SaveMusic(music *model.Music) error {
	if err := mr.db.Create(music).Error; err != nil {
		return fmt.Errorf("failed to create music record: %w", err)
	}
	return nil
}

func (mr *musicRepository) CreateMusicWithDiary(diary *model.Diary, musicReq *model.MusicRequest) (*model.MusicResponse, error) {
	var result model.MusicResponse

	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// 1. まずDiaryを保存
		if err := tx.Create(diary).Error; err != nil {
			return err
		}

		// 2. API呼び出し
		jsonData, err := json.Marshal(musicReq)
		if err != nil {
			return err
		}

		httpReq, err := http.NewRequest("POST", mr.baseURL+"/music", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("x-api-key", os.Getenv("API_KEY"))

		response, err := mr.httpClient.Do(httpReq)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return err
		}

		// 3. Music情報をDBに保存
		music := model.Music{
			DiaryID:      diary.ID, // 保存されたDiaryのIDを使用
			Title:        result.Data[0].Title,
			AudioFile:    result.Data[0].AudioFile,
			ImageFile:    result.Data[0].ImageFile,
			ItemUUID:     result.Data[0].ItemUUID,
			Lyrics:       result.Data[0].Lyric,
			IsAuto:       musicReq.IsAuto,
			Prompt:       musicReq.Prompt,
			Instrumental: musicReq.Instrumental,
		}

		if err := tx.Create(&music).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ページネーション
func (mr *musicRepository) GetMusicsList(page int, limit int) ([]model.Music, error) {
	var musics []model.Music
	if err := mr.db.Offset((page - 1) * limit).Limit(limit).Find(&musics).Error; err != nil {
		return nil, err
	}
	return musics, nil
}
