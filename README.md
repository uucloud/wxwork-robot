# wxwork-robot

企业微信机器人SDK

支持发送消息和接收消息并注册处理方法，支持多机器人管理


## 使用方式
#### 初始化机器人
``` golang
robot, err := client.M.AddRobot("your robot name", "95037c0d-f7ec-4cbb-aeec-xxxx", "your receive token", "encoding aes key")
if err != nil {
    panic(err)
}
```
#### 发送
```golang
if err := robot.Send.Text("", "hello world", []string{}); err!=nil{
	panic(err)
}
```

#### 注册消息处理
```
//机器人注册接收处理事件
//注册顺序决定了触发优先级
robot.Rec.Register(receiver.NewReply(
    //触发类型
    []receiver.MsgType{receiver.MText, receiver.MMixed, receiver.MImage},

    //触发条件
    func(event *receiver.CallEvent) bool {
        //任意回复都会触发
        return true
    },

    //群聊时的回复内容
    func(event *receiver.CallEvent) error {
        return robot.Send.Text(event.ChatID(), "你好", []string{event.UserID()})
    },

    //单聊时的回复内容
    func(event *receiver.CallEvent) (string, receiver.PassiveReplyType) {
        return event.Texts(), receiver.PRTMarkdown
    },
))

```

#### 注册事件处理
```
robot.Rec.SetGroupEvent(receiver.EventDeleteFromChat, func(event *receiver.CallEvent) error {
	    fmt.Println("机器人从群里删除了", event.ChatID())
		return nil
	})
```

## 功能
### 发送
* [x] 发送文字
* [X] 发送markdown
* [x] 发送图片
* [x] 发送图文链接
* [x] 发送TextNotice模版卡片
* [ ] 发送NewsNotice模版卡片
* [ ] 发送Button模版卡片
* [ ] 发送Vote模版卡片
* [x] 发送多项选择模版卡片
### 接收
* [x] 回调校验
* [x] 回调处理，包括单聊和群聊
* [x] 事件处理
* [ ] 卡片事件回调
### 其他
* [x] 多机器人管理
* [x] http server
* [ ] 日志及端口配置化


## 一些说明
* 企业微信机器人不能主动私聊，只能在5s内响应，且只支持文字内容。
