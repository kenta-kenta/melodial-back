package main

import (
	"fmt"

	"github.com/kenta-kenta/diary-music/db"
	"github.com/kenta-kenta/diary-music/model"
)

func main() {
	// DB接続
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	// マイグレーション
	dbConn.AutoMigrate(&model.User{}, &model.Diary{})
}
