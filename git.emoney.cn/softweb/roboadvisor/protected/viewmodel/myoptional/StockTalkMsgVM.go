package myoptional

import (
	"github.com/devfeel/mapper"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
)

//yqq推送数据-微股吧
type StockTalkMsgVM struct {
	//直播间Id
	LiveRoomId int64
	//直播间名称
	LiveRoomName string
	//股票代码集合
	StockInfoList []model.StockInfo
	//图片数组
	ImageList []string
	//评论内容
	Content string
	//益圈圈消息发送时间
	SendTime mapper.JSONTime
}
