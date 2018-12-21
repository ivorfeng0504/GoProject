package myoptional

import "github.com/devfeel/mapper"

//微股吧
type StockTalk struct {
	//主键Id
	StockTalkId int64
	//用户昵称
	NickName string
	//头像
	Avatar string
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//图片信息集合 JSON格式
	ImageInfoList string
	//评论内容
	Content string
	//创建时间
	CreateTime mapper.JSONTime
}
