package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"im/pkg"
	"im/pkg/config"
	"im/pkg/models"
	"im/pkg/redis"
	"im/pkg/tools"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"regexp"
	"time"
)

type LoginParams struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegisterParams struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	NikeName string `json:"nikename"`
	Avatar   string `json:"avatar"`
}

type TokenParams struct {
	Token string `json:"token"`
}

type GroupUserListParams struct {
	TokenParams TokenParams
	GroupID     uint `json:"groupId"`
}

type ResponseData map[string]interface{}

type ResponseMsg struct {
	Status int          `json:"status"`
	Msg    string       `json:"msg"`
	Data   ResponseData `json:"data"`
}

func NewSuccessMsg(data ResponseData) ResponseMsg {
	return ResponseMsg{
		Status: 1,
		Msg:    "success",
		Data:   data,
	}
}

func NewErrorMsg(code int, err error) ResponseMsg {
	return ResponseMsg{
		Status: code,
		Msg:    err.Error(),
		Data:   ResponseData{},
	}
}

func CreateRpcServer() {
	rpc.RegisterName("userServer", new(UserServer))
	config := config.NewConfig()
	addr := fmt.Sprintf("%s:%d", config.RpcHost, config.RpcPort)
	s := &http.Server{
		Addr:           addr,
		Handler:        http.HandlerFunc(rpcHander),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func rpcHander(w http.ResponseWriter, r *http.Request) {
	CorsHeaderSet(w)
	var conn io.ReadWriteCloser = struct {
		io.Writer
		io.ReadCloser
	}{
		ReadCloser: r.Body,
		Writer:     w,
	}

	rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
}

func CorsHeaderSet(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type")
	header.Set("Access-Control-Allow-Origin", "*")
}

type UserServer struct {
}

func (i *UserServer) Login(request LoginParams, reply *ResponseMsg) error {
	user := models.User{UserName: request.Username, Password: request.Password}
	if err := user.CheckUser(); err == nil {
		userInfo := pkg.UserInfo{User: user}
		secondNum := 3600
		expire := time.Second * time.Duration(secondNum)
		token, err := mkLoginToken(userInfo, expire)
		if err != nil {
			*reply = NewErrorMsg(-1, err)
			return nil
		}
		*reply = NewSuccessMsg(ResponseData{"token": token, "expire": secondNum})
		return nil
	}
	*reply = NewErrorMsg(-1, errors.New("登录失败"))
	return nil
}

func (i *UserServer) refreshToken(request TokenParams, reply *ResponseMsg) error {
	userInfo := pkg.CheckUserToken(request.Token)
	if userInfo == nil {
		*reply = NewErrorMsg(-2, errors.New("用户信息不存在"))
		return nil
	}
	secondNum := 3600
	expire := time.Second * time.Duration(secondNum)
	token, err := mkLoginToken(*userInfo, expire)
	if err != nil {
		*reply = NewErrorMsg(-1, err)
		return nil
	}
	*reply = NewSuccessMsg(ResponseData{"token": token, "expire": secondNum})
	return nil
}

func mkLoginToken(u pkg.UserInfo, expire time.Duration) (string, error) {
	token := tools.Md5("user_" + fmt.Sprint("%d", u.User.ID) + "_" + fmt.Sprint("%d", time.Now().UnixNano()))
	cacheKey := "login_user_" + token
	data, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	redis := redis.NewRedis()
	_, err = redis.Client.Set(cacheKey, string(data), expire).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (i *UserServer) UserInit(request TokenParams, reply *ResponseMsg) error {
	userInfo := pkg.CheckUserToken(request.Token)
	if userInfo == nil {
		*reply = NewErrorMsg(-1, errors.New("用户信息不存在"))
		return nil
	}
	data := ResponseData{
		"user": ResponseData{
			"id":       userInfo.User.ID,
			"userName": userInfo.User.UserName,
			"avatar":   userInfo.User.Avatar,
			"mobile":   userInfo.User.Mobile,
			"nikeName": userInfo.User.NikeName,
		},
		"group": models.GetGroupsByUserId(userInfo.User.ID),
	}
	*reply = NewSuccessMsg(data)
	return nil
}

func (i *UserServer) GroupUserList(request GroupUserListParams, reply *ResponseMsg) error {
	userInfo := pkg.CheckUserToken(request.TokenParams.Token)
	if userInfo == nil {
		*reply = NewErrorMsg(-1, errors.New("用户信息不存在"))
		return nil
	}
	group := models.Group{}
	group.ID = request.GroupID
	if err := group.GetOne(); err != nil {
		*reply = NewErrorMsg(-1, err)
		return nil
	}
	inGroup := false
	for _, v := range group.Members {
		if v == userInfo.User.ID {
			inGroup = true
			break
		}
	}
	if !inGroup {
		*reply = NewErrorMsg(-1, errors.New("非法获取群信息"))
		return nil
	}
	userList := models.UserListByIDs(group.Members...)
	data := ResponseData{
		"group":    group,
		"userList": userList,
	}
	*reply = NewSuccessMsg(data)
	return nil
}

func (i *UserServer) UserRegister(request RegisterParams, reply *ResponseMsg) error {
	if err := validateRegisterParams(request); err != nil {
		*reply = NewErrorMsg(-1, err)
		return nil
	}
	user := models.User{
		UserName: request.Username,
		Password: request.Password,
		Mobile:   request.Mobile,
		NikeName: request.NikeName,
		Avatar:   request.Avatar,
	}
	if _, err := user.Create(); err != nil {
		*reply = NewErrorMsg(-1, err)
	}
	*reply = NewSuccessMsg(ResponseData{})
	return nil
}

func validateRegisterParams(p RegisterParams) error {
	if p.Username == "" || p.Password == "" || p.Mobile == "" || p.NikeName == "" || p.Avatar == "" {
		return errors.New("参数错误")
	}
	if len(p.Username) > 20 || len(p.Username) < 6 {
		return errors.New("用户名长度范围6~20")
	}
	r, _ := regexp.Compile(`1\d{10}`)
	if !r.MatchString(p.Mobile) {
		return errors.New("手机号码格式错误")
	}
	return nil
}
