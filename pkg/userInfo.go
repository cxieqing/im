package pkg

import (
	"fmt"
	"im/pkg/models"
	"im/pkg/tools"
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

func UserHashToken(uid uint) string {
	return tools.Md5("user" + fmt.Sprint("%d", uid))
}
