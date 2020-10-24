package models

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestGroup(t *testing.T) {
	ownerId := uint(12)
	memberId := uint(123)
	otherMemberId := uint(11)
	group := Group{Name: "先遣小分队1"}
	t.Run("Group create", func(t *testing.T) {
		group.GroupOwner = ownerId
		if id, _ := group.Create(memberId); id < 1 {
			t.Errorf("unexpect new recode id %d", id)
		}
	})

	t.Run("Group getOne", func(t *testing.T) {
		g := &Group{}
		g.ID = group.ID
		if err := g.GetOne(); err != nil {
			t.Errorf("unexpect GetOne error  %s", err)
		}
	})

	t.Run("Group delete", func(t *testing.T) {
		if err := group.Delete(); err != nil {
			t.Errorf("unexpect error %s", err)
		}
	})

	t.Run("Group add member", func(t *testing.T) {
		group = Group{Name: "先遣小分队2"}
		group.GroupOwner = ownerId
		if id, _ := group.Create(memberId); id < 1 {
			t.Errorf("unexpect new recode id %d", id)
		}
		err := group.AddMember(memberId)
		if err == nil {
			t.Errorf("expect error  duplicate")
		}
		err = group.AddMember(otherMemberId)
		if err != nil {
			t.Errorf("unexpect AddMember error  %s", err)
		}
		memberNum := group.MemberNum()
		if memberNum != 3 {
			t.Errorf("got member num %d want  %d", memberNum, 3)
		}
	})

	t.Run("Groups get by user id", func(t *testing.T) {
		groups := GetGroupsByUserId(ownerId)
		if len(groups) < 1 {
			t.Errorf("unexpect groups len %d", len(groups))
		}
		hasNewRecode := false
		for _, v := range groups {
			if v.ID == group.ID {
				hasNewRecode = true
			}
		}
		if !hasNewRecode {
			t.Errorf("new Recode not exists")
		}
	})

	t.Run("Group member out", func(t *testing.T) {
		err := group.OutMember(otherMemberId)
		if err != nil {
			t.Errorf("unexpect OutMember err %s", err)
		}
		memberNum := group.MemberNum()
		if memberNum != 2 {
			t.Errorf("got member num %d want  %d", memberNum, 2)
		}
		err = group.OutMember(memberId)
		if err != nil {
			t.Errorf("unexpect OutMember err %s", err)
		}
		g := &Group{}
		g.ID = group.ID
		err = g.GetOne()
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expect delete group err %s", gorm.ErrRecordNotFound)
		}
	})
}
