package models

import (
	"errors"
	"im/pkg/db"

	"gorm.io/gorm"
)

type UserGroupMessage struct {
	Model
	GroupMessageID uint
	ToUserID       uint
	IsSend         SendType
	IsRead         ReadType
	GroupMessage   GroupMessage
}

func (m *UserGroupMessage) Create() (uint, error) {
	result := db.MysqlDB.Create(m)
	return m.ID, result.Error
}

func (m *UserGroupMessage) Delete() error {
	return db.MysqlDB.Delete(m).Error
}

func (m *UserGroupMessage) UnSendUserGroupMsgList() []UserGroupMessage {
	var messages []UserGroupMessage
	result := db.MysqlDB.Where("is_send=? AND to_user_id=?", UnSend, m.ToUserID).Find(&messages)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return messages
}
