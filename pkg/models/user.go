package models

import (
	"errors"
	"im/pkg/db"
	"im/pkg/tools"

	"gorm.io/gorm"
)

type User struct {
	UserName string
	Password string
	Mobile   string
	NikeName string
	Avatar   string
	Model
}

func (u *User) Create() (uint, error) {
	result := db.MysqlDB.Create(u)
	return u.ID, result.Error
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password = tools.Md5(u.Password)
	return nil
}

func (u *User) Update() (int64, error) {
	result := db.MysqlDB.Save(u)
	return result.RowsAffected, result.Error
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {

	return nil
}

func (u *User) ResetPassword(newPassword string) (int64, error) {
	u.Password = tools.Md5(newPassword)
	return u.Update()
}

func (u *User) GetOne() error {
	result := db.MysqlDB.First(u, u.ID)
	return result.Error
}

func (u *User) CheckUser() error {
	result := db.MysqlDB.Where("user_name =? AND password=?", u.UserName, tools.Md5(u.Password)).First(u)
	return result.Error
}

func UserListByIDs(ids ...uint) []User {
	var users []User
	result := db.MysqlDB.Where("id in(?)", tools.ImplodeUint(",", ids...)).Find(&users)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return users
}
