package receiver

import "encoding/xml"

type MsgType string
type ChatType string
type PassiveReplyType string
type EventType string

const (
	MText       MsgType = "text"
	MImage      MsgType = "image"
	MMixed      MsgType = "mixed"
	MEvent      MsgType = "event"
	MTCardEvent MsgType = "template_card_event"

	MSingle ChatType = "single"
	MGroup  ChatType = "group"

	PRTText     PassiveReplyType = "text"
	PRTMarkdown PassiveReplyType = "markdown"

	EventAddToChat      EventType = "add_to_chat"
	EventDeleteFromChat EventType = "delete_from_chat"
	EventEnterChat      EventType = "enter_chat"
)

// 用户群里@机器人或者单聊中向机器人发送文本消息或图文混排消息
type CallEvent struct {
	XMLName    xml.Name `xml:"xml"`
	WebhookUrl CDATA    `xml:"WebhookUrl"`
	MsgId      CDATA    `xml:"MsgId"`
	ChatId     CDATA    `xml:"ChatId"`
	ChatType   CDATA    `xml:"ChatType"`
	From       struct {
		UserId CDATA `xml:"UserId"`
	} `xml:"From"`
	MsgType      CDATA `xml:"MsgType"` // text/image/mixed/event
	Text         Text
	Image        Image
	MixedMessage MixedMessage
	Event        struct {
		EventType CDATA `xml:"EventType"`
	} `xml:"Event"`
}

type MixedMessage struct {
	MsgItem []MsgItem `xml:"MsgItem"`
}

type MsgItem struct {
	MsgType CDATA
	Text    Text
	Image   Image
}

type Text struct {
	Content CDATA
}

type Image struct {
	ImageURL CDATA
}

func (cm *CallEvent) MsgID() string {
	return cm.MsgId.Value
}

func (cm *CallEvent) ChatID() string {
	return cm.ChatId.Value
}

func (cm *CallEvent) UserID() string {
	return cm.From.UserId.Value
}

func (cm *CallEvent) Texts() string {
	ret := ""
	ret += cm.Text.Content.Value
	for _, mix := range cm.MixedMessage.MsgItem {
		if mix.MsgType.Value == string(MText) {
			ret += mix.Text.Content.Value
		}
	}
	return ret
}

func (cm *CallEvent) ImageURLs() []string {
	var ret []string
	if cm.Image.ImageURL.Value != "" {
		ret = append(ret, cm.Image.ImageURL.Value)
	}
	for _, mix := range cm.MixedMessage.MsgItem {
		if mix.MsgType.Value == string(MImage) {
			ret = append(ret, mix.Image.ImageURL.Value)
		}
	}
	return ret
}

func (cm *CallEvent) GetMsgType() MsgType {
	return MsgType(cm.MsgType.Value)
}

func (cm *CallEvent) GetChatType() ChatType {
	return ChatType(cm.ChatType.Value)
}

func (cm *CallEvent) GetEventType() EventType {
	return EventType(cm.Event.EventType.Value)
}

// 群机器人被添加到或被移除出群聊
type AddOrRemoveResp struct {
}

// 用户进入机器人单聊界面
type EnterPrivateChatResp struct {
}

// 用户在单聊界面中给机器人发送图片消息
type PrivateChatPicResp struct {
}

// 用户点击模版卡片中的交互控件
type TemplateCarReactResp struct {
}

type PassiveReplyResp struct {
	XMLName  xml.Name `xml:"xml"`
	MsgType  string
	Text     PRText
	Markdown PRMarkdown
}

type PRText struct {
	Content       CDATA
	MentionedList MentionedList
}

type MentionedList struct {
	Item []CDATA `xml:"Item"`
}

type PRMarkdown struct {
	Content CDATA
}
