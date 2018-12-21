package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

type UserAgent struct {
}

type AddUserRequest struct {
	Account  string
	NickName string
	UID      int64
	Source   string
}

// GetUserById 通过用户Id获取用户信息
func GetUserById(userId int) (user *livemodel.User, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userId
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/getuserbyuserid", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &user)
	return user, err
}

// GetUserByAccount 通过账号获取用户信息
func GetUserByAccount(account string) (user *livemodel.User, err error) {
	req := contract.NewApiRequest()
	req.RequestData = account
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/getuserbyaccount", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &user)
	return user, err
}

// GetUserByUID 通过UID账号获取用户信息
func GetUserByUID(uid int64) (user *livemodel.User, err error) {
	req := contract.NewApiRequest()
	req.RequestData = uid
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/getuserbyuid", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &user)
	return user, err
}

// AddUser 通过账号和UID自动注册用户
func AddUser(account string, nickName string, uid int64) (userId int, err error) {
	req := contract.NewApiRequest()
	userReq := AddUserRequest{
		Account:  account,
		UID:      uid,
		NickName: nickName,
		Source:   _const.UserSource_PC,
	}
	req.RequestData = userReq
	response, err := PostNoCache(config.CurrentConfig.LiveWebApiHost+"/api/live/adduser", req)
	if err != nil {
		return 0, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return 0, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return 0, nil
	}
	userId = int(response.Message.(float64))
	return userId, err
}
