package myoptional

import "github.com/devfeel/mapper"

//微股吧
type StockTalkForPCClient struct {
	//主键Id
	StockTalkId int64
	//用户昵称
	NickName string
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//评论内容
	Content string
	//创建时间
	CreateTime mapper.JSONTime
	//详情URL
	DetailUrl string
	//置顶等级
	TalkLevel int
}
