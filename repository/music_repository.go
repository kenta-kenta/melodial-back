package repository

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/kenta-kenta/diary-music/model"
)

type IMusicRepository interface {
	CreateMusic(req *model.MusicRequest) (*model.MusicResponse, error)
}

type musicRepository struct {
	baseURL    string
	httpClient *http.Client
}

func NewMusicRepository() *musicRepository {
	return &musicRepository{
		baseURL:    "https://api.topmediai.com/v1",
		httpClient: &http.Client{},
	}
}

func (mr *musicRepository) CreateMusic(req *model.MusicRequest) (*model.MusicResponse, error) {
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

	return &result, nil
}
