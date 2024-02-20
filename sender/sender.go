package sender

import (
	"fmt"
	"github.com/NoahAmethyst/dispatch-center/dingding"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
	"github.com/pkg/errors"
)

type ISender interface {
	Push(message *dispatch_pb.Message) error
}

func Push(message *dispatch_pb.Message) error {
	var iSender ISender
	switch message.T {
	case dispatch_pb.Bot_Dingding:
		iSender = dingding.DingCli
	default:
		return errors.New(fmt.Sprintf("Unsupported bot type:%+v", message.GetT()))
	}

	return iSender.Push(message)
}
