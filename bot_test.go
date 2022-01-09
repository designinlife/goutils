package goutils

import (
	"testing"
)

func TestFeishuTextMessage(t *testing.T) {
	msg, err := NewFeishuTextMessage("This is an content.")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(msg))
}

func TestFeishuRichMessage(t *testing.T) {
	msg := NewFeishuRichMessage("This is an title.")
	msg.NewContent()
	msg.AddText("这是内容")
	msg.AddAt("@designinlife")
	msg.AddHref("网易", "https://www.163.com")

	t.Logf("%s", msg)
}
