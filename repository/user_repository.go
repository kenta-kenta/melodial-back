package repository

import (
	"github.com/kenta-kenta/diary-music/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
	GetUserById(user *model.User, userId uint) error
}
type UserRepository struct {
	db *gorm.DB
}

// コンストラクタ
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// ユーザー情報の取得
func (ur *UserRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email = ?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

// ユーザー情報の作成
func (ur *UserRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserById(user *model.User, userId uint) error {
	if err := ur.db.First(user, userId).Error; err != nil {
		return err
	}
	return nil
}
