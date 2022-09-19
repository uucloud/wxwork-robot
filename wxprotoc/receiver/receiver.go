package receiver

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"sync"
	"time"
	"wework-robot/utils"

	"go.uber.org/zap"
)

type Receiver interface {
	Verify(msgSign, timestamp, nonce, echoStr string) ([]byte, error)
	DecryptMsg(msgSign, timestamp, nonce string, data []byte) ([]byte, error)
	Register(reply BaseReply)
	Reply(event *CallEvent) ([]byte, error)

	SetGroupEvent(e EventType, f func(event *CallEvent) error)
	SetSingleEvent(f func(event *CallEvent) (string, PassiveReplyType))
}

type receiver struct {
	token          string
	encodingAESKey string

	msgCrypt *WXBizMsgCrypt

	rLock   sync.RWMutex
	replies map[MsgType][]BaseReply

	ctx context.Context
}

func NewReceiver(ctx context.Context, token, encodingAESKey string) Receiver {
	r := &receiver{
		token:          token,
		encodingAESKey: encodingAESKey,
		msgCrypt:       NewWXBizMsgCrypt(token, encodingAESKey, "", XmlType),
		replies:        map[MsgType][]BaseReply{},
		ctx:            ctx,
	}

	return r
}

func (r *receiver) Verify(msgSign, timestamp, nonce, echoStr string) ([]byte, error) {
	echo, cErr := r.msgCrypt.VerifyURL(msgSign, timestamp, nonce, echoStr)
	if cErr != nil {
		utils.Logger.Error("verify error", zap.String("err", cErr.Error()))
		return nil, cErr
	}
	utils.Logger.Info("verify successful")
	return echo, nil
}

func (r *receiver) DecryptMsg(msgSign, timestamp, nonce string, data []byte) ([]byte, error) {
	return r.msgCrypt.DecryptMsg(msgSign, timestamp, nonce, data)
}

func (r *receiver) Reply(event *CallEvent) (b []byte, err error) {
	//单聊只能直接回
	if event.GetChatType() == MSingle {
		b, err = r.passiveReply(event)
		if err != nil {
			return
		}
	} else {
		go r.reply(event)
		b = []byte("ok")
	}
	return
}

func (r *receiver) reply(event *CallEvent) {
	r.rLock.RLock()
	replies := r.replies[event.GetMsgType()]
	r.rLock.RUnlock()

	for _, reply := range replies {
		if reply.Trigger(event) {
			err := reply.Reply(event)
			if err != nil {
				utils.Logger.Error("reply error", zap.Error(err))
				continue
			}
			return
		}
	}
	utils.Logger.Error("not match any reply", zap.Any("event", event))
}

func (r *receiver) passiveReply(event *CallEvent) ([]byte, error) {
	r.rLock.RLock()
	replies := r.replies[event.GetMsgType()]
	r.rLock.RUnlock()

	for _, reply := range replies {
		if reply.SupportPassiveReply() && reply.Trigger(event) {
			return r.doPassiveReply(reply.PassiveReply(event))
		}
	}
	return nil, fmt.Errorf("not match any reply, event:(%+v)", event)
}

func (r *receiver) doPassiveReply(content string, replyType PassiveReplyType) ([]byte, error) {
	var reply PassiveReplyResp

	switch replyType {
	case PRTText:
		reply = PassiveReplyResp{
			MsgType: "text",
			Text: PRText{
				Content: CDATA{Value: content},
			}}
	case PRTMarkdown:
		reply = PassiveReplyResp{
			MsgType:  "markdown",
			Markdown: PRMarkdown{Content: CDATA{content}},
		}
	default:
		return nil, fmt.Errorf("not support passive reply type:%v", replyType)
	}
	utils.Logger.Info("PassiveReply", zap.Any("replyType", replyType), zap.String("content", content))
	return r.encryptMsg(reply)
}

func (r *receiver) encryptMsg(content interface{}) ([]byte, error) {
	b, err := xml.Marshal(content)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	nonce := RandString(10)
	return r.msgCrypt.EncryptMsg(string(b), timestamp, nonce)
}

func (r *receiver) Event() error {
	return nil
}

func (r *receiver) CardEvent() error {
	return nil
}

func (r *receiver) Register(reply BaseReply) {
	r.rLock.Lock()
	defer r.rLock.Unlock()

	for _, mt := range reply.MsgType() {
		r.replies[mt] = append(r.replies[mt], reply)
	}
}

func (r *receiver) SetGroupEvent(e EventType, f func(event *CallEvent) error) {
	reply := NewReply(
		[]MsgType{MEvent},
		func(event *CallEvent) bool {
			return event.GetEventType() == e
		},
		f,
		nil,
	)

	r.Register(reply)
}

func (r *receiver) SetSingleEvent(f func(event *CallEvent) (string, PassiveReplyType)) {
	reply := NewReply(
		[]MsgType{MEvent},
		func(event *CallEvent) bool {
			return event.GetEventType() == EventEnterChat
		},
		nil,
		f,
	)
	r.Register(reply)
}
