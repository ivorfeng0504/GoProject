package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"time"
	"git.emoney.cn/softweb/roboadvisor/const"
)

// GetStockNewsInfo 获取股票相关资讯
func GetStockNewsInfo(stockList []string) (result []*myoptional_model.StockNewsInformation, err error) {
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	requestData := StockNewsInformationRequest{
		Top:       8,
		StartTime: time.Now().AddDate(0, 0, -20).Format("2006-01-02"),
		StockList: stockList,
	}
	req.RequestData = requestData
	response, err := Post(config.CurrentConfig.StockTalkWebApiHost+"/api/stocknews/getstocknewsinfo", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetStockNewsInfoByHSet  根据股票代码获取相关资讯（直接读HSet缓存）
func GetStockNewsInfoByHSet(stockList []string) (result []*myoptional_model.StockNewsInformation, err error) {
	if stockList == nil || len(stockList) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	requestData := StockNewsInformationRequest{
		Top:       30,
		StockList: stockList,
		PerTop:    _const.StockNewsInformationPerTop,
	}
	req.RequestData = requestData
	response, err := Post(config.CurrentConfig.StockTalkWebApiHost+"/api/stocknews/getstocknewsinfobyhset", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

type StockNewsInformationRequest struct {
	StartTime string
	StockList []string
	//返回的数据总数
	Top int
	//每只股票缓存的数目，需要与Task中的一致
	PerTop int
}
