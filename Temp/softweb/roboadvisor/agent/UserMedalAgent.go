package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetUserMedal 获取用户勋章
func GetUserMedal(user model.UserInfo) (userMedalInfo UserMedalInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = user
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getusermedal", req)
	if err != nil {
		return userMedalInfo, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return userMedalInfo, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return userMedalInfo, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userMedalInfo)
	//非智盈大师版本 不显示种子用户勋章
	if (user.PID == 888020000 || user.PID == 888020400) == false {
		userMedalInfo.SeedUserLevel = 0
	}
	return userMedalInfo, nil
}

//用户勋章信息
type UserMedalInfo struct {
	//勋章-文娱标兵
	WenYuBiaoBingLevel int
	//勋章-学习委员
	XueXiWeiYuanLevel int
	//勋章-资讯大咖
	ZiXunDaKaLevel int
	//勋章-股市大亨
	GuShiDaHengLevel int
	//勋章-好问博士
	HaoWenBoShiLevel int
	//勋章-高朋满座
	GaoPengManZuoLevel int
	//勋章-互动达人
	HuDongDaRenLevel int
	//勋章-种子用户
	SeedUserLevel int
	//勋章-合格投资者
	QualifiedInvestorLevel int
}
