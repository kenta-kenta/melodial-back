package main

import (
	"github.com/kenta-kenta/diary-music/controller"
	"github.com/kenta-kenta/diary-music/db"
	"github.com/kenta-kenta/diary-music/repository"
	"github.com/kenta-kenta/diary-music/router"
	"github.com/kenta-kenta/diary-music/usecase"
	"github.com/kenta-kenta/diary-music/validator"
)

func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	diaryValidator := validator.NewDiaryValidator()
	userRepository := repository.NewUserRepository(db)
	diaryRepository := repository.NewDiaryRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	diaryUsecase := usecase.NewDiaryUsecase(diaryRepository, diaryValidator)
	userController := controller.NewUserController(userUsecase)
	diaryController := controller.NewDiaryController(diaryUsecase)
	e := router.NewRouter(userController, diaryController)
	e.Logger.Fatal(e.Start(":8080"))
}
