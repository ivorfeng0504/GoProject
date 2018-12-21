package strategyservice

import (
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
	"time"
)

type NewsInfoList []*model.NewsInfo

func (list NewsInfoList) Len() int {
	return len(list)
}

func (list NewsInfoList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list NewsInfoList) Less(i, j int) bool {
	return time.Time(list[j].LastModifyTime).Before(time.Time(list[i].LastModifyTime))
}
