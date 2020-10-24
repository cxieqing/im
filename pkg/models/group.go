package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"im/pkg/db"
	"im/pkg/tools"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Members []uint

func (m *Members) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to type value:", value))
	}
	arr := tools.ExplodeUint(",", string(bytes))
	*m = arr
	return nil
}

func (m Members) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "", nil
	}
	return tools.ImplodeUint(",", m...), nil
}

type Group struct {
	Model
	GroupOwner uint
	Name       string
	Members    Members
}

func (g *Group) Create(mids ...uint) (uint, error) {
	g.Members = append(Members{g.GroupOwner}, mids...)
	result := db.MysqlDB.Create(g)
	return g.ID, result.Error
}

func (g *Group) GetOne() error {
	result := db.MysqlDB.First(g, g.ID)
	return result.Error
}

func (g *Group) OutMember(uid uint) error {
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		result := db.MysqlDB.Clauses(clause.Locking{Strength: "UPDATE"}).First(g, g.ID)
		if result.Error != nil {
			return result.Error
		}
		if g.GroupOwner == uid {
			return g.Delete()
		}
		for k, v := range g.Members {
			if v == uid {
				g.Members = tools.UintArrayDel(g.Members, k)
			}
		}
		if len(g.Members) == 1 {
			g.DeletedAt.Scan(time.Now())
		}
		return db.MysqlDB.Save(g).Error
	})
	return err
}

func (g *Group) AddMember(uids ...uint) error {
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		result := db.MysqlDB.Clauses(clause.Locking{Strength: "UPDATE"}).First(g, g.ID)
		if result.Error != nil {
			return result.Error
		}
		for _, v := range g.Members {
			for _, m := range uids {
				if v == m {
					return errors.New("duplicate add member in group")
				}
			}
		}
		g.Members = append(g.Members, uids...)
		return db.MysqlDB.Save(g).Error
	})
	return err
}

func (g *Group) MemberNum() int {
	return len(g.Members)
}

func (g *Group) Delete() error {
	return db.MysqlDB.Delete(g).Error
}

func (g *Group) AfterDelete(tx *gorm.DB) (err error) {
	return nil
}

func GetGroupsByUserId(id uint) []Group {
	var groups []Group
	db.MysqlDB.Raw("select * from `group` where FIND_IN_SET(?, members) and deleted_at is null", id).Find(&groups)
	return groups
}
