package usecase

import (
	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/repository"
)

type IMusicUsecase interface {
	CreateMusic(prompt, lyrics, title string, isAuto, instrumental int, diaryID uint) (*model.MusicResponse, error)
	GetMusicsList(page int, limit int, userId uint) ([]model.Music, error)
}

type MusicUsecase struct {
	mr repository.IMusicRepository
}

func NewMusicUsecase(mr repository.IMusicRepository) *MusicUsecase {
	return &MusicUsecase{mr}
}

func (mu *MusicUsecase) CreateMusic(prompt, lyrics, title string, isAuto, instrumental int, diaryID uint) (*model.MusicResponse, error) {
	request := &model.MusicRequest{
		IsAuto:       isAuto,
		Prompt:       prompt,
		Lyrics:       lyrics,
		Title:        title,
		Instrumental: instrumental,
	}

	return mu.mr.CreateMusic(request, diaryID)
}

func (mu *MusicUsecase) GetMusicsList(page int, limit int, userId uint) ([]model.Music, error) {
	return mu.mr.GetMusicsList(page, limit, userId)
}
