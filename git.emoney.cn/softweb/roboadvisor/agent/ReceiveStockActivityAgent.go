package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// InsertStockPool 新增一个股票池
func InsertStockPool(stockList []*userhome_model.StockInfo) (err error) {
	req := contract.NewApiRequest()
	req.RequestData = stockList
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/activity/insertstockpool", req)
	if err != nil {
		return errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return errors.New(response.RetMsg)
	}
	return nil
}

// GetNewstStockPool 获取最新的股票池
func GetNewstStockPool(userInfoId int) (stock *userhome_model.ReceiveStockActivity, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/activity/getnewststockpool", req)
	if err != nil {
		return stock, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return stock, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return stock, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &stock)
	return stock, nil
}
