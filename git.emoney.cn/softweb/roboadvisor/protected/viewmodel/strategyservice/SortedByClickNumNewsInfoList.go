package strategyservice

import (
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
)

type SortedByClickNumNewsInfoList []*model.NewsInfo

func (list SortedByClickNumNewsInfoList) Len() int {
	return len(list)
}

func (list SortedByClickNumNewsInfoList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list SortedByClickNumNewsInfoList) Less(i, j int) bool {
	return list[i].ClickNum > list[j].ClickNum
}
