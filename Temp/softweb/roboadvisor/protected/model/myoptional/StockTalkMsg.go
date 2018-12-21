package myoptional

import "github.com/devfeel/mapper"

//yqq同步过来的评论数据-for 微股吧
type StockTalkMsg struct {
	//主键Id
	StockTalkMsgId int64
	//直播间Id
	LiveRoomId int64
	//直播间名称
	LiveRoomName string
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//图片数组 []string 的json字符串
	ImageList string
	//评论内容
	Content string
	//益圈圈消息发送时间
	SendTime mapper.JSONTime
	//是否已经推送到微股吧
	IsSendStockTalk bool
	//创建时间
	CreateTime mapper.JSONTime
	//是否删除
	IsDeleted bool
	//修改时间
	ModifyTime mapper.JSONTime
	//最后修改用户
	ModifyUser string
}

type ImageInfo struct {
	ImageUrl  string
	ImageDesc string
}
