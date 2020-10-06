package user

import (
	"errors"
	"im/pkg/db"

	"gorm.io/gorm"
)

type User struct {
	UserName string
	Password string
	NikeName string
	Icon     string
	gorm.Model
}

func (u *User) Add() (uint, error) {
	result := db.MysqlDB.Create(u)
	return u.ID, result.Error
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	return nil
}

func (u *User) Update() (int64, error) {
	result := db.MysqlDB.Save(u)
	return result.RowsAffected, result.Error
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {

	return nil
}

func (u *User) getOne() error {
	result := db.MysqlDB.First(u, u.ID)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//@todo log
	}
	return result.Error
}
