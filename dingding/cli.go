package dingding

import (
	"github.com/NoahAmethyst/dispatch-center/constant"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
	"github.com/NoahAmethyst/dispatch-center/utils/log"
	"github.com/blinkbean/dingtalk"
	"os"
	"sync"
	"time"
)

type dingCli struct {
	Dingtalk *dingtalk.DingTalk
	Interval time.Duration
	sync.Once
	revcMessage chan *dispatch_pb.Message
}

var DingCli *dingCli

func init() {
	DingCli = &dingCli{
		Interval:    time.Second * 3,
		Once:        sync.Once{},
		revcMessage: make(chan *dispatch_pb.Message, 100),
	}
	token := os.Getenv(constant.DING_TOKEN)
	secret := os.Getenv(constant.DING_SECRET)
	if len(secret) == 0 {
		log.Error().Msgf("Dingding secret not set.")
		return
	}
	DingCli.Dingtalk = dingtalk.InitDingTalkWithSecret(token, secret)
}

func (c *dingCli) Send(title, content, referenceUrl string, t dispatch_pb.DingMType) error {
	var err error
	switch t {
	case dispatch_pb.DingMType_Text:
		err = c.Dingtalk.SendTextMessage(content)
	case dispatch_pb.DingMType_Link:
		err = c.Dingtalk.SendLinkMessage(title, content, "", referenceUrl)
	case dispatch_pb.DingMType_Markdown:
		err = c.Dingtalk.SendMarkDownMessage(title, content)
	default:
		log.Warn().Msgf("Not support message type:%+v", t)
	}

	if err != nil {
		log.Error().Msgf("Push ding message failed:%s", err.Error())
	}
	return err

}

func (c *dingCli) Push(message *dispatch_pb.Message) error {

	go func(_message *dispatch_pb.Message) {
		c.revcMessage <- _message
	}(message)

	c.Once.Do(func() {
		for {
			select {
			case msg := <-c.revcMessage:
				time.Sleep(c.Interval)
				if err := c.Send(msg.Meta.GetTitle(), msg.Meta.GetContent(), msg.Meta.GetReferenceUrl(), msg.Dingding.GetMt()); err != nil {
					log.Error().Msgf("Send dingtalk failed:%s", err.Error())
					_ = c.Push(msg)
				}
			}
		}
	})

	return nil

}
