package gotest

import (
	"github.com/NoahAmethyst/dispatch-center/dingding"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
	"github.com/NoahAmethyst/dispatch-center/utils/log"
	"testing"
)

func Test_DingMarkdown(t *testing.T) {
	if err := dingding.DingCli.Send("Test", "Test content", "https://www.google.com", dispatch_pb.DingMType_Markdown); err != nil {
		log.Error().Msgf(err.Error())
	}

}
