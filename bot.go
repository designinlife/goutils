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

type feishuCardMessageElement struct {
	Tag  string `json:"tag"`
	Text struct {
		Content string `json:"content"`
		Tag     string `json:"tag"`
	} `json:"text,omitempty"`
	Actions []feishuCardMessageElementAction `json:"actions,omitempty"`
}

type feishuCardMessageElementAction struct {
	Tag  string `json:"tag"`
	Text struct {
		Content string `json:"content"`
		Tag     string `json:"tag"`
	} `json:"text"`
	URL   string   `json:"url"`
	Type  string   `json:"type"`
	Value struct{} `json:"value"`
}

type FeishuCardMessage struct {
	FeishuMessage
	Card struct {
		Config struct {
			WideScreenMode bool `json:"wide_screen_mode"`
			EnableForward  bool `json:"enable_forward"`
		} `json:"config"`
		Elements []feishuCardMessageElement `json:"elements"`
		Header   struct {
			Title struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"title"`
		} `json:"header"`
	} `json:"card"`
}

func (s *FeishuCardMessage) String() string {
	v, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(v)
}

func NewFeishuCardMessage(title string) *FeishuCardMessage {
	msg := &FeishuCardMessage{}
	msg.MsgType = "interactive"
	msg.Card.Config.EnableForward = true
	msg.Card.Config.WideScreenMode = true
	msg.Card.Header.Title.Tag = "plain_text"
	msg.Card.Header.Title.Content = title

	return msg
}

func (s *FeishuCardMessage) AddLineContent(v string) {
	elem := feishuCardMessageElement{
		Tag: "div",
	}
	elem.Text.Tag = "lark_md"
	elem.Text.Content = v

	s.Card.Elements = append(s.Card.Elements, elem)
}

func (s *FeishuCardMessage) AddSplitLine() {
	elem := feishuCardMessageElement{
		Tag: "hr",
	}

	s.Card.Elements = append(s.Card.Elements, elem)
}

func (s *FeishuCardMessage) AddButton(label, href string) {
	action := feishuCardMessageElementAction{
		Tag:  "button",
		Type: "default",
		URL:  href,
	}
	action.Text.Tag = "lark_md"
	action.Text.Content = label

	elem := feishuCardMessageElement{
		Tag: "action",
	}
	elem.Actions = append(elem.Actions, action)

	s.Card.Elements = append(s.Card.Elements, elem)
}

type FeishuBotSender struct {
}

func (s *FeishuBotSender) Send(v []byte) error {
	return nil
}
