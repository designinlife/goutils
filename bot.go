package goutils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	AccessToken       string
	SecretKey         string
	TenantAccessToken string // 租户访问凭证, 用于上传图片
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

func (s *FeishuBotSender) UploadImage(filename string) (string, error) {
	if !IsFile(filename) {
		return "", errors.New("Source file does not exist.")
	}

	file, _ := os.Open(filename)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, _ := writer.CreateFormFile("image", filepath.Base(file.Name()))
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	writer.WriteField("image_type", "message")

	r, _ := http.NewRequest("POST", "https://open.feishu.cn/open-apis/im/v1/images", body)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.TenantAccessToken))
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.Errorf("Response status error: %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// fmt.Println(fmt.Sprintf("Result: %s", string(content)))
	logger.Debugf("Result: %s", string(content))

	r1 := struct {
		Code int `json:"code"`
		Data struct {
			ImageKey string `json:"image_key"`
		} `json:"data"`
		Msg string `json:"msg"`
	}{}

	if err = json.Unmarshal(content, &r1); err != nil {
		return "", err
	}

	if r1.Code != 0 {
		return "", errors.Errorf("Upload error: %s (%d)", r1.Msg, r1.Code)
	}

	return r1.Data.ImageKey, nil
}

func (s *FeishuBotSender) Send(v BotMessage) error {
	if s.AccessToken == "" {
		return errors.New("Access token is invalid.")
	}

	s.sign(v)

	data, err := v.Body()
	if err != nil {
		return err
	}

	client := NewHttpClient()
	resp, err := client.Post(fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", s.AccessToken), &HttpRequest{
		JSON: data,
	})
	if err != nil {
		return err
	}

	r1 := struct {
		StatusCode    int    `json:"StatusCode"`
		StatusMessage string `json:"StatusMessage"`
		Code          int    `json:"code"`
		Msg           string `json:"msg"`
	}{}

	// logger.Debugf("Response: %v", resp)

	if err = json.Unmarshal(resp.Body, &r1); err != nil {
		return errors.Errorf("Response parse error: %v", err)
	}

	if r1.Code != 0 {
		return errors.Errorf("%s (%d)", r1.Msg, r1.Code)
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

type DingtalkActionCardSingleMessage struct {
	ActionCard struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		BtnOrientation string `json:"btnOrientation"`
		SingleTitle    string `json:"singleTitle"`
		SingleURL      string `json:"singleURL"`
	} `json:"actionCard"`
	Msgtype string `json:"msgtype"`
}

func (s *DingtalkActionCardSingleMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewDingtalkActionCardSingleMessage(title, content, singleTitle, singleURL string) *DingtalkActionCardSingleMessage {
	msg := &DingtalkActionCardSingleMessage{}
	msg.Msgtype = "actionCard"
	msg.ActionCard.BtnOrientation = "0"
	msg.ActionCard.Title = title
	msg.ActionCard.Text = content
	msg.ActionCard.SingleTitle = singleTitle
	msg.ActionCard.SingleURL = singleURL

	return msg
}

type DingtalkActionCardMessage struct {
	Msgtype    string `json:"msgtype"`
	ActionCard struct {
		Title          string                            `json:"title"`
		Text           string                            `json:"text"`
		BtnOrientation string                            `json:"btnOrientation,omitempty"`
		Btns           []DingtalkActionCardMessageButton `json:"btns,omitempty"`
	} `json:"actionCard"`
}

type DingtalkActionCardMessageButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

func NewDingtalkActionCardMessage(title, content string) *DingtalkActionCardMessage {
	msg := &DingtalkActionCardMessage{}
	msg.Msgtype = "actionCard"
	msg.ActionCard.Title = title
	msg.ActionCard.Text = content
	msg.ActionCard.BtnOrientation = "0"

	return msg
}

func (s *DingtalkActionCardMessage) AddButton(title, actionUrl string) {
	s.ActionCard.Btns = append(s.ActionCard.Btns, DingtalkActionCardMessageButton{
		Title:     title,
		ActionURL: actionUrl,
	})
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
		Links []DingtalkFeedCardMessageLink `json:"links"`
	} `json:"feedCard"`
}

type DingtalkFeedCardMessageLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}

func NewDingtalkFeedCardMessage() *DingtalkFeedCardMessage {
	msg := &DingtalkFeedCardMessage{}
	msg.Msgtype = "feedCard"

	return msg
}

func (s *DingtalkFeedCardMessage) AddLink(title, messageURL, picURL string) {
	s.FeedCard.Links = append(s.FeedCard.Links, DingtalkFeedCardMessageLink{
		Title:      title,
		MessageURL: messageURL,
		PicURL:     picURL,
	})
}

func (s *DingtalkFeedCardMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *DingtalkBotSender) Send(v BotMessage) error {
	if s.AccessToken == "" {
		return errors.New("Access token is invalid.")
	}

	data, err := v.Body()
	if err != nil {
		return err
	}

	value := url.Values{}

	if s.SecretKey != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, s.SecretKey)
		h := hmac.New(sha256.New, []byte(s.SecretKey))
		_, err := h.Write([]byte(stringToSign))
		if err != nil {
			return err
		}
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		value.Set("access_token", s.AccessToken)
		value.Set("timestamp", fmt.Sprintf("%d", timestamp))
		value.Set("sign", signature)
	} else {
		value.Set("access_token", s.AccessToken)
	}

	client := NewHttpClient()
	resp, err := client.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?%s", value.Encode()), &HttpRequest{
		JSON: data,
	})
	if err != nil {
		return err
	}

	r1 := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}

	if err = json.Unmarshal(resp.Body, &r1); err != nil {
		return err
	}

	if r1.Errcode != 0 {
		return errors.Errorf("%s (%d)", r1.Errmsg, r1.Errcode)
	}

	logger.Debugf("Response: %v", resp)

	return nil
}

type WxWorkTextMessage struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content             string   `json:"content"`
		MentionedList       []string `json:"mentioned_list"`
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

