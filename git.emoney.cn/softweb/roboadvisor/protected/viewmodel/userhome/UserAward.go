package viewmodel


type UserAward struct {
	//主键Id
	UserAwardId int64
	//奖品名称
	AwardName string
	//奖品Id
	AwardId int64
	//奖品图片
	AwardImg string
	//创建时间
	CreateTime string
	//介绍视频
	IntroduceVideo string
	//有效期
	AvailableDay   string
	//到期时间
	ExpireTime string
	//奖品时效状态
	AwardStateDesc string
	//奖品类型描述
	AwardTypeDesc string
	//是否是视频
	IsVideo bool
}

