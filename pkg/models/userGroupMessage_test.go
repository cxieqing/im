package models

import "testing"

func TestUserGroupMessage(t *testing.T) {
	messag := UserGroupMessage{
		ToUserID: 1,
		IsSend:   UnSend,
		IsRead:   UnRead,
		GroupMessage: GroupMessage{
			ID: uint(1),
		},
	}

	t.Run("Test UserGroupMessage Create", func(t *testing.T) {
		if id, _ := messag.Create(); id < 1 {
			t.Errorf("unexpect new recode id %d", id)
		}
	})

	t.Run("Test UnSendUserGroupMsgList", func(t *testing.T) {
		messages := messag.UnSendUserGroupMsgList()
		if len(messages) < 1 {
			t.Errorf("want least 1 got %d", len(messages))
		}
		hasNewRecode := false
		for _, v := range messages {
			if v.ID == messag.ID {
				hasNewRecode = true
			}
		}
		if !hasNewRecode {
			t.Errorf("new Recode not exists")
		}
	})

	t.Run("Test UserGroupMessage Delete", func(t *testing.T) {
		if err := messag.Delete(); err != nil {
			t.Errorf("unexpect error %d", err)
		}
	})
}
