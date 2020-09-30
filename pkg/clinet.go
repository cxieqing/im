package pkg

import (
	"encoding/json"
	"im/pkg/message"
	"im/pkg/redis"
	"im/pkg/user"

	"github.com/gorilla/websocket"
)

type Clinet struct {
	Connet     *websocket.Conn
	Group      map[int]*user.Group
	UserInfo   *user.User
	MessageChn chan ImMessage
}

type ImMessage struct {
	Token       string              `json:"token"`
	Content     string              `json:"content"`
	ContentType message.ContenType  `json:"contentType"`
	To          int                 `json:"to"`
	From        int                 `json:"from"`
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
	go func() {
		c.InitMessage()
	}()
}

func (c *Clinet) InitGroup() {
	groups := user.GetGroupsByUserId(c.UserInfo.Id)
	for _, g := range groups {
		c.Group[g.Id] = &g
	}
}

func (c *Clinet) InitMessage() {
	messages := message.GetUserUnsendMsgByUserId(c.UserInfo.Id)
	for _, m := range messages {
		c.MessageChn <- MsgTransToIm(m, message.UserMessage)
	}
}

func (c *Clinet) InitGroupMessage() {
	for gid := range c.Group {
		messages := message.GetUserUnsendMsgByGroupId(gid)
		for _, m := range messages {
			c.MessageChn <- MsgTransToIm(m, message.GroupMessage)
		}
	}
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
