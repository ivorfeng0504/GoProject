package userloginlog

import (
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/user"
)

// GetUserLoginCount 获取用户累计登天数
func GetUserLoginCount(cid int64) (loginDays int, err error) {
	userSrv := user.NewUserService()
	loginDays, _, _, _, _, err = userSrv.LoginDaysAndProductByCID(cid)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetUserLoginCount 获取用户累计登天数 异常 cid=%d count=%d", cid, loginDays)
	} else {
		global.InnerLogger.DebugFormat("GetUserLoginCount 获取用户累计登天数 cid=%d count=%d", cid, loginDays)
	}
	return loginDays, err
}

// GetUserLoginCountSerial 获取用户连登记录
func GetUserLoginCountSerial(cid int64) (loginLog UserLoginLog, err error) {
	userSrv := user.NewUserService()
	_, continueLoginDays, _, _, _, err := userSrv.LoginDaysAndProductByCID(cid)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetUserLoginCountSerial 获取用户连登记录 异常 cid=%d", cid)
		return loginLog, err
	}
	global.InnerLogger.DebugFormat("GetUserLoginCountSerial 获取用户连登记录 cid=%d continueLoginDays=%d", cid, continueLoginDays)
	if continueLoginDays == 0 {
		global.InnerLogger.DebugFormat("GetUserLoginCountSerial 获取用户连登记录 异常 设置为默认登录天数 cid=%d", cid)
		continueLoginDays = 1
	}
	loginLog = UserLoginLog{
		LoginCountSerial: continueLoginDays,
		//LastLoginTime:    time.Now(),
	}
	return loginLog, err
}
