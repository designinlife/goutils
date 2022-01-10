package goutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type BotSender interface {
	Send(v BotMessage) error
}

type BotMessage interface {
	Body() ([]byte, error)
}

type feishuMessage struct {
	Timestamp string `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
	MsgType   string `json:"msg_type"`
}

type FeishuTextMessage struct {
	feishuMessage
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (s *FeishuTextMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
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
	feishuMessage
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

func (s *FeishuRichMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
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
	feishuMessage
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

func (s *FeishuCardMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
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
	AccessToken string
	SecretKey   string
}

func (s *FeishuBotSender) sign(v interface{}) error {
	if s.SecretKey == "" {
		return nil
	}

	timestamp := time.Now().Unix()
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + s.SecretKey
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	switch vtype := v.(type) {
	case *FeishuCardMessage:
		vtype.Timestamp = strconv.FormatInt(timestamp, 10)
		vtype.Sign = signature
	case *FeishuRichMessage:
		vtype.Timestamp = strconv.FormatInt(timestamp, 10)
		vtype.Sign = signature
	case *FeishuTextMessage:
		vtype.Timestamp = strconv.FormatInt(timestamp, 10)
		vtype.Sign = signature
	default:
		return errors.New("非法的类型参数。")
	}

	return nil
}

func (s *FeishuBotSender) Send(v BotMessage) error {
	s.sign(v)

	data, err := v.Body()
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Data: %s", string(data)))

	client := NewHttpClient()
	resp, err := client.Post(fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", s.AccessToken), &HttpRequest{
		JSON: data,
	})
	if err != nil {
		return err
	}

	logger.Debugf("Response: %v", resp)

	return nil
}

type DingtalkBotSender struct {
	AccessToken string
	SecretKey   string
}

type DingtalkTextMessage struct {
	At struct {
		AtMobiles []string `json:"atMobiles,omitempty"`
		AtUserIds []string `json:"atUserIds,omitempty"`
		IsAtAll   bool     `json:"isAtAll,omitempty"`
	} `json:"at,omitempty"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Msgtype string `json:"msgtype"`
}

func NewDingtalkTextMessage(title, content string, atAll bool) *DingtalkTextMessage {
	msg := &DingtalkTextMessage{}
	msg.Msgtype = "text"
	if title != "" {
		msg.Text.Content = fmt.Sprintf("%s\n%s", title, content)
	} else {
		msg.Text.Content = content
	}
	if atAll {
		msg.At.IsAtAll = true
	}

	return msg
}

func (s *DingtalkTextMessage) AtMobiles(mobiles ...string) {
	if !s.At.IsAtAll {
		s.At.AtMobiles = append(s.At.AtMobiles, mobiles...)
	}
}

func (s *DingtalkTextMessage) AtUserIds(userIds ...string) {
	if !s.At.IsAtAll {
		s.At.AtMobiles = append(s.At.AtMobiles, userIds...)
	}
}

func (s *DingtalkTextMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type DingtalkLinkMessage struct {
	Msgtype string `json:"msgtype"`
	Link    struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicURL     string `json:"picUrl,omitempty"`
		MessageURL string `json:"messageUrl"`
	} `json:"link"`
}

func NewDingtalkLinkMessage(title, content, messageUrl, picUrl string) *DingtalkLinkMessage {
	msg := &DingtalkLinkMessage{}
	msg.Msgtype = "link"
	msg.Link.Title = title
	msg.Link.Text = content
	msg.Link.MessageURL = messageUrl

	if picUrl != "" {
		msg.Link.PicURL = picUrl
	}

	return msg
}

func (s *DingtalkLinkMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type DingtalkMarkdownMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles,omitempty"`
		AtUserIds []string `json:"atUserIds,omitempty"`
		IsAtAll   bool     `json:"isAtAll,omitempty"`
	} `json:"at,omitempty"`
}

func NewDingtalkMarkdownMessage(title, content string, atAll bool) *DingtalkMarkdownMessage {
	msg := &DingtalkMarkdownMessage{}
	msg.Msgtype = "markdown"
	msg.Markdown.Title = title
	msg.Markdown.Text = content

	if atAll {
		msg.At.IsAtAll = true
	}

	return msg
}

func (s *DingtalkMarkdownMessage) AtMobiles(mobiles ...string) {
	if !s.At.IsAtAll {
		s.At.AtMobiles = append(s.At.AtMobiles, mobiles...)
	}
}

func (s *DingtalkMarkdownMessage) AtUserIds(userIds ...string) {
	if !s.At.IsAtAll {
		s.At.AtMobiles = append(s.At.AtMobiles, userIds...)
	}
}

func (s *DingtalkMarkdownMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type DingtalkActionCardMessage struct {
	Msgtype    string `json:"msgtype"`
	ActionCard struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		BtnOrientation string `json:"btnOrientation,omitempty"`
		Btns           []struct {
			Title     string `json:"title"`
			ActionURL string `json:"actionURL"`
		} `json:"btns"`
	} `json:"actionCard"`
}

func (s *DingtalkActionCardMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type DingtalkFeedCardMessage struct {
	Msgtype  string `json:"msgtype"`
	FeedCard struct {
		Links []struct {
			Title      string `json:"title"`
			MessageURL string `json:"messageURL"`
			PicURL     string `json:"picURL"`
		} `json:"links"`
	} `json:"feedCard"`
}

func (s *DingtalkFeedCardMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *DingtalkBotSender) Send(v BotMessage) error {
	data, err := v.Body()
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Data: %s", string(data)))

	var dingtalkApiUrl string

	if s.SecretKey != "" {
		timestamp := time.Now().Unix()
		stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + s.SecretKey
		var data []byte
		h := hmac.New(sha256.New, []byte(stringToSign))
		_, err := h.Write(data)
		if err != nil {
			return err
		}
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		dingtalkApiUrl = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%s&sign=%s", s.AccessToken, strconv.FormatInt(timestamp, 10), signature)
	} else {
		dingtalkApiUrl = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", s.AccessToken)
	}

	client := NewHttpClient()
	resp, err := client.Post(dingtalkApiUrl, &HttpRequest{
		JSON: data,
	})
	if err != nil {
		return err
	}

	logger.Debugf("Response: %v", resp)

	return nil
}
