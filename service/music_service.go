package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kenta-kenta/diary-music/model"
)

type IMusicService interface {
	CreateMusic(prompt string) (*model.Music, error)
}

type musicService struct {
	baseURL    string
	httpClient *http.Client
}

func NewMusicService() IMusicService {
	return &musicService{
		baseURL:    "https://api.topmediai.com/v1",
		httpClient: &http.Client{},
	}
}

func (s *musicService) CreateMusic(prompt string) (*model.Music, error) {
	// リクエストの作成
	req := &model.MusicRequest{
		Prompt: prompt,
		IsAuto: 1,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエストの作成
	httpReq, err := http.NewRequest("POST", s.baseURL+"/music", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", os.Getenv("API_KEY"))

	// リクエストの実行
	response, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	// レスポンスの解析
	var result model.MusicResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Musicモデルの作成
	music := &model.Music{
		Title:        result.Data[0].Title,
		AudioFile:    result.Data[0].AudioFile,
		ImageFile:    result.Data[0].ImageFile,
		ItemUUID:     result.Data[0].ItemUUID,
		Lyrics:       result.Data[0].Lyric,
		IsAuto:       req.IsAuto,
		Prompt:       req.Prompt,
		Instrumental: req.Instrumental,
	}

	return music, nil
}
