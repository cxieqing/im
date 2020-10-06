package user

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupOwner int
	Members    []uint
}

func (g *Group) Creat() bool {
	return true
}

func (g *Group) Out(uid uint) bool {
	return true
}

func GetGroupsByUserId(id uint) []Group {
	return make([]Group, 1)
}
