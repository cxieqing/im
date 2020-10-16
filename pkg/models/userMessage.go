package models

import (
	"errors"
	"im/pkg/db"

	"gorm.io/gorm"
)

type MessageType uint8

const (
	UserType MessageType = 1

	GroupType MessageType = 2
)

type ReadType uint8

const (
	HasRead ReadType = 1

	UnRead ReadType = 0
)

type SendType uint8

const (
	HasSend SendType = 1

	UnSend SendType = 0
)

type ContenType uint8

const (
	ImageContent ContenType = iota
	TextContent
	VdioContent
)

type UserMessage struct {
	ContentType ContenType
	Content     string
	From        uint
	To          uint
	IsSend      SendType
	Len         float32
	IsRead      ReadType
	Model
}

func (m *UserMessage) create() (uint, error) {
	result := db.MysqlDB.Create(m)
	return m.ID, result.Error
}

func (m *UserMessage) CreateSendMsg() (uint, error) {
	m.IsSend = HasSend
	return m.create()
}

func (m *UserMessage) CreateUnSendMsg() (uint, error) {
	m.IsRead = UnRead
	m.IsSend = UnSend
	return m.create()
}

func (m *UserMessage) UnSendMsgList() []UserMessage {
	var messages []UserMessage
	result := db.MysqlDB.Where("is_send=? AND `to`=?", UnSend, m.To).Find(&messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return messages
}

func (m *UserMessage) HistoryMsgList(page, pageSize int) ([]UserMessage, int64) {
	var messages []UserMessage
	result := db.MysqlDB.Where("is_send=? AND ((`to`=? AND `from`=?) or (`to`=? AND `from`=?))", HasSend, m.To, m.From, m.From, m.To).Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	var total int64
	db.MysqlDB.Model(&UserMessage{}).Where("is_send=? AND ((`to`=? AND `from`=?) or (`to`=? AND `from`=?))", HasSend, m.To, m.From, m.From, m.To).Count(&total)
	return messages, total
}

func (m *UserMessage) delete() error {
	return db.MysqlDB.Delete(m).Error
}
