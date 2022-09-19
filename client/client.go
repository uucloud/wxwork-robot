package client

import (
	"context"
	"time"
	"wework-robot/wxprotoc/receiver"
	"wework-robot/wxprotoc/sender"
)

type Robot struct {
	Name string

	Rec  receiver.Receiver
	Send sender.Sender
}

func newRobot(ctx context.Context, name, sendKey, receiveToken, receiveEncodingAESKey string, timeout time.Duration) *Robot {
	r := &Robot{
		Name: name,
		Rec:  receiver.NewReceiver(ctx, receiveToken, receiveEncodingAESKey),
		Send: sender.NewSender(ctx, sendKey, timeout),
	}
	return r
}
