package message

import (
	"errors"
	"im/pkg/db"
	"time"

	"gorm.io/gorm"
)

type MessageType uint8

const (
	User MessageType = 1

	Group MessageType = 2
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

type Message struct {
	ID          uint `gorm:"primarykey"`
	ContentType ContenType
	Content     string
	From        uint
	To          uint
	IsSend      SendType
	Len         float32
	CreatedAt   time.Time
}

type UserMessage struct {
	Message
	MessageType MessageType
	IsRead      ReadType
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type GroupMessage struct {
	Message
}

func GetUserUnsendMsgByUserId(uid uint) []Message {
	var messages []Message
	result := db.MysqlDB.Where("message_type=? AND is_send=? AND to=?", User, UnSend, uid).Find(messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return messages
}

func GetUserUnsendMsgByGroupId(gid uint, uid uint) []Message {
	var messages []Message
	result := db.MysqlDB.Where("message_type=? AND is_send=? AND from=? AND to=?", Group, UnSend, gid, uid).Find(messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return messages
}

func SaveUnsendGroupMsg(m Message) (uint, error) {
	now := time.Now()
	umsg := UserMessage{
		Message: m,
	}
	umsg.CreatedAt = now
	umsg.IsSend = UnSend
	umsg.IsRead = UnRead
	umsg.MessageType = Group
	result := db.MysqlDB.Save(umsg)
	return umsg.ID, result.Error
}

func SaveUserMsg(m Message, s SendType) (uint, error) {
	now := time.Now()
	umsg := UserMessage{
		Message: m,
	}
	umsg.CreatedAt = now
	umsg.IsSend = s
	umsg.IsRead = UnRead
	umsg.MessageType = User
	result := db.MysqlDB.Save(umsg)
	return umsg.ID, result.Error
}

func SaveGroupMsg(m Message) (uint, error) {
	now := time.Now()
	gmsg := GroupMessage{
		Message: m,
	}
	gmsg.CreatedAt = now
	gmsg.IsSend = HasSend
	result := db.MysqlDB.Save(gmsg)
	return gmsg.ID, result.Error
}
