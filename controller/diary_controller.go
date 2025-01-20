package controller

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/usecase"
	"github.com/labstack/echo/v4"
)

type IDiaryController interface {
	GetAllDiaries(c echo.Context) error
	GetDiaryById(c echo.Context) error
	CreateDiary(c echo.Context) error
	UpdateDiary(c echo.Context) error
	DeleteDiary(c echo.Context) error
	GetDiaryDates(c echo.Context) error
}

type diaryController struct {
	du usecase.IDiaryUsecase
}

func NewDiaryController(du usecase.IDiaryUsecase) IDiaryController {
	return &diaryController{du}
}

func (dc *diaryController) GetAllDiaries(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))
	// Pageを取得
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	// PageSizeを取得
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	// GetAllDiariesメソッドを呼び出し
	response, err := dc.du.GetAllDiaries(userId, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (dc *diaryController) GetDiaryById(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	id := c.Param("diaryId")
	diaryId, _ := strconv.Atoi(id)
	diaries, err := dc.du.GetDiaryById(uint(userId.(float64)), uint(diaryId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, diaries)
}

func (dc *diaryController) GetDiaryDates(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	year := c.QueryParam("year")
	month := c.QueryParam("month")

	dates, err := dc.du.GetDiaryDates(userId, year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dates)
}

func (dc *diaryController) CreateDiary(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	diary := model.Diary{}
	if err := c.Bind(&diary); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	diary.UserId = uint(userId.(float64))

	// diaryRes := model.MusicResponse{
	// 	Status:  0,
	// 	Message: "",
	// 	Data: []model.MusicData{
	// 		{
	// 			AudioFile: "https://files.topmediai.com/aimusic/api/e630a4ff-8dae-493e-9e72-d8a3a4cc29f2-audio.mp3",
	// 			ImageFile: "https://files.topmediai.com/aimusic/8941427/32527de8-a833-431c-a56f-24e193bb4c65-image.png",
	// 			ItemUUID:  "e630a4ff-8dae-493e-9e72-d8a3a4cc29f2",
	// 			Title:     "晴れのキャンパス",
	// 			Lyric:     "[Verse]\n今日は大学に行った\n青空が広がった\n自転車こいで\n風が気持ちいいね\n\n[Verse 2]\n教室はまぶしい笑顔\n仲間たちとジョーク交わす\nアスファルトの道\n皆で一緒に歩く\n\n[Chorus]\n晴れのキャンパスが呼んでる\n夢を追いかけよう 皆と共に\n輝く未来が待っている\n一歩一歩 強くなれる\n\n[Verse 3]\n図書館で本を読んだ\n新しい知恵を得た\nカフェで一息\nコーヒーの香り\n\n[Bridge]\n夕焼けがきれいで\n寂しさも吹き飛ぶ\n星が瞬く夜\nまた明日が楽しみ\n\n[Chorus]\n晴れのキャンパスが呼んでる\n夢を追いかけよう 皆と共に\n輝く未来が待っている\n一歩一歩 強くなれる",
	// 			Tags:      "",
	// 		},
	// 	},
	// }

	diaryRes, err := dc.du.CreateDiaryWithMusic(&diary)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, diaryRes)
}

func (dc *diaryController) UpdateDiary(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	id := c.Param("diaryId")
	taskId, _ := strconv.Atoi(id)
	diary := model.Diary{}
	if err := c.Bind(&diary); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	diaryRes, err := dc.du.UpdateDiary(uint(userId.(float64)), uint(taskId), diary)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, diaryRes)
}

func (dc *diaryController) DeleteDiary(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	id := c.Param("diaryId")
	taskId, _ := strconv.Atoi(id)
	err := dc.du.DeleteDiary(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
