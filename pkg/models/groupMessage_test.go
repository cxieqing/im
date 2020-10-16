package models

import "testing"

func TestGroupMessage(t *testing.T) {

	content := "i am is group message"
	messag := GroupMessage{
		ContentType: TextContent,
		Content:     content,
		From:        uint(5),
		To:          uint(10),
		Len:         float32(len(content)),
	}

	t.Run("Test GroupMessage Create", func(t *testing.T) {
		if id, _ := messag.Create(); id < 1 {
			t.Errorf("unexpect new recode id %d", id)
		}
	})
	t.Run("Test GroupMessage List", func(t *testing.T) {
		if messags := messag.List(1, 1); len(messags) == 1 {
			t.Errorf("want 1 got %d", len(messags))
		}
	})
}
