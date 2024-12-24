package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kenta-kenta/diary-music/model"
	"github.com/kenta-kenta/diary-music/usecase"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	CsrfToken(c echo.Context) error
	GetUser(c echo.Context) error
}

type UserController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &UserController{uu}
}

func (uc *UserController) SignUp(c echo.Context) error {
	user := model.User{}
	// リクエストボディのバインド
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// ユーザー登録
	resUser, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, resUser)
}

func (uc *UserController) Login(c echo.Context) error {
	user := model.User{}
	// リクエストボディのバインド
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// ログイン
	token, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)                      // Cookieの生成
	cookie.Name = "token"                           // Cookie名
	cookie.Value = token                            // Cookie値
	cookie.Expires = time.Now().Add(24 * time.Hour) // 有効期限
	cookie.Path = "/"                               // パス
	cookie.Domain = os.Getenv("API_DOMAIN")         // ドメイン
	cookie.Secure = true
	cookie.HttpOnly = true                  // JavaScriptからのアクセスを禁止
	cookie.SameSite = http.SameSiteNoneMode // SameSite属性
	c.SetCookie(cookie)                     // Cookieの設定
	return c.NoContent(http.StatusOK)
}

func (uc *UserController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)              // Cookieの生成
	cookie.Name = "token"                   // Cookie名
	cookie.Value = ""                       // Cookie値
	cookie.Expires = time.Now()             // 有効期限
	cookie.Path = "/"                       // パス
	cookie.Domain = os.Getenv("API_DOMAIN") // ドメイン
	cookie.Secure = true
	cookie.HttpOnly = true                  // JavaScriptからのアクセスを禁止
	cookie.SameSite = http.SameSiteNoneMode // SameSite属性
	c.SetCookie(cookie)                     // Cookieの設定
	return c.NoContent(http.StatusOK)
}

// CsrfToken returns a CSRF token
func (uc *UserController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{"csrf_token": token})
}

func (uc *UserController) GetUser(c echo.Context) error {
	// JWTトークンからユーザーID取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	// DBからユーザー情報取得
	userInfo := model.User{}
	if err := uc.uu.GetUserById(&userInfo, userId); err != nil {
		return c.JSON(http.StatusNotFound, "ユーザーが見つかりません")
	}

	// レスポンス用の構造体
	response := model.UserResponse{
		ID:       userInfo.ID,
		Email:    userInfo.Email,
		UserName: userInfo.UserName,
	}

	return c.JSON(http.StatusOK, response)
}
