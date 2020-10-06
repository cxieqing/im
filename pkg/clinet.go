package pkg

import (
	"encoding/json"
	"im/pkg/message"
	"im/pkg/redis"
	"im/pkg/user"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 512

type Clinet struct {
	Connet     *websocket.Conn
	Group      map[uint]*user.Group
	UserInfo   *user.User
	MessageChn chan ImMessage
}

type ImMessage struct {
	Token       string              `json:"token"`
	Content     string              `json:"content"`
	ContentType message.ContenType  `json:"contentType"`
	To          uint                `json:"to"`
	From        uint                `json:"from"`
	MessageType message.MessageType `json:"messageType"`
	Len         float32             `json:"len"`
}

func NewClinet(connet *websocket.Conn, token string) *Clinet {
	userPtr := checkToken(token)
	if userPtr == nil {
		return nil
	}
	return &Clinet{Connet: connet, UserInfo: userPtr}
}

func (c *Clinet) Init() {
	c.InitGroup()
	go c.InitMessage()
	go c.InitGroupMessage()
}

func (c *Clinet) InitGroup() {
	groups := user.GetGroupsByUserId(c.UserInfo.ID)
	for _, g := range groups {
		c.Group[g.ID] = &g
	}
}

func (c *Clinet) InitMessage() {
	messages := message.GetUserUnsendMsgByUserId(c.UserInfo.ID)
	for _, m := range messages {
		c.WriteMessage(MsgTransToIm(m, message.User))
	}
}

func (c *Clinet) InitGroupMessage() {
	for gid := range c.Group {
		messages := message.GetUserUnsendMsgByGroupId(gid, c.UserInfo.ID)
		for _, m := range messages {
			c.WriteMessage(MsgTransToIm(m, message.Group))
		}
	}
}

func (c *Clinet) ReadMessage() (ImMessage, error) {
	msg := ImMessage{}
	c.Connet.SetReadLimit(maxMessageSize)
	err := c.Connet.ReadJSON(msg)
	return msg, err
}

func (c *Clinet) WriteMessage(m ImMessage) {
	defer func() {
		if err := recover(); err != nil {
			dm := ImTransTomsg(m)
			dm.IsSend = message.UnSend
			if m.MessageType == message.User {
				message.SaveUserMsg(dm, message.UnSend)
			} else {
				message.SaveUnsendGroupMsg(dm)
			}
		}
	}()
	c.MessageChn <- m
}

func checkToken(token string) *user.User {
	redis := redis.NewRedis()
	val, err := redis.Client.Get(token).Result()
	if err != nil {
		return nil
	}
	var user = user.User{}
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil
	}
	return &user
}

func MsgTransToIm(msg message.Message, t message.MessageType) ImMessage {
	data := ImMessage{
		Content:     msg.Content,
		ContentType: msg.ContentType,
		To:          msg.To,
		Len:         msg.Len,
	}
	data.MessageType = t
	return data
}

func ImTransTomsg(msg ImMessage) message.Message {
	data := message.Message{
		Content:     msg.Content,
		ContentType: msg.ContentType,
		From:        msg.From,
		To:          msg.To,
		Len:         msg.Len,
	}
	return data
}
