package router

import (
	"net/http"
	"os"

	"github.com/kenta-kenta/diary-music/controller"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, dc controller.IDiaryController, mc controller.IMusicController) *echo.Echo {
	e := echo.New()
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},                                                                              // 許可するオリジン
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken}, // 許可するヘッダー
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},                                                                                            // 許可するメソッド
		AllowCredentials: true,                                                                                                                                // クレデンシャル情報（Cookieなど）の送信を許可
	}))
	// CSRF
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",                     // Cookieのパス
		CookieDomain:   os.Getenv("API_DOMAIN"), // Cookieのドメイン
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieSameSite: http.SameSiteNoneMode,
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.Login)
	e.POST("/logout", uc.Logout)
	e.GET("/csrf", uc.CsrfToken)

	auth := e.Group("")
	auth.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	auth.GET("/user", uc.GetUser)

	diaries := auth.Group("/diaries")
	diaries.GET("", dc.GetAllDiaries) // クエリパラメータが必要(?page=1&page_size=10)
	diaries.GET("/:diaryId", dc.GetDiaryById)
	diaries.GET("/dates", dc.GetDiaryDates) // クエリパラメータが必要(?year=2021&month=1)
	diaries.POST("", dc.CreateDiary)
	diaries.PUT("/:diaryId", dc.UpdateDiary)
	diaries.DELETE("/:diaryId", dc.DeleteDiary)

	musics := auth.Group("/musics")
	musics.GET("", mc.GetMusicsList) // クエリパラメータが必要(?page=1&limit=10)

	// diaries.POST("/:diaryId/musics", mc.CreateMusic)

	return e
}
