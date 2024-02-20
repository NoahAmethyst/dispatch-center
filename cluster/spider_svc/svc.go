package spider_svc

import (
	"github.com/NoahAmethyst/dispatch-center/cluster/rpc"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/spider_pb"
)

func SvcCli() spider_pb.SpiderServiceClient {
	return spider_pb.NewSpiderServiceClient(rpc.GetConn(rpc.CLISpider))
}
