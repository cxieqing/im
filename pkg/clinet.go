package pkg

import (
	"im/pkg/message"
	"im/pkg/user"

	"github.com/gorilla/websocket"
)

type Clinet struct {
	connet   *websocket.Conn
	Group    map[int]*user.Group
	UserInfo *user.User
	Message  chan message.Message
}

func NewClinet(connet *websocket.Conn, token string) *Clinet {
	user := checkToken(token)
	if user == nil {
		return nil
	}
	return &Clinet{connet: connet}
}

func (c *Clinet) initGroup() {
	groups := user.GetGroupsByUserId(c.UserInfo.Id)
	for _, g := range groups {
		c.Group[g.Id] = &g
	}
}

func (c *Clinet) initMessage(sc chan<- message.Message) {
	messages := message.GetUserUnReadMsgById(c.UserInfo.Id)
	for _, m := range messages {
		sc <- m
	}
}

func checkToken(token string) *user.User {
	redis := NewRedis()

}
