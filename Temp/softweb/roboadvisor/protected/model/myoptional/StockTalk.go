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
	//用户UID
	UID string
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
	//是否审核通过
	IsValid bool
	//是否删除
	IsDeleted bool
	//修改时间
	ModifyTime mapper.JSONTime
	//最后修改用户
	ModifyUser string
	//置顶等级
	TalkLevel int
	//微股吧数据类型0用户数据 1基本面数据 2策略股池 3益圈圈推送 4管理后台发送
	StockTalkType int
}
