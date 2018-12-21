package model

import (
	"time"
	"git.emoney.cn/softweb/roboadvisor/protected/model/user"
)

type UserInfo struct {
	//主键 自增长
	ID int

	//UID 客户端UID
	UID int64

	//sso组ID  绑定使用
	GID int64

	//用户类型（0：游客 1：手机 :2：微信 3：QQ  ）
	UserType int

	//密文手机
	MobileX string

	//掩码手机
	MobileMask string

	//昵称
	NickName string

	//头像
	Headportrait string

	//微信openID
	OpenID_WeChat string

	//QQ openID
	OpenID_QQ string

	//用户等级
	UserLevel int

	//是否删除
	IsDeleted int

	//创建时间
	CreateTime time.Time

	//最后更新时间
	LastLoginTime time.Time

	Account string

	PID int

	//版本(免费：1  体验：2   收费：3)
	PIDType int

	SID int

	TID int

	//绑定相关业务使用（sso解密的用户名em  明文手机号  wx_uniqueID  qq_openID）
	UserName string

	//已绑定账号列表
	BindAccountList []*model.BoundAccount_Response

	//用户到期日
	EndDate string

	//所在省份
	Province string

	//所在城市
	City string

	//服务网点
	ServiceAgentName string

	ClientVersion int

}
