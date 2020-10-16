package models

import "testing"

func TestUser(t *testing.T) {
	user := User{
		UserName: "nihao",
		NikeName: "你好",
		Password: "nihao",
		Icon:     "//www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png",
	}
	t.Run("User create", func(t *testing.T) {
		if id, _ := user.Create(); id < 1 {
			t.Errorf("unexpect new recode id %d", id)
		}
	})

	t.Run("User GetOne", func(t *testing.T) {
		u := User{}
		u.ID = user.ID
		if err := u.GetOne(); err != nil {
			t.Errorf("unexpect GetOne error  %s", err)
		}
	})

	t.Run("User update", func(t *testing.T) {
		user.NikeName = `你好你好`
		user.UserName = "nihaonihao"
		user.Icon = "https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=1984420730,272709403&fm=26&gp=0.jpg"
		rowsAffected, _ := user.Update()
		if rowsAffected != 1 {
			t.Errorf("want rowsAffected 1 got %d", rowsAffected)
		}
	})

	t.Run("User ResetPassword", func(t *testing.T) {
		if rowsAffected, _ := user.ResetPassword("nihaonihao"); rowsAffected != 1 {
			t.Errorf("want rowsAffected 1 got %d", rowsAffected)
		}
	})
}
