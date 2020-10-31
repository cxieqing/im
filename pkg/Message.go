package pkg

import (
	"encoding/json"
	"im/pkg/models"
)

type Message struct {
	ID          uint               `json:"id"`
	Content     string             `json:"content"`
	ContentType models.ContenType  `json:"contentType"`
	To          uint               `json:"to"`
	From        uint               `json:"from"`
	FromName    string             `json:"fromName"`
	FromAvatar  string             `json:"fromAvatar"`
	GroupID     uint               `json:"groupId"`
	GroupName   string             `json:"groupName"`
	MessageType models.MessageType `json:"messageType"`
	Len         float32            `json:"len"`
	IsRead      models.ReadType    `json:"isRead"`
	CreateAt    int64              `json:"createAt"`
}

func (m *Message) Complete() {
	fromUser := models.User{}
	fromUser.ID = m.From
	err := fromUser.GetOne()
	if err == nil {
		m.FromAvatar = fromUser.Avatar
		m.FromName = fromUser.UserName
	}
	if m.MessageType == models.GroupType {
		group := models.Group{}
		group.ID = m.To
		err := group.GetOne()
		if err != nil {
			m.GroupName = group.Name
			m.GroupID = m.To
		}
	}
}

func MsgDispatch(msg interface{}, s *ImServer) {
	m, err := clientMsgToMsg(msg)
	if err != nil {
		s.Log.Info("消息解析失败", err.Error())
		return
	}
	m.Complete()
	if m.MessageType == models.GroupType {
		msg := MsgToGroupMsg(m)
		_, err := msg.Create()
		if err != nil {
			return
		}
		groupBucket, err := GlobalGroupMap.Group(m.To)
		if err != nil {
			groupBucket.SendMessage(&m, func(mid uint) {
				msg := models.UserGroupMessage{
					ToUserID:       mid,
					GroupMessageID: msg.ID,
					IsSend:         models.UnSend,
					IsRead:         models.UnRead,
				}
				msg.Create()
			})
		}
	} else {
		msg := MsgToUserMsg(m)
		if client := GetOnlineUserById(msg.To, s); client != nil {
			_, err := msg.CreateSendMsg()
			if err != nil {
				return
			}
			client.MessageChn <- m
		} else {
			msg.CreateUnSendMsg()
		}
	}
}

func UserMsgToMsg(msg models.UserMessage) Message {
	data := Message{
		ID:          msg.ID,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		To:          msg.To,
		From:        msg.From,
		Len:         msg.Len,
		IsRead:      msg.IsRead,
		MessageType: models.UserType,
		CreateAt:    msg.Model.CreatedAt,
	}
	return data
}

func MsgToUserMsg(msg Message) models.UserMessage {
	data := models.UserMessage{
		Content:     msg.Content,
		ContentType: msg.ContentType,
		From:        msg.From,
		To:          msg.To,
		Len:         msg.Len,
	}
	return data
}

func GroupMsgToMsg(msg models.GroupMessage) Message {
	data := Message{
		ID:          msg.ID,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		To:          msg.To,
		Len:         msg.Len,
		MessageType: models.GroupType,
		IsRead:      models.UnRead,
		CreateAt:    msg.CreatedAt,
	}
	return data
}

func MsgToGroupMsg(msg Message) models.GroupMessage {
	data := models.GroupMessage{
		Content:     msg.Content,
		ContentType: msg.ContentType,
		From:        msg.From,
		To:          msg.To,
		Len:         msg.Len,
	}
	return data
}

func UserGroupMsgToMsg(msg models.UserGroupMessage) Message {
	data := Message{
		ID:          msg.GroupMessage.ID,
		Content:     msg.GroupMessage.Content,
		ContentType: msg.GroupMessage.ContentType,
		To:          msg.ToUserID,
		GroupID:     msg.GroupMessage.To,
		Len:         msg.GroupMessage.Len,
		MessageType: models.GroupType,
		IsRead:      msg.IsRead,
	}
	return data
}

func clientMsgToMsg(msg interface{}) (Message, error) {
	var m Message
	resByre, resByteErr := json.Marshal(msg)
	if resByteErr != nil {
		return m, resByteErr
	}
	jsonRes := json.Unmarshal(resByre, &m)
	if jsonRes != nil {
		return m, jsonRes
	}
	return m, nil
}
