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

	// createdAt, err := time.Parse(time.RFC3339, "2025-01-24T14:22:35.282635+09:00")
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// updatedAt, err := time.Parse(time.RFC3339, "2025-01-24T14:22:35.282635+09:00")
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	// diaryRes := model.DiaryResponse{
	// 	ID:      9,
	// 	Content: "昔、飼っていた猫のことを思い出した。とても甘えん坊で、いつも私のそばにいた。一緒に遊んだり、寝たりするのが大好きだった。ある日、突然いなくなってしまい、とても悲しかった。あの時の気持ちを、今でも鮮明に覚えている。あの猫との出会いは、私にとってかけがえのない宝物だ。",
	// 	MusicData: []model.MusicData{
	// 		{
	// 			AudioFile: "https://files.topmediai.com/aimusic/api/7f4f8313-4bc5-49b2-9ecc-751006f693f4-audio.mp3",
	// 			ImageFile: "https://files.topmediai.com/aimusic/9441627/1d5194f1-c210-4c92-a8f5-edfbb6bcee9a-image.png",
	// 			ItemUUID:  "7f4f8313-4bc5-49b2-9ecc-751006f693f4",
	// 			Title:     "かけがえのない宝物",
	// 			Lyric:     "[Verse]\n小さな足跡 追いかけていた\n甘えん坊の君 そばにいた\n一緒に遊んで 寝転んでいた\nあの日々は煌めく宝物\n\n[Verse 2]\n日の光浴びて 輝く瞳\n君の温もり 忘れられない\nどこに行ったのか 一人彷徨う\n戻ってくる日を ずっと待つよ\n\n[Chorus]\n君がいなくなった あの日の泣き声\n胸に響いて 今も忘れない\n君との出会いは かけがえのない\n宝物 心にずっと在る\n\n[Verse 3]\n朝日の中で 君と笑った\n柔らかな毛皮 触れた時の\nぬくもり覚えて 胸が痛むよ\n君の名前を 何度も呼んだ\n\n[Bridge]\n月夜の下で 夢を見てる\n君が帰る日を 待っているよ\n星空見上げて 君を感じる\nあの時の気持ち ここに在る\n\n[Chorus]\n君がいなくなった あの日の泣き声\n胸に響いて 今も忘れない\n君との出会いは かけがえのない\n宝物 心にずっと在る",
	// 			Tags:      "",
	// 		},
	// 	},
	// 	CreatedAt: createdAt,
	// 	UpdatedAt: updatedAt,
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
