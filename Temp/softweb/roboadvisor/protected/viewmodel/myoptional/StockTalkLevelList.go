package myoptional

import (
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"time"
)

type StockTalkLevelList []*myoptional_model.StockTalk

func (list StockTalkLevelList) Len() int {
	return len(list)
}

func (list StockTalkLevelList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// 根据置顶等级优先排序 如果等级相同则 时间最新的优先
func (list StockTalkLevelList) Less(i, j int) bool {
	if list[j].TalkLevel == list[i].TalkLevel {
		if time.Time(list[j].CreateTime) == time.Time(list[i].CreateTime) {
			return list[j].StockTalkId < list[i].StockTalkId
		}
		return time.Time(list[j].CreateTime).Before(time.Time(list[i].CreateTime))
	} else {
		return list[j].TalkLevel < list[i].TalkLevel
	}
}
