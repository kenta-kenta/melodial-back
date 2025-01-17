package controller

import (
	"net/http"
	"strconv"

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
	diaryID := c.QueryParam("diary_id")
	if diaryID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Diary ID is required",
		})
	}

	dID, err := strconv.ParseUint(diaryID, 10, 32)
	// JSONリクエストをバインド
	request := new(model.MusicRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// バリデーション
	if request.Prompt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Prompt, lyrics and title are required",
		})
	}

	// デフォルト値の設定
	if request.IsAuto == 0 {
		request.IsAuto = 1
	}
	if request.Instrumental == 0 {
		request.Instrumental = 0
	}

	// テストデータ
	// response := model.MusicResponse{
	// 	Status:  200,
	// 	Message: "Success",
	// 	Data: []model.MusicData{
	// 		{
	// 			AudioFile: "https://files.topmediai.com/aimusic/api/be2151dc-cbff-4be6-ab1b-afa02ea772f1-audio.mp3",
	// 			ImageFile: "https://files.topmediai.com/aimusic/9960040/8f2fd590-95d0-40ff-8e26-ba1b2224cc33-image.png",
	// 			ItemUUID:  "be2151dc-cbff-4be6-ab1b-afa02ea772f1",
	// 			Title:     request.Title,
	// 			Lyric:     "lyric",
	// 			Tags:      "uplifting, danceable, pop",
	// 		},
	// 	},
	// }

	// usecaseの呼び出し
	response, err := mc.mu.CreateMusic(request.Prompt, request.Lyrics, request.Title, request.IsAuto, request.Instrumental, uint(dID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response)
}
