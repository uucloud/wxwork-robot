package sender

import (
	"context"
	"strconv"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
)

func newTestSender() Sender {
	return NewSender(context.Background(), "", 0)
}

func TestSendTxt(t *testing.T) {
	s := newTestSender()
	if err := s.Text("", "test", []string{}); err != nil {
		t.Fatal(err)
	}

	if err := s.Markdown("", "# 标题1"); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 30; i++ {
		if err := s.Markdown("", "# limit "+strconv.Itoa(i)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTemplateCard(t *testing.T) {
	s := newTestSender()
	uid, _ := uuid.GenerateUUID()
	err := s.TextNoticeCard("", &TemplateCard{
		Source: &TemplateCardSource{
			IconURL:   "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0",
			Desc:      "企业微信",
			DescColor: 0,
		},
		ActionMenu: &TemplateCardActionMenu{
			Desc: "消息气泡副交互辅助文本说明",
			ActionList: []struct {
				Text string `json:"text"`
				Key  string `json:"key"`
			}{
				{
					Text: "接收推送",
					Key:  "action_key1",
				},
				{
					Text: "不再推送",
					Key:  "action_key2",
				},
			},
		},
		MainTitle: &TemplateCardMainTitle{
			Title: "欢迎使用企业微信",
			Desc:  "您的好友正在邀请您加入企业微信",
		},
		EmphasisContent: &TemplateCardEmphasisContent{
			Title: "100",
			Desc:  "数据含义",
		},
		QuoteArea: &TemplateCardQuoteArea{
			Type:      1,
			URL:       "https://work.weixin.qq.com/?from=openApi",
			AppID:     "APPID",
			Pagepath:  "PAGEPATH",
			Title:     "引用文本标题",
			QuoteText: "Jack：企业微信真的很好用~\nBalian：超级好的一款软件！",
		},
		SubTitleText: "下载企业微信还能抢红包！",
		HorizontalContentList: []*TemplateCardHorizontal{
			{
				Keyname: "邀请人",
				Value:   "张三",
			},
			{
				Keyname: "企微官网",
				Value:   "点击访问",
				Type:    1,
				URL:     "https://work.weixin.qq.com/?from=openApi",
			},
		},
		JumpList: []*TemplateJumpList{
			{
				Type:  1,
				URL:   "https://work.weixin.qq.com/?from=openApi",
				Title: "企业微信官网",
			},
			//{  //parent department not found
			//	Type:     2,
			//	AppID:    "APPID",
			//	Pagepath: "PAGEPATH",
			//	Title:    "跳转小程序",
			//},
		},
		CardAction: &TemplateCardAction{
			Type:     1,
			URL:      "https://work.weixin.qq.com/?from=openApi",
			AppID:    "APPID",
			Pagepath: "PAGEPATH",
		},
		TaskID: uid,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultipleCard(t *testing.T) {
	s := newTestSender()
	uid, _ := uuid.GenerateUUID()
	err := s.MultipleCard("", &TemplateCard{
		Source: &TemplateCardSource{
			IconURL: "https://www.bilibili.com/favicon.ico?v=1",
			Desc:    "企业微信",
		},

		MainTitle: &TemplateCardMainTitle{
			Title: "欢迎使用企业微信",
			Desc:  "您的好友正在邀请您加入企业微信",
		},
		SelectList: []*TemplateCardSelectList{
			{
				QuestionKey: "question_key1",
				Title:       "选择器标签1",
				SelectedID:  "id_one",
				OptionList: []TemplateCardSelectOption{
					{
						ID:   "id_one",
						Text: "选项1",
					},
					{
						ID:   "id_two",
						Text: "选项2",
					},
				},
			},
		},
		SubmitButton: &TemplateCardSubmitButton{
			Text: "提交",
			Key:  "submit_key",
		},

		TaskID: uid,
	})

	if err != nil {
		t.Fatal(err)
	}
}
