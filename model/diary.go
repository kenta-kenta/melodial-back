package model

import "time"

type Diary struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id"`
	Content   string    `json:"content"`
	Music     []Music   `json:"music" gorm:"foreignKey:DiaryID"` // 一対一の関係
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DiaryResponse struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Content   string      `json:"content" gorm:"not null"`
	MusicData []MusicData `json:"music_data" gorm:"foreignKey:DiaryID"` // 一対一の関係
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type DiaryDate struct {
	Date time.Time `json:"date"`
}

type DiaryDateResponse struct {
	Dates []time.Time `json:"dates"`
}

type DiaryDateCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DiaryDateCountResponse struct {
	Dates []DiaryDateCount `json:"dates"`
}

// PaginationQueryはページネーションのクエリを表します。
type PaginationQuery struct {
	Page     int
	PageSize int
}

type PaginationResponse struct {
	DiaryResponse []DiaryResponse `json:"diary_response"`
	TotalItems    int64           `json:"total_items"`
	Page          int             `json:"page"`
	PageSize      int             `json:"page_size"`
	TotalPages    int             `json:"total_pages"`
}
