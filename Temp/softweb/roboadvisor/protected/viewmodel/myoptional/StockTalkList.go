package myoptional

import (
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"time"
)

type StockTalkList []*myoptional_model.StockTalk

func (list StockTalkList) Len() int {
	return len(list)
}

func (list StockTalkList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list StockTalkList) Less(i, j int) bool {
	return time.Time(list[j].CreateTime).Before(time.Time(list[i].CreateTime))
}
