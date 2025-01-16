package controller

import (
	"net/http"

	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/usecase"
	"github.com/labstack/echo/v4"
)

type IMusicController interface {
	CreateMusic(c echo.Context) error
}

type MusicController struct {
	mu usecase.IMusicUsecase
}

func NewMusicController(mu usecase.IMusicUsecase) *MusicController {
	return &MusicController{mu}
}

func (mc *MusicController) CreateMusic(c echo.Context) error {
	// JSONリクエストをバインド
	request := new(model.MusicRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// バリデーション
	if request.Prompt == "" || request.Lyrics == "" || request.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Prompt, lyrics and title are required",
		})
	}

	// デフォルト値の設定
	if request.IsAuto == 0 {
		request.IsAuto = 1
	}
	if request.Instrumental == 0 {
		request.Instrumental = 1
	}

	// usecaseの呼び出し
	response, err := mc.mu.CreateMusic(request.Prompt, request.Lyrics, request.Title, request.IsAuto, request.Instrumental)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response)
}
