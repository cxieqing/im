package models

import "testing"

func TestUserMessage(t *testing.T) {
	fromUserid := uint(2)
	toUserid := uint(3)
	message := UserMessage{
		ContentType: TextContent,
		Content:     "hello hello",
		From:        fromUserid,
		To:          toUserid,
		IsSend:      HasSend,
		Len:         float32(len("hello hello")),
	}
	t.Run("create user message", func(t *testing.T) {
		message.ID, _ = message.CreateSendMsg()
		if message.ID < 1 {
			t.Errorf("unexpect new recode id %d", message.ID)
		}
	})

	t.Run("message history list", func(t *testing.T) {
		var (
			messages []UserMessage
			total    int64
		)
		messages, total = message.HistoryMsgList(1, 1)
		if total < 1 {
			t.Errorf("want total at least 1 got %d", total)
		}
		if messages[0].ID != message.ID {
			t.Errorf("want message id %d  got %d", message.ID, messages[0].ID)
		}
	})
}
