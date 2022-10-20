/*
 * @Author: uucloud
 * @Date: 2022-08-25 18:00:35
 * @Last Modified by: uucloud
 * @Last Modified time: 2022-08-25 18:00:56
 */

package main

import (
	"flag"
	"time"

	"github.com/uucloud/wxwork-robot/client"
	"github.com/uucloud/wxwork-robot/utils"
	"github.com/uucloud/wxwork-robot/wxprotoc/receiver"

	"go.uber.org/zap"
)

func main() {
	flag.Parse()

	initRobot()

	ctx := utils.SetupSignalContext()

	utils.Logger.Info("start", zap.String("time", time.Now().String()))
	s := utils.NewServer(&utils.ServerConfig{
		Addr: ":9000",
	}, client.M.Handler())

	go s.Server.ListenAndServe()
	<-ctx.Done()
}

func initRobot() {
	robot, err := client.M.AddRobot("your robot name", "95037c0d-f7ec-4cbb-aeec-xxxx", "your receive token", "encoding aes key")
	if err != nil {
		panic(err)
	}

	robot.Rec.Register(receiver.NewReply(
		//回复类型
		[]receiver.MsgType{receiver.MText, receiver.MMixed, receiver.MImage},

		//触发条件
		func(event *receiver.CallEvent) bool {
			return true
		},

		//群聊时的回复内容
		func(event *receiver.CallEvent) error {
			return robot.Send.Text(event.ChatID(), utils.TrimAt(robot.Name, event.Texts()), []string{event.UserID()})
		},

		//单聊时的回复内容
		func(event *receiver.CallEvent) (string, receiver.PassiveReplyType) {
			return event.Texts(), receiver.PRTMarkdown
		},
	))
}
