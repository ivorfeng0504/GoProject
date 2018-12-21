package expertnews

import "github.com/devfeel/mapper"

type ExpertNews_Topic struct {
	//主键 自增长
	ID int

	//主题名称
	TopicName	string

	//主题简介
	TopicSummary string

	//主题图片
	TopicImg string

	//关注指数
	ConcernIndex float32

	//关联板块信息-json
	RelatedBKInfo string

	//关联板块个股-json
	RelatedStockInfo string

	//关联板块数量
	RelatedBKCount int

	//关联个股数量
	RelatedStockCount int

	//入选日期
	InputDate mapper.JSONTime

	//周期
	Period string

	IsDeleted bool

	CreateUser string

	CreateTime mapper.JSONTime

	LastModifyTime mapper.JSONTime

	//关联板块列表
	RelatedBKList []*Topic_BK

	//关联个股列表
	RelatedStockList []*Topic_Stock
}

type Topic_BK struct {
	//板块名称
	BKName string

	//板块code
	BKCode string

	//板块涨跌幅
	BKF float32
}

type Topic_Stock struct {
	//排序编号
	StockSortIndex string

	//个股名称
	StockName string

	//个股代码
	StockCode string

	//个股简介
	StockSummary string

	StockF float32

	StockPrice float32

	StockPE string
}