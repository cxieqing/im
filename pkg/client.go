package pkg

import (
	"errors"
	"fmt"
	"im/pkg/models"
	"im/pkg/tools"
	"time"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 512

var pingPeriod = 5 * time.Second

type MsgType uint8

const (
	ReadyConnect MsgType = iota
	Normal
	CloseConnect
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
	PingNum    int
}

func NewClient(connet *websocket.Conn, token string) *Client {
	userPtr := CheckUserToken(token)
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
	c.Connet.SetCloseHandler(func(code int, text string) error {
		s := NewImServer()
		s.WaitCloseClient <- c.Token
		return nil
	})
	c.Connet.SetPongHandler(func(appData string) error {
		c.PingNum = 0
		return nil
	})
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

func (c *Client) ReadMessage() (ClinetMsg, error) {
	msg := ClinetMsg{}
	//c.Connet.SetReadLimit(maxMessageSize)
	err := c.Connet.ReadJSON(&msg)
	return msg, err
}

func (c *Client) ReceiveMessage(msg Message) {

}

func (c *Client) HandMessage(s *ImServer) {
	go c.HandReadMessage(s)
	go c.HandWriteMessage(s)
}

func (c *Client) HandReadMessage(s *ImServer) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				s.WaitCloseClient <- c.Token
				return
			}
		} else {
			if msg.MsgType == Normal {
				MsgDispatch(msg.Data, s)
			}
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
				//c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Connet.SetWriteDeadline(time.Now().Add(writeDeadline))
			c.Connet.WriteJSON(msg)
		case <-ticker.C:
			if c.Connet != nil {
				c.Connet.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(writeDeadline))
				if c.PingNum > 3 {
					c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				c.PingNum++
			} else {
				return
			}
		}
	}
}

func (c *Client) Free() {
	if c == nil {
		return
	}
	if c.Connet != nil {
		c.Connet.Close()
		c.Connet = nil
	}
	if c.MessageChn != nil {
		close(c.MessageChn)
	}
	GlobalGroupMap.ClientFree(c)
	c.UserInfo = nil
	UserTokenClear(c)
}

func WaitClientReady(conn *websocket.Conn, s *ImServer) error {
	if len(s.Clients) >= s.MaxClientNum {
		s.Log.Warn("服务端链接用户过多")
		return errors.New(" server crowd")
	}
	msg := ClinetMsg{}
	conn.SetReadLimit(maxMessageSize)
	if err := conn.ReadJSON(&msg); err != nil {
		s.Log.Info("获取信息失败：", err.Error())
		return err
	}
	if msg.MsgType == ReadyConnect {
		client := NewClient(conn, msg.Token)
		if client != nil {
			s.Log.Info("新上线用户：", client.UserInfo.User.ID)
			client.Token = ClientHashToken(client.UserInfo.User.ID)
			client.Init()
			s.NewClient <- client
		}
	}
	return nil
}

func GetOnlineUserById(id uint, s *ImServer) *Client {
	token := ClientHashToken(id)
	clientPtr, ok := s.Clients[token]
	if !ok {
		return nil
	}
	return clientPtr
}

func ClientHashToken(uid uint) string {
	return tools.Md5("user" + fmt.Sprint("%d", uid))
}
