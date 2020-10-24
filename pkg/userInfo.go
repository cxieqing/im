package pkg

import (
	"encoding/json"
	"im/pkg/models"
	"im/pkg/redis"
)

type UserInfo struct {
	User   models.User
	Groups []*GroupBucket
}

func NewUserInfo(uid uint) (*UserInfo, error) {
	user := models.User{}
	user.ID = uid
	if err := user.GetOne(); err != nil {
		return nil, err
	}
	return &UserInfo{User: user, Groups: make([]*GroupBucket, 0)}, nil
}

func (u *UserInfo) InitGroup(c *Client) {
	groups := models.GetGroupsByUserId(u.User.ID)
	var (
		gb  *GroupBucket
		err error
	)
	for _, g := range groups {
		gb, err = GlobalGroupMap.Group(g.ID)
		if err != nil {
			gb.AddClient(c)
			u.Groups = append(u.Groups, gb)
		}
	}
}

func CheckUserToken(token string) *UserInfo {
	redis := redis.NewRedis()
	cacheKey := "login_user_" + token
	val, err := redis.Client.Get(cacheKey).Result()
	if err != nil {
		return nil
	}
	var user = UserInfo{}
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil
	}
	return &user
}

func UserTokenClear(c *Client) {
	redis := redis.NewRedis()
	redis.Client.Do("delete", c.Token)
}
