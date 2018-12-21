package live

import "time"

type User struct {
	//用户编号
	UserId int

	//用户名  手机会是掩码 180****1234 仅用显示不做标识
	UserName string

	//创建时间
	CreateTime time.Time

	//状态 -1：屏蔽 0：正常
	Status int

	//最后登录时间
	LastLoginTime time.Time

	//用户类型 1邮箱 2手机 3EM(通行证) 4QQ 5微博 6智投EM注册
	UserType string

	//昵称
	NickName string

	//账号(智投为手机号)
	Account string

	//金融平台UID
	UID int64

	//用户直播间权限-非数据库字段
	RoomList []int

	//pid
	PID int

	//服务商ID
	AgentID string



	//通行证ID
	//PassportId int64

	// 通行证加密手机号码
	//Mobile string

	//本地应用加密手机号码
	//appMobile string

	// 应用AppId
	//AppId int

	//用户掩码用户名
	//MaskUserName string

	//真实姓名
	//RealName string

	//风险承受能力
	//RiskAbility string

	//适当性匹配
	//RiskStatus int

	//加密用户名
	//EncryptUserName string

	// 加密手机号码
	//EncryptMobile string

	//加密手机号码
	//BindTime time.Time

	//手机掩码 180****1234
	//MobileMask string

	//是否播主
	//isRoom int

	//用户的房间号
	//RoomId int

	//是否来自移动端 1来自移动手机端 0来自PC端
	//IsFromMobile int

	//来源 PC Mobile
	//From string

	//IP
	//IP string

	// 渠道码
	//ChannelCode string

	//渠道名称
	//ChannelName string

	// 推广当前页面Url
	//ChannelUrl string

	//是否来自推广用户
	//IsChannel int

	//是否黑名单
	//IsBlackList int

	//加入黑名单时间
	//BlackDate time.Time

	//是否黑名单
	//BlId int

	// 头像编号
	//Avator string

	//第三方用户编号
	//ExUserId string
	//QQ string
	//userFrom int
	//userLevel int
}
