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
	diaryRes, err := dc.du.CreateDiary(diary)
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
