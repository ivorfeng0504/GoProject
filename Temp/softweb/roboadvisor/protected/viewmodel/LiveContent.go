package viewmodel

import (
	"github.com/devfeel/mapper"
)

//直播内容
type LiveContent struct {

	//发表内容
	Content string

	//发表内容(图片)
	//ContentImg string

	//创建时间
	CreateTimeStr string

	//创建时间
	CreateTime mapper.JSONTime
	// 来源（1播主自发 2 播主回答推荐 3打赏显示）
	LiveFrom int
	// 消息内容
	Message string

	//消息性质
	// 0 无  1利好 2利空
	MessageType int

	//关联个股/板块
	// 股票代码|股票名称|涨幅|股票类型,股票代码|股票名称|涨幅|股票类型
	StockList string

	//关联个股/板块
	StockInfoList []*StockInfo

	// 是否置顶
	IsTop bool

	//直播名称
	LiveName string

}
