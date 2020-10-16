package pkg

import "im/pkg/models"

type Message struct {
	ID          uint               `json:"id"`
	Content     string             `json:"content"`
	ContentType models.ContenType  `json:"contentType"`
	To          uint               `json:"to"`
	From        uint               `json:"from"`
	GroupID     uint               `json:"groupId"`
	MessageType models.MessageType `json:"messageType"`
	Len         float32            `json:"len"`
	IsRead      models.ReadType    `json:"isRead"`
	CreateAt    int                `json:"createAt"`
}

func (m *Message) Dispatch(s *ImServer) {
	if m.MessageType == models.GroupType {
		msg := MsgToGroupMsg(*m)
		_, err := msg.Create()
		if err != nil {
			return
		}
		groupBucket, err := GlobalGroupMap.Group(m.To)
		if err != nil {
			groupBucket.SendMessage(m, func(mid uint) {
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
		msg := MsgToUserMsg(*m)
		if client := GetOnlineUserById(msg.To, s); client != nil {
			_, err := msg.CreateSendMsg()
			if err != nil {
				return
			}
			client.MessageChn <- *m
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
		Len:         msg.Len,
		IsRead:      msg.IsRead,
		MessageType: models.UserType,
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
