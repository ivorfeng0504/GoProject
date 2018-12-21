package myoptional

type StockTalkSimpleList []*StockTalkSimple

func (list StockTalkSimpleList) Len() int {
	return len(list)
}

func (list StockTalkSimpleList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list StockTalkSimpleList) Less(i, j int) bool {
	return list[j].CreateTime < list[i].CreateTime
}
