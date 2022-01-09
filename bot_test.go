package goutils

import (
	"testing"
)

func TestFeishuTextMessage(t *testing.T) {
	msg := NewFeishuTextMessage("This is an content.")

	t.Logf("%s", msg)
}

func TestFeishuRichMessage(t *testing.T) {
	msg := NewFeishuRichMessage("This is an title.")
	msg.NewContent()
	msg.AddText("这是内容")
	msg.AddAt("@designinlife")
	msg.AddHref("网易", "https://www.163.com")

	t.Logf("%s", msg)
}

func TestFeishuCardMessage(t *testing.T) {
	msg := NewFeishuCardMessage("This is an title.")
	msg.AddLineContent("这是内容")
	msg.AddSplitLine()
	// msg.AddButton("按钮", "https://www.163.com")

	t.Logf("%s", msg)
}
