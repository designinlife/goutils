package goutils

import (
	"os"
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
	msg.AddButton("按钮", "https://www.163.com")

	t.Logf("%s", msg)
}

func TestFeishuTextMessage_Send(t *testing.T) {
	msg := NewFeishuTextMessage("This is an content.")

	// logger.SetLevel(logger.DebugLevel)

	sender := &FeishuBotSender{AccessToken: os.Getenv("FS_ACCESS_TOKEN"), SecretKey: os.Getenv("FS_SECRET_KEY")}
	err := sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("%s", msg)
}

func TestFeishuRichMessage_Send(t *testing.T) {
	msg := NewFeishuRichMessage("This is an title.")
	msg.NewContent()
	msg.AddText("这是内容")
	msg.AddAt("@designinlife")
	msg.AddHref("网易", "https://www.163.com")

	sender := &FeishuBotSender{AccessToken: os.Getenv("FS_ACCESS_TOKEN"), SecretKey: os.Getenv("FS_SECRET_KEY")}
	err := sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("%s", msg)
}

func TestFeishuCardMessage_Send(t *testing.T) {
	msg := NewFeishuCardMessage("This is an title.")
	msg.AddLineContent("这是内容")
	msg.AddSplitLine()
	msg.AddButton("按钮", "https://www.163.com")

	sender := &FeishuBotSender{AccessToken: os.Getenv("FS_ACCESS_TOKEN"), SecretKey: os.Getenv("FS_SECRET_KEY")}
	err := sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("%s", msg)
}

func TestDingtalkTextMessage_Send(t *testing.T) {
	msg := NewDingtalkTextMessage("", "This is an content.", false)

	body, _ := msg.Body()

	t.Logf("%s", string(body))

	sender := &DingtalkBotSender{AccessToken: os.Getenv("DT_ACCESS_TOKEN"), SecretKey: os.Getenv("DT_SECRET_KEY")}
	err := sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestWxWorkTextMessage_Send(t *testing.T) {
	msg := NewWxWorkTextMessage("This is an content.")

	body, _ := msg.Body()

	t.Logf("%s", string(body))

	sender := &WxWorkBotSender{AccessToken: os.Getenv("WX_ACCESS_TOKEN")}
	err := sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestWxWorkImageMessage_Send(t *testing.T) {
	msg, err := NewWxWorkImageMessage("d:/1.jpeg")
	if err != nil {
		t.Error(err)
	}

	body, _ := msg.Body()

	t.Logf("%s", string(body))
	return

	sender := &WxWorkBotSender{AccessToken: os.Getenv("WX_ACCESS_TOKEN")}
	err = sender.Send(msg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