func (s *WxWorkTextMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkTextMessage(content string) *WxWorkTextMessage {
	msg := &WxWorkTextMessage{}
	msg.Msgtype = "text"
	msg.Text.Content = content

	return msg
}

func (s *WxWorkTextMessage) AtAll() {
	s.Text.MentionedList = append(s.Text.MentionedList, "@all")
}

func (s *WxWorkTextMessage) AtUserIds(userIds ...string) {
	s.Text.MentionedList = append(s.Text.MentionedList, userIds...)
}

func (s *WxWorkTextMessage) AtMobiles(mobiles ...string) {
	s.Text.MentionedMobileList = append(s.Text.MentionedMobileList, mobiles...)
}

type WxWorkMarkdownMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func (s *WxWorkMarkdownMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkMarkdownMessage(content string) *WxWorkMarkdownMessage {
	msg := &WxWorkMarkdownMessage{}
	msg.Msgtype = "markdown"
	msg.Markdown.Content = content

	return msg
}

type WxWorkImageMessage struct {
	Msgtype string `json:"msgtype"`
	Image   struct {
		Base64 string `json:"base64"`
		Md5    string `json:"md5"`
	} `json:"image"`
}

func (s *WxWorkImageMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkImageMessage(filePath string) (*WxWorkImageMessage, error) {
	if !IsFile(filePath) {
		return nil, errors.New("Source file does not exist.")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileStat.Size() > 2*1024*1024 {
		return nil, errors.New("The maximum image (before encoding) cannot exceed 2Mb, and supports JPG and PNG formats.")
	}

	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)

	content, err2 := ioutil.ReadAll(tee)
	if err2 != nil {
		return nil, err2
	}
	encoded := base64.StdEncoding.EncodeToString(content)

	hash := md5.New()
	if _, err = io.Copy(hash, &buf); err != nil {
		return nil, err
	}
	hashInBytes := hash.Sum(nil)
	md5Str := hex.EncodeToString(hashInBytes)

	msg := &WxWorkImageMessage{}
	msg.Msgtype = "image"
	msg.Image.Base64 = encoded
	msg.Image.Md5 = md5Str
	// msg.Image.Md5 = "e14ee653af8f0e2cf3ce74d3b211353f"

	return msg, nil
}

type WxWorkNewsMessage struct {
	Msgtype string `json:"msgtype"`
	News    struct {
		Articles []wxWorkNewsMessageArticle `json:"articles"`
	} `json:"news"`
}

type wxWorkNewsMessageArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
}

func (s *WxWorkNewsMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *WxWorkNewsMessage) AddArticle(title, description, url, picurl string) {
	s.News.Articles = append(s.News.Articles, wxWorkNewsMessageArticle{
		Title:       title,
		Description: description,
		URL:         url,
		Picurl:      picurl,
	})
}

func NewWxWorkNewsMessage() *WxWorkNewsMessage {
	msg := &WxWorkNewsMessage{}
	msg.Msgtype = "news"

	return msg
}

type WxWorkFileMessage struct {
	Msgtype string `json:"msgtype"`
	File    struct {
		MediaID string `json:"media_id"`
	} `json:"file"`
}

func (s *WxWorkFileMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkFileMessage(mediaId string) *WxWorkFileMessage {
	msg := &WxWorkFileMessage{}
	msg.Msgtype = "file"
	msg.File.MediaID = mediaId

	return msg
}

type WxWorkTextNoticeMessage struct {
	Msgtype      string `json:"msgtype"`
	TemplateCard struct {
		CardType string `json:"card_type"`
		Source   struct {
			IconURL   string `json:"icon_url,omitempty"`
			Desc      string `json:"desc,omitempty"`
			DescColor int    `json:"desc_color,omitempty"`
		} `json:"source,omitempty"`
		MainTitle struct {
			Title string `json:"title,omitempty"`
			Desc  string `json:"desc,omitempty"`
		} `json:"main_title"`
		EmphasisContent struct {
			Title string `json:"title,omitempty"`
			Desc  string `json:"desc,omitempty"`
		} `json:"emphasis_content,omitempty"`
		QuoteArea struct {
			Type      int    `json:"type,omitempty"`
			URL       string `json:"url,omitempty"`
			Appid     string `json:"appid,omitempty"`
			Pagepath  string `json:"pagepath,omitempty"`
			Title     string `json:"title,omitempty"`
			QuoteText string `json:"quote_text,omitempty"`
		} `json:"quote_area,omitempty"`
		SubTitleText          string `json:"sub_title_text,omitempty"`
		HorizontalContentList []struct {
			Keyname string `json:"keyname"`
			Value   string `json:"value,omitempty"`
			Type    int    `json:"type,omitempty"`
			URL     string `json:"url,omitempty"`
			MediaID string `json:"media_id,omitempty"`
		} `json:"horizontal_content_list,omitempty"`
		JumpList []struct {
			Type     int    `json:"type,omitempty"`
			URL      string `json:"url,omitempty"`
			Title    string `json:"title"`
			Appid    string `json:"appid,omitempty"`
			Pagepath string `json:"pagepath,omitempty"`
		} `json:"jump_list,omitempty"`
		CardAction struct {
			Type     int    `json:"type"`
			URL      string `json:"url,omitempty"`
			Appid    string `json:"appid,omitempty"`
			Pagepath string `json:"pagepath,omitempty"`
		} `json:"card_action"`
	} `json:"template_card"`
}

func (s *WxWorkTextNoticeMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkTextNoticeMessage() *WxWorkTextNoticeMessage {
	msg := &WxWorkTextNoticeMessage{}
	msg.Msgtype = "template_card"
	msg.TemplateCard.CardType = "text_notice"

	return msg
}

type WxWorkNewsNoticeMessage struct {
	Msgtype      string `json:"msgtype"`
	TemplateCard struct {
		CardType string `json:"card_type"`
		Source   struct {
			IconURL   string `json:"icon_url,omitempty"`
			Desc      string `json:"desc,omitempty"`
			DescColor int    `json:"desc_color,omitempty"`
		} `json:"source,omitempty"`
		MainTitle struct {
			Title string `json:"title"`
			Desc  string `json:"desc,omitempty"`
		} `json:"main_title"`
		CardImage struct {
			URL         string  `json:"url"`
			AspectRatio float64 `json:"aspect_ratio,omitempty"`
		} `json:"card_image"`
		ImageTextArea struct {
			Type     int    `json:"type,omitempty"`
			URL      string `json:"url,omitempty"`
			Title    string `json:"title,omitempty"`
			Desc     string `json:"desc,omitempty"`
			ImageURL string `json:"image_url"`
		} `json:"image_text_area,omitempty"`
		QuoteArea struct {
			Type      int    `json:"type,omitempty"`
			URL       string `json:"url,omitempty"`
			Appid     string `json:"appid,omitempty"`
			Pagepath  string `json:"pagepath,omitempty"`
			Title     string `json:"title,omitempty"`
			QuoteText string `json:"quote_text,omitempty"`
		} `json:"quote_area,omitempty"`
		VerticalContentList []struct {
			Title string `json:"title"`
			Desc  string `json:"desc,omitempty"`
		} `json:"vertical_content_list,omitempty"`
		HorizontalContentList []struct {
			Keyname string `json:"keyname"`
			Value   string `json:"value,omitempty"`
			Type    int    `json:"type,omitempty"`
			URL     string `json:"url,omitempty"`
			MediaID string `json:"media_id,omitempty"`
		} `json:"horizontal_content_list,omitempty"`
		JumpList []struct {
			Type     int    `json:"type,omitempty"`
			URL      string `json:"url,omitempty"`
			Title    string `json:"title"`
			Appid    string `json:"appid,omitempty"`
			Pagepath string `json:"pagepath,omitempty"`
		} `json:"jump_list,omitempty"`
		CardAction struct {
			Type     int    `json:"type"`
			URL      string `json:"url,omitempty"`
			Appid    string `json:"appid,omitempty"`
			Pagepath string `json:"pagepath,omitempty"`
		} `json:"card_action"`
	} `json:"template_card"`
}

func (s *WxWorkNewsNoticeMessage) Body() ([]byte, error) {
	v, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewWxWorkNewsNoticeMessage() *WxWorkNewsNoticeMessage {
	msg := &WxWorkNewsNoticeMessage{}
	msg.Msgtype = "template_card"
	msg.TemplateCard.CardType = "news_notice"

	return msg
}

type WxWorkBotSender struct {
	AccessToken string
}

func (s *WxWorkBotSender) UploadMedia(filename string) (string, error) {
	if !IsFile(filename) {
		return "", errors.New("Source file does not exist.")
	}

	file, _ := os.Open(filename)
	defer file.Close()
	stat, _ := file.Stat()

	fileName := filepath.Base(filename)
	fileSize := stat.Size()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.WriteField("filename", fileName)
	writer.WriteField("filelength", fmt.Sprintf("%d", fileSize))

	r, _ := http.NewRequest("POST", fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file", s.AccessToken), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.Errorf("Response status error: %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// fmt.Println(fmt.Sprintf("Result: %s", string(content)))
	logger.Debugf("Result: %s", string(content))

	r1 := struct {
		Errcode   int    `json:"errcode"`
		Errmsg    string `json:"errmsg"`
		Type      string `json:"type"`
		MediaID   string `json:"media_id"`
		CreatedAt string `json:"created_at"`
	}{}

	if err = json.Unmarshal(content, &r1); err != nil {
		return "", err
	}

	if r1.Errcode != 0 {
		return "", errors.Errorf("Upload error: %s (%d)", r1.Errmsg, r1.Errcode)
	}

	return r1.MediaID, nil
}

func (s *WxWorkBotSender) Send(v BotMessage) error {
	if s.AccessToken == "" {
		return errors.New("Access token is invalid.")
	}

	data, err := v.Body()
	if err != nil {
		return err
	}

	value := url.Values{}
	value.Set("key", s.AccessToken)

	client := NewHttpClient()
	resp, err := client.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?%s", value.Encode()), &HttpRequest{
		JSON: data,
	})
	if err != nil {
		return err
	}

	r1 := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}

	if err = json.Unmarshal(resp.Body, &r1); err != nil {
		return err
	}

	if r1.Errcode != 0 {
		return errors.Errorf("%s (%d)", r1.Errmsg, r1.Errcode)
	}

	logger.Debugf("Response: %v", resp)

	return nil
}
