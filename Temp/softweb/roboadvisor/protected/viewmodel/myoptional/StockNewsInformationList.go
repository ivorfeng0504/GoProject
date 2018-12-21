package myoptional

import (
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"time"
)

type StockNewsInformationList []*myoptional_model.StockNewsInformation

func (list StockNewsInformationList) Len() int {
	return len(list)
}

func (list StockNewsInformationList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list StockNewsInformationList) Less(i, j int) bool {
	return time.Time(list[j].PublishTime).Before(time.Time(list[i].PublishTime))
}
