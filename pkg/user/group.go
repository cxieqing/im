package user

type Group struct {
	Id         int
	GroupOwner int
	Members    []int
}

func (g *Group) Creat() bool {
	return true
}

func (g *Group) Out(uid int) bool {
	return true
}

func GetGroupsByUserId(id int) []Group {
	return make([]Group, 1)
}
