package user

type User struct {
	Id       int
	UserName string
	Password string
	NikeName string
	Icon     string
}

func Add(u User) bool {
	return true
}

func Update(u User) bool {
	return true
}
