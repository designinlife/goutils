package goutils

import (
	"encoding/json"
)

type BotSender interface {
	Send(v []byte) error
}

type DingtalkBotSender struct {
}

func (s *DingtalkBotSender) Send(v []byte) error {
	return nil
}

type FeishuMessage struct {
	Timestamp string `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
	MsgType   string `json:"msg_type"`
}

type FeishuTextMessage struct {
	FeishuMessage
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (s *FeishuTextMessage) String() string {
	v, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(v)
}

// NewFeishuTextMessage 创建飞书文本消息。
func NewFeishuTextMessage(content string) *FeishuTextMessage {
	msg := &FeishuTextMessage{}
	msg.Content.Text = content
	msg.MsgType = "text"

	return msg
}

type feishuRichMessageContent struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type FeishuRichMessage struct {
	FeishuMessage
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string                       `json:"title"`
				Content [][]feishuRichMessageContent `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
	ContentPosition int `json:"-"`
}

func (s *FeishuRichMessage) String() string {
	v, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(v)
}

func NewFeishuRichMessage(title string) *FeishuRichMessage {
	msg := &FeishuRichMessage{}
	msg.MsgType = "post"
	msg.Content.Post.ZhCn.Title = title
	msg.Content.Post.ZhCn.Content = make([][]feishuRichMessageContent, 0)

	return msg
}

func (s *FeishuRichMessage) NewContent() {
	s.Content.Post.ZhCn.Content = append(s.Content.Post.ZhCn.Content, []feishuRichMessageContent{})

	s.ContentPosition = len(s.Content.Post.ZhCn.Content) - 1
}

func (s *FeishuRichMessage) AddText(v string) {
	s.Content.Post.ZhCn.Content[s.ContentPosition] = append(s.Content.Post.ZhCn.Content[s.ContentPosition], feishuRichMessageContent{
		Tag:  "text",
		Text: v,
	})
}

func (s *FeishuRichMessage) AddHref(label, href string) {
	s.Content.Post.ZhCn.Content[s.ContentPosition] = append(s.Content.Post.ZhCn.Content[s.ContentPosition], feishuRichMessageContent{
		Tag:  "a",
		Text: label,
		Href: href,
	})
}

func (s *FeishuRichMessage) AddAt(userId string) {
	s.Content.Post.ZhCn.Content[s.ContentPosition] = append(s.Content.Post.ZhCn.Content[s.ContentPosition], feishuRichMessageContent{
		Tag:    "at",
		UserID: userId,
	})
}

type FeishuBotSender struct {
}

func (s *FeishuBotSender) Send(v []byte) error {
	return nil
}
