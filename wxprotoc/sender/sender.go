package sender

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/time/rate"
)

type Sender interface {
	Text(chatID string, msg string, mentionedList []string) error
	Markdown(chatID string, msg string) error
	Image(chatID string, b []byte) error
	News(chatID string, arts []Article) error

	TextNoticeCard(chatID string, t *TemplateCard) error
	NewsNoticeCard() error
	ButtonCard() error
	VoteCard() error
	MultipleCard() error
}

type sender struct {
	key string

	url string

	limiter *rate.Limiter
	client  *gorequest.SuperAgent

	timeout time.Duration
	ctx     context.Context
}

func NewSender(ctx context.Context, key string, timeout time.Duration) Sender {
	if timeout <= 0 {
		timeout = time.Second * 20
	}
	req := gorequest.New()
	req.Client.Timeout = timeout

	s := &sender{
		key:     key,
		url:     sendURL(key),
		ctx:     ctx,
		client:  req,
		limiter: rate.NewLimiter(0.16, 15),
		timeout: timeout,
	}
	return s
}

func (s *sender) Text(ChatID string, msg string, mentionedList []string) error {
	req := SendReq{
		CommonSend: CommonSend{
			ChatID:  ChatID,
			MsgType: "text",
		},
		Text: &TextType{
			Content:       msg,
			MentionedList: mentionedList,
		},
	}

	return s.send(req)
}

func (s *sender) Markdown(chatID string, msg string) error {
	req := SendReq{
		CommonSend: CommonSend{
			ChatID:  chatID,
			MsgType: "markdown",
		},
		Markdown: &MarkdownType{
			Content: msg,
		},
	}
	return s.send(req)
}

func (s *sender) Image(chatID string, b []byte) error {
	b64 := getBase64(b)
	m5 := getMd5(b)

	req := SendReq{
		CommonSend: CommonSend{
			ChatID:  chatID,
			MsgType: "image",
		},
		Image: &ImageType{
			Base64: b64,
			Md5:    m5,
		},
	}
	return s.send(req)
}

func (s *sender) News(chatID string, arts []Article) error {
	req := SendReq{
		CommonSend: CommonSend{
			ChatID:  chatID,
			MsgType: "news",
		},
		News: &NewsType{
			Articles: arts,
		},
	}
	return s.send(req)
}

func (s *sender) TextNoticeCard(chatID string, t *TemplateCard) error {
	t.CardType = "text_notice"
	req := SendReq{
		CommonSend: CommonSend{
			ChatID:  chatID,
			MsgType: "template_card",
		},
		TemplateCard: t,
	}

	return s.send(req)
}

func (s *sender) NewsNoticeCard() error {
	panic("need to implement")
	return nil
}

func (s *sender) ButtonCard() error {
	panic("need to implement")
	return nil
}

func (s *sender) VoteCard() error {
	panic("need to implement")
	return nil
}

func (s *sender) MultipleCard() error {
	panic("need to implement")
	return nil
}

func (s *sender) send(content interface{}) error {
	if !s.limiter.Allow() {
		ctx, _ := context.WithTimeout(s.ctx, s.timeout)
		if err := s.limiter.Wait(ctx); err != nil {
			return err
		}
	}

	var (
		resp SendResponse
		errs []error
	)
	_, _, errs = s.client.Post(s.url).SendStruct(content).EndStruct(&resp)
	if len(errs) > 0 {
		return errs[0]
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("%d/%s", resp.ErrCode, resp.ErrMsg)
	}
	return nil
}

func sendURL(key string) string {
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
	return fmt.Sprintf(url, key)
}

func getMd5(b []byte) string {
	has := md5.Sum(b)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func getBase64(b []byte) string {
	sEnc := base64.StdEncoding.EncodeToString(b)
	return sEnc
}
