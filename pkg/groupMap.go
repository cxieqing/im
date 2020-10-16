package pkg

import (
	"im/pkg/models"
	"im/pkg/tools"
	"strconv"
	"sync"
)

var mu sync.Mutex

type GroupMap map[string]*GroupBucket

var GlobalGroupMap = make(GroupMap)

func (gm *GroupMap) hashToken(gid uint) string {
	return tools.Md5(strconv.FormatUint(uint64(gid), 10))
}

func (gm *GroupMap) Group(gid uint) (*GroupBucket, error) {
	token := gm.hashToken(gid)
	mu.Lock()
	defer mu.Unlock()
	if gb, ok := (*gm)[token]; ok {
		return gb, nil
	}
	var err error
	(*gm)[token], err = NewGroupBucket(gid)
	if err != nil {
		return nil, err
	}
	return (*gm)[token], nil
}

func (gm *GroupMap) FreeGroup(gid uint) {
	token := gm.hashToken(gid)
	mu.Lock()
	defer mu.Unlock()
	if _, ok := (*gm)[token]; ok {
		delete((*gm), token)
	}
}

func (gm *GroupMap) ClientFree(c *Client) {
	for _, gb := range c.UserInfo.Groups {
		for k, mc := range gb.membersClient {
			if mc.Token == c.Token {
				gb.membersClient = append(gb.membersClient[:k], gb.membersClient[k+1:]...)
			}
		}
	}
}

type GroupBucket struct {
	mu            sync.Mutex
	group         *models.Group
	membersClient []*Client
}

func NewGroupBucket(gid uint) (*GroupBucket, error) {
	group := models.Group{}
	group.ID = gid
	if err := group.GetOne(); err != nil {
		return nil, err
	}
	return &GroupBucket{group: &group, membersClient: make([]*Client, 0)}, nil
}

func NewGroupBucketWithNewGroup(owner uint, mids ...uint) (*GroupBucket, error) {
	var (
		err error
		gid uint
		gb  *GroupBucket
	)
	group := models.Group{GroupOwner: owner}
	gid, err = group.Create(mids...)
	if err != nil {
		return nil, err
	}
	gb, err = GlobalGroupMap.Group(gid)
	if err != nil {
		return nil, err
	}
	return gb, nil
}

func (g *GroupBucket) AddClient(c *Client) {
	g.mu.Lock()
	g.membersClient = append(g.membersClient, c)
	g.mu.Unlock()
}

func (g *GroupBucket) SendMessage(m *Message, unOnlineUserMessageSave func(mid uint)) {
	hasSendMembers := make([]uint, 0, len(g.membersClient))
	for _, c := range g.membersClient {
		if m.From != c.UserInfo.User.ID {
			c.MessageChn <- *m
		}
		hasSendMembers = append(hasSendMembers, c.UserInfo.User.ID)
	}
	unOnLineUserIDs := tools.ArrayDiffUint(g.group.Members, hasSendMembers)
	for _, v := range unOnLineUserIDs {
		unOnlineUserMessageSave(v)
	}
}

func (g *GroupBucket) AddMember(s *ImServer, mids ...uint) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.group.Members = append(g.group.Members, tools.ArrayDiffUint(mids, g.group.Members)...)
	err := g.group.AddMember(mids...)
	if err != nil {
		return err
	}
	for _, uid := range mids {
		if c := GetOnlineUserById(uid, s); c != nil {
			g.membersClient = append(g.membersClient, c)
		}
	}
	return nil
}

func (g *GroupBucket) DeleteMember(mid uint) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.group.OutMember(mid); err != nil {
		return err
	}
	if g.group.GroupOwner == mid {
		g.Free()
		return nil
	}
	for k, c := range g.membersClient {
		if c.UserInfo.User.ID == mid {
			g.membersClient = append(g.membersClient[:k], g.membersClient[k+1:]...)
		}
	}
	if len(g.membersClient) < 2 {
		g.Free()
	}
	return nil
}

func (g *GroupBucket) Free() {
	GlobalGroupMap.FreeGroup(g.group.ID)
	g.group = nil
	g.membersClient = nil

}
