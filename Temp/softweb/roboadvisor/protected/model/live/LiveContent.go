package live

import (
	"github.com/devfeel/mapper"
)

//直播内容
type LiveContent struct {

	//唯一标识
	Id int

	//话题ID
	TId int

	// 直播室ID
	LId int

	//消息ID
	MesId int

	//发表内容
	Content string

	//发表内容(图片)
	//ContentImg string

	//创建时间
	CreateTime mapper.JSONTime

	//最后修改时间
	LastModifyTime mapper.JSONTime

	// 来源（1播主自发 2 播主回答推荐 3打赏显示）
	LiveFrom int

	//用户ID
	UserId int

	//用户ID
	AdminUserId int64

	// 通行证ID
	PassportId int64

	// 消息内容
	Message string

	//金额
	Money float64

	//红包记录是否来自其VIP直播室
	IsFromVIP int

	//是否删除
	IsDelete int

	//删除时间
	DeleteTime mapper.JSONTime

	//消息性质
	// 0 无  1利好 2利空
	MessageType int

	//关联个股/板块
	// 股票代码|股票名称|涨幅|股票类型,股票代码|股票名称|涨幅|股票类型
	StockList string

	// 是否置顶
	IsTop bool

	//直播名称
	LiveName string

	// 用户名  手机会是掩码 180****1234 仅用显示不做标识
	UserName string

	// 昵称
	NickName string
}
