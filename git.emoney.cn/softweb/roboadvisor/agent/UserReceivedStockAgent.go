package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// InsertUserStock 用户领取股票
func InsertUserStock(userinfo *userhome_model.UserInfo) (userStock *userhome_model.UserReceivedStock, err error) {
	req := contract.NewApiRequest()
	requestData := InsertUserStockRequest{
		UserInfo:   userinfo,
		ActivityId: config.CurrentConfig.Activity_ReceiveStock,
	}
	req.RequestData = requestData
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/activity/insertuserstock", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return userStock, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userStock)
	return userStock, err
}

// GetUserStockToday 获取用户当日领取的股票
func GetUserStockToday(userInfoId int) (userStock *userhome_model.UserReceivedStock, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfoId
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/getuserstocktoday", req)
	if err != nil {
		return userStock, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return userStock, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return userStock, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userStock)
	return userStock, err
}

// GetUserStockHistory 获取用户领取的股票历史
func GetUserStockHistory(userInfoId int) (userStockList []*userhome_model.UserReceivedStock, err error) {
	req := contract.NewApiRequest()
	requestData := UserReceivedStockHistoryRequest{
		UserInfoId: userInfoId,
		Top:        10,
	}
	req.RequestData = requestData
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/getuserstockhistory", req)
	if err != nil {
		return userStockList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return userStockList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return userStockList, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &userStockList)
	return userStockList, err
}

type UserReceivedStockHistoryRequest struct {
	UserInfoId int
	Top        int
}

type InsertUserStockRequest struct {
	UserInfo   *userhome_model.UserInfo
	ActivityId int64
}
