package rpc

import (
	"im/pkg"
	"im/pkg/models"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type LoginParams struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type TokenParams struct {
	Token string `json:"token"`
}

type ResponseMsg struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewSuccessMsg(data interface{}) ResponseMsg {
	return ResponseMsg{
		Status: 1,
		Msg:    "success",
		Data:   data,
	}
}

func UserRpc(w http.ResponseWriter, r *http.Request) {
	rpc.RegisterName("ImService", new(ImService))
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "OPTIONS,POST,GET")
	header.Add("Access-Control-Allow-Headers", "Content-type")
	header.Add("Access-Control-Allow-Origin", "*")
	var conn io.ReadWriteCloser = struct {
		io.Writer
		io.ReadCloser
	}{
		ReadCloser: r.Body,
		Writer:     w,
	}

	rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
}

type ImService struct {
}

func (i *ImService) Login(request LoginParams, reply *ResponseMsg) error {
	user := models.User{UserName: request.Username, Password: request.Password}
	if err := user.CheckUser(); err != nil {
		token := pkg.UserHashToken()
	}
	*reply = NewSuccessMsg("hello word")
	return nil
}

func (i *ImService) UserInit(request TokenParams, reply *ResponseMsg) {

}
