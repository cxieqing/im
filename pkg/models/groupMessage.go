package models

import (
	"errors"
	"im/pkg/db"
	"time"

	"gorm.io/gorm"
)

type GroupMessage struct {
	ID          uint `gorm:"primarykey"`
	ContentType ContenType
	Content     string
	From        uint
	To          uint
	Len         float32
	CreatedAt   time.Time
}

func (g *GroupMessage) Create() (uint, error) {
	result := db.MysqlDB.Create(g)
	return g.ID, result.Error
}

func (g *GroupMessage) List(page, pageSize int) []GroupMessage {
	var messages []GroupMessage
	result := db.MysqlDB.Where("to=? AND ", g.To).Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return messages
}
