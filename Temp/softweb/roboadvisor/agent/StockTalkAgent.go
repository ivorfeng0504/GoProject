package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	myoptional_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"time"
)

// InsertStockTalk 插入一个微股吧评论
func InsertStockTalk(model myoptional_model.StockTalk) (err error) {
	req := contract.NewApiRequest()
	req.RequestData = model
	response, err := PostNoCache(config.CurrentConfig.StockTalkWebApiHost+"/api/stocktalk/insertstocktalk", req)
	if err != nil {
		return errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return errors.New(response.RetMsg)
	}
	return nil
}

// GetStockTalkListByDate 获取指定日期的微股吧评论
func GetStockTalkListByDate(isIncre bool, dateStr string) (stockTalkList []*myoptional_vmmodel.StockTalk, err error) {
	req := contract.NewApiRequest()
	requestData := StockTalkRequest{
		IsIncre: isIncre,
	}
	if len(dateStr) == 0 {
		requestData.CheckDate = false
	} else {
		requestData.Date, err = _time.ParseTime(dateStr)
		if err == nil {
			requestData.CheckDate = true
		}
	}
	req.RequestData = requestData
	response, err := Post(config.CurrentConfig.StockTalkWebApiHost+"/api/stocktalk/getstocktalklistbydate", req)
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
	err = _json.Unmarshal(jsonStr, &stockTalkList)
	return stockTalkList, err
}

// GetStockTalkListByPage 分页获取微股吧评论
func GetStockTalkListByPage(stockCodeList []string, pageIndex int, pageSize int) (stockTalkList []*myoptional_vmmodel.StockTalk, err error) {
	req := contract.NewApiRequest()
	if stockCodeList == nil || len(stockCodeList) == 0 || pageIndex < 0 || pageSize <= 0 || (len(stockCodeList) == 1 && len(stockCodeList[0]) == 0) {
		return nil, nil
	}
	if pageSize > 30 {
		return nil, errors.New("分页大小不能大于30")
	}
	requestData := StockTalkRequest{
		StockCodeList: stockCodeList,
		PageIndex:     pageIndex,
		PageSize:      pageSize,
	}
	req.RequestData = requestData
	response, err := Post(config.CurrentConfig.StockTalkWebApiHost+"/api/stocktalk/getstocktalklistbypage", req)
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
	err = _json.Unmarshal(jsonStr, &stockTalkList)
	return stockTalkList, err
}

type StockTalkRequest struct {
	//是否查询增量
	IsIncre bool
	//日期
	Date time.Time
	//是否传递了正确的日期
	CheckDate bool
	//分页大小
	PageSize int
	//分页索引 索引从0开始
	PageIndex int
	//股票列表
	StockCodeList []string
}
