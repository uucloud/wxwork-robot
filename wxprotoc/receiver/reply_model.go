package receiver

type BaseReply interface {
	Trigger(event *CallEvent) bool

	MsgType() []MsgType

	SupportPassiveReply() bool
	PassiveReply(event *CallEvent) (string, PassiveReplyType)
	Reply(event *CallEvent) error
}

type baseReply struct {
	mt []MsgType

	trigger      func(event *CallEvent) bool
	reply        func(event *CallEvent) error
	passiveReply func(event *CallEvent) (string, PassiveReplyType)
}

func NewReply(
	msgType []MsgType,
	trigger func(event *CallEvent) bool,
	reply func(event *CallEvent) error,
	passiveReply func(event *CallEvent) (string, PassiveReplyType),
) BaseReply {
	if len(msgType) == 0 {
		panic("msgType can't be empty")
	}

	if trigger == nil {
		panic("trigger can't be empty")
	}

	if reply == nil && passiveReply == nil {
		panic("at least one of reply and passiveReply can be used")
	}

	return &baseReply{
		mt:           msgType,
		trigger:      trigger,
		reply:        reply,
		passiveReply: passiveReply,
	}
}

func (br *baseReply) Trigger(event *CallEvent) bool {
	return br.trigger(event)
}

func (br *baseReply) MsgType() []MsgType {
	return br.mt
}

func (br *baseReply) SupportPassiveReply() bool {
	return br.passiveReply != nil
}

func (br *baseReply) PassiveReply(event *CallEvent) (string, PassiveReplyType) {
	return br.passiveReply(event)
}

func (br *baseReply) Reply(event *CallEvent) error {
	return br.reply(event)
}
