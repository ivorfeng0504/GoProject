package live

import (
	"github.com/devfeel/mapper"
)

//直播室
type LiveRoom struct {
	//唯一标识
	Id int

	// 关注ID
	Fid int

	//直播名称
	LiveName string

	//首拼字母
	FirstLetter string

	//通行证Id
	PassportId int64

	//通行证帐号
	UserName string

	//统一后台帐号
	AdminUserName string

	//直播状态(0关闭1开启)
	LiveStatus int

	//图片地址url
	LiveImg string

	// 播主简介
	LiveIntro string

	// 创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	// 标签名称
	TagNameStr string

	//标签Id字符串
	TagIdStr string

	//粉丝数
	FansNum int

	//访问量
	VisitAddrNum int

	// 是否被关注 1关注
	IsFocus int

	//关联普通直播室Id
	PId int

	//是否为VIP直播室（0否1是）
	IsVIP int

	//是否有对应的VIP直播室（0否1是）
	hasVIP int

	// 是否隐藏（0否1是）
	IsHide int

	//贡献值
	Contribution float64

	//剩余贡献值
	RemainContribution float64

	//剩余贡献值是否满足VIP条件
	IsSatisfyContribution int

	//用户是否有VIP权限
	IsAwardVIP int

	//直播内容
	Content string

	//最后更新时间
	NewContentTime mapper.JSONTime

	//到期时间
	ExpiredTime mapper.JSONTime

	//投资顾问
	AdviserName string

	//投资编号
	AdviserNo string

	//直播话题
	TopicName string

	// 是否删除
	IsDelete int

	// 删除时间
	DeleteTime mapper.JSONTime
}
