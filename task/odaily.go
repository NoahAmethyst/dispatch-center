package task

import (
	"context"
	"fmt"
	"github.com/NoahAmethyst/dispatch-center/cluster/spider_svc"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/spider_pb"
	"github.com/NoahAmethyst/dispatch-center/sender"
	cron2 "github.com/NoahAmethyst/dispatch-center/utils/cron"
	"github.com/NoahAmethyst/dispatch-center/utils/file_util"
	"github.com/NoahAmethyst/dispatch-center/utils/log"
	"sync"
	"time"
)

type odailySent struct {
	ids []int64
	max int
	sync.RWMutex
}

var OdailySentRecord odailySent

func (o *odailySent) Put(id int64) {
	o.Lock()
	defer o.Unlock()
	o.ids = append(o.ids, id)
	if len(o.ids) > o.max {
		for i := 0; i < len(o.ids)-o.max; i++ {
			o.ids = o.ids[1:]
		}
	}
}

func (o *odailySent) Exist(id int64) bool {
	o.RLock()
	defer o.RUnlock()
	exist := false
	for _, _id := range o.ids {
		if _id == id {
			exist = true
			break
		}
	}
	return exist
}

func (o *odailySent) SaveData() {
	log.Info().Msgf("Save odaily sent record data")
	path := file_util.GetFileRoot()
	if _, err := file_util.WriteJsonFile(o.ids, path, "odailySentRecord", false); err != nil {
		log.Error().Msgf("Save odaily news faield:%s", err.Error())
	} else {
		_ = file_util.TCCosUpload("cache", "odailySentRecord.json", fmt.Sprintf("%s/%s", path, "odailySentRecord.json"))
	}
}

var markdownTemplate = `
### [%s](%s)

%s

发布时间：%s
`

func PushOdailyNews(ctx context.Context, bot dispatch_pb.Bot) {
	cli := spider_svc.SvcCli()
	req := new(spider_pb.SpiderReq)
	resp, err := cli.OdailyFeeds(ctx, req)
	if err != nil {
		log.Error().Msgf("Fetch Odaily feeds failed:%s", err.Error())
		return
	}

	df := "2006-01-02 15:04:05"

	for i := len(resp.OdailyFeeds) - 1; i >= 0; i-- {
		odailyFeed := resp.OdailyFeeds[i]
		if OdailySentRecord.Exist(odailyFeed.Id) {
			log.Debug().Msgf("Odaily Feed 【%s】has been sent", odailyFeed.Title)
			continue
		}
		message := dispatch_pb.Message{
			Meta: &dispatch_pb.MetaMessage{
				Title: odailyFeed.Title,
			},
			Qq: nil,
			Dingding: &dispatch_pb.DingdingMessage{
				Mt: dispatch_pb.DingMType_Markdown,
			},
			T: bot,
		}

		if len(odailyFeed.ReferenceUrl) > 0 {
			message.Meta.ReferenceUrl = odailyFeed.ReferenceUrl
		} else {
			message.Meta.ReferenceUrl = odailyFeed.Url
		}

		publishdeAt := time.Unix(odailyFeed.PublishedAt, 0)

		message.Meta.Content = fmt.Sprintf(markdownTemplate, odailyFeed.Title, message.Meta.ReferenceUrl,
			odailyFeed.Description, publishdeAt.Format(df))

		if err := sender.Push(&message); err != nil {
			log.Error().Msgf("Task:Send odaily feed failed:%s", err.Error())
		}
		OdailySentRecord.Put(odailyFeed.Id)
	}

}

func RegisterTask(ctx context.Context, bot dispatch_pb.Bot, taskDuration string) {
	task := func() {
		PushOdailyNews(ctx, bot)
	}
	cron2.AddCronJob(task, taskDuration)
}

func init() {
	OdailySentRecord = odailySent{
		ids:     make([]int64, 0, 200),
		max:     200,
		RWMutex: sync.RWMutex{},
	}
	sentRecord := make([]int64, 0, 200)
	path := file_util.GetFileRoot()
	if err := file_util.LoadJsonFile(fmt.Sprintf("%s/odailySentRecord.json", path), &sentRecord); err != nil {
		log.Info().Msgf("retry load odailySentRecord json from tencent cos")
		_err := file_util.TCCosDownload("cache", "odailySentRecord.json", fmt.Sprintf("%s/%s", path, "odailySentRecord.json"))
		if _err == nil {
			_ = file_util.LoadJsonFile(fmt.Sprintf("%s/odailySentRecord.json", path), &sentRecord)
		} else {
			log.Error().Msgf("Load odaily sent record data from Tencent COS faied:%s", _err.Error())
		}
	}
}
