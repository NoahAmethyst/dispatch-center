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
	sync.RWMutex
	Dingtalk          *dingtalk.DingTalk
	Interval          int64
	LastPushTimestamp int64
}

var DingCli *dingCli

func init() {
	DingCli = &dingCli{
		RWMutex:  sync.RWMutex{},
		Interval: 60 / 20,
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
	c.Lock()
	var err error
	defer func() {
		if err == nil {
			c.LastPushTimestamp = time.Now().Unix()
		}
		c.Unlock()
	}()
	now := time.Now().Unix()
	interval := now - c.LastPushTimestamp
	if interval < c.Interval {
		time.Sleep(time.Duration(interval) * time.Second)
	}

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
	return c.Send(message.Meta.GetTitle(), message.Meta.GetContent(), message.Meta.GetReferenceUrl(), message.Dingding.GetMt())
}
