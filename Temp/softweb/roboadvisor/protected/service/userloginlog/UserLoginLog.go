package userloginlog

type UserLoginLog struct {

	//用户UID
	UID int64
	//用户最后登录时间
	//LastLoginTime time.Time
	//用户累计登录次数
	LoginCount int
	//用户最后连登天数
	LoginCountSerial int
}
