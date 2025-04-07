package todos

import (
	"time"
)

type Todo struct {
	ID        uint   `gorm:"primarykey"`
	Title     string `json:"text" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Todo) TableName() string {
	return "todos"
}
