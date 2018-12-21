package userproduct

import (
	"git.emoney.cn/softweb/roboadvisor/protected/service/user"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"time"
)

// GetUserPCProductList 获取用户PC端产品列表
func GetUserPCProductList(cid int64) (userProductList []*UserProduct, err error) {
	userSrv := user.NewUserService()
	_, _, product, activateTimeStr, endTimeStr, err := userSrv.LoginDaysAndProductByCID(cid)
	if err != nil {
		return userProductList, err
	}
	if len(product) == 0 || len(endTimeStr) == 0 {
		return userProductList, err
	}
	stateDesc := ""
	now := time.Now()
	endTime, _ := _time.ParseTime(endTimeStr)
	if endTime.After(now) {
		stateDesc = "未过期"
	} else {
		stateDesc = "已过期"
	}
	userProductList = append(userProductList, &UserProduct{
		UID:          cid,
		ProductName:  product,
		ActivateTime: activateTimeStr,
		ExpireTime:   endTimeStr,
		StateDesc:    stateDesc,
	})
	return userProductList, nil
}
