package model

import (
	"github.com/devfeel/mapper"
)

type StrategyInfo struct{

	/// 编号
	ID int

	/// 策略编号
	StrategyID int

	/// 策略名称
	StrategyName string

	/// 产品平台
	ProductPlat string

	///Logo图片地址
	LogoUrl string

	/// 策略简介
	Summary string
	
	/// 策略描述
	Description string

	/// 策略风格
	Style string
	
	/// 配图图片地址
	ImagesUrl string
	
	/// 微信服务号
	WebchatServiceNO string
	
	/// 客服电话
	ServiceTel  string

	/// 客服QQ
	ServiceQQ string
	
	/// 客服QQ链接
	ServiceQQLink string
	
	/// 客服EMail
	EMail string
	
	/// 创建时间
	DateTime mapper.JSONTime
	
	/// 修改时间
	ModifyTime mapper.JSONTime

	/// 是否删除
	IsDeleted bool
	
	/// 状态
	Status int
	
	/// 创建者
	Creator string

}