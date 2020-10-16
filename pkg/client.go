package pkg

import (
	"encoding/json"
	"fmt"
	"im/pkg/models"
	"im/pkg/redis"
	"time"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 512

var pingPeriod = 5 * time.Second

type MsgType uint8

const (
	ReadConnect MsgType = iota
	Normal
)

type ClinetMsg struct {
	MsgType MsgType     `json:"msgType"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
}

type Client struct {
	Token      string
	Connet     *websocket.Conn
	UserInfo   *UserInfo
	MessageChn chan Message
}

func NewClient(connet *websocket.Conn, token string) *Client {
	userPtr := checkToken(token)
	if userPtr == nil {
		return nil
	}
	UserInfoPtr, err := NewUserInfo(userPtr.User.ID)
	if err != nil {
		return nil
	}
	return &Client{Connet: connet, UserInfo: UserInfoPtr}
}

func (c *Client) Init() {
	c.UserInfo.InitGroup(c)
	go c.InitMessage()
	go c.InitGroupMessage()
}

func (c *Client) InitMessage() {
	um := models.UserMessage{To: c.UserInfo.User.ID}
	messages := um.UnSendMsgList()
	for _, m := range messages {
		c.ReceiveMessage(UserMsgToMsg(m))
	}
}

func (c *Client) InitGroupMessage() {
	ugm := models.UserGroupMessage{ToUserID: c.UserInfo.User.ID}
	messages := ugm.UnSendUserGroupMsgList()
	for _, msg := range messages {
		c.ReceiveMessage(UserGroupMsgToMsg(msg))
	}
}

func (c *Client) ReadMessage() (Message, error) {
	msg := Message{}
	c.Connet.SetReadLimit(maxMessageSize)
	err := c.Connet.ReadJSON(msg)
	return msg, err
}

func (c *Client) ReceiveMessage(msg Message) {

}

func (c *Client) HandReadMessage(s *ImServer) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			//@todo log
		} else {
			msg.Dispatch(s)
		}
	}
}

func (c *Client) HandWriteMessage(s *ImServer) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.WaitCloseClient <- c.Token
	}()
	writeDeadline := time.Duration(s.WriteDeadline)
	for {
		select {
		case msg, ok := <-c.MessageChn:
			if !ok {
				c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Connet.SetWriteDeadline(time.Now().Add(writeDeadline))
			c.Connet.WriteJSON(msg)
		case <-ticker.C:
			err := c.Connet.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(writeDeadline))
			if err != nil {
				c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		}
	}
}

func (c *Client) Free() {
	c.Connet.Close()
	close(c.MessageChn)
	c.UserInfo = nil
	GlobalGroupMap.ClientFree(c)
	tokenClear(c)
}

func GetOnlineUserById(id uint, s *ImServer) *Client {
	token, err := getUserTokenById(id)
	if err != nil {
		return nil
	}
	clientPtr, ok := s.Clients[token]
	if !ok {
		return nil
	}
	return clientPtr
}

func getUserTokenById(id uint) (string, error) {
	cacheKey := "userid:token:" + fmt.Sprint(id)
	redis := redis.NewRedis()
	return redis.Client.Get(cacheKey).Result()
}

func checkToken(token string) *UserInfo {
	redis := redis.NewRedis()
	val, err := redis.Client.Get(token).Result()
	if err != nil {
		return nil
	}
	var user = UserInfo{}
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil
	}
	return &user
}

func tokenClear(c *Client) {
	cacheKey := "userid:token:" + fmt.Sprint(c.UserInfo.User.ID)
	redis := redis.NewRedis()
	redis.Client.Do("delete", cacheKey)
	redis.Client.Do("delete", c.Token)
}
