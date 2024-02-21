package gotest

import (
	"testing"
	"time"
)

func Test_Timeformat(t *testing.T) {
	time.Local = time.FixedZone("UTC", 8*60*60)
	now := time.Now()
	df := "2006-01-02 15:04:05"
	//if loc, err := time.LoadLocation("Asia/Shanghai"); err != nil {
	//	t.Errorf("Set time location failed:%s", err.Error())
	//} else {
	//	// Convert Time to East 8 District
	//	now = now.In(loc)
	//}
	t.Logf("%s", now.Format(df))
}
