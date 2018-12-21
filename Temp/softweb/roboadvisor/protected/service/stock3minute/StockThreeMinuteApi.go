package stock3minute

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetStockListInfo 获取个股三分钟股票列表
func GetStockListInfo() (stockList []*model.StockInfo, err error) {
	apiUrl := config.CurrentConfig.StockThreeMinuteAPI_GetStockListInfo
	if len(apiUrl) == 0 {
		global.InnerLogger.ErrorFormat(err, "GetStockListInfo 获取个股三分钟股票列表 接口地址配置不正确 configkey=StockThreeMinuteAPI_GetStockListInfo")
		return nil, errors.New("获取个股三分钟股票列表 接口地址配置不正确")
	}
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockListInfo 获取个股三分钟股票列表 接口调用异常apiUrl=%s", apiUrl)
		return nil, errReturn
	}
	global.InnerLogger.DebugFormat("GetStockListInfo 获取个股三分钟股票列表 请求地址为：%s  contentType=%s intervalTime=%d 结果为：%s ", apiUrl, contentType, intervalTime, body)
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockListInfo 获取个股三分钟股票列表 ApiGatewayResponse反序列化异常 body=%s", body)
		return nil, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		global.InnerLogger.ErrorFormat(err, "GetStockListInfo 获取个股三分钟股票列表 RetCode异常 body=%s RetCode=%d", body, apiGatewayResp.RetCode)
		return nil, err
	}
	err = _json.Unmarshal(apiGatewayResp.Message, &stockList)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockListInfo 获取个股三分钟股票列表 apiGatewayResp.Message反序列化异常 body=%s", body)
		return nil, err
	}
	return stockList, err
}

// GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息
func GetStockThreeMinuteInfo(stockCode string) (info *StockThreeMinuteInfo, err error) {
	apiUrl := config.CurrentConfig.StockThreeMinuteAPI_GetStockThreeMinuteInfo
	if len(apiUrl) == 0 {
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 接口地址配置不正确 configkey=StockThreeMinuteAPI_GetStockThreeMinuteInfo")
		return nil, errors.New("根据股票代码获取最新的个股三分钟信息 接口地址配置不正确")
	}
	apiUrl = fmt.Sprintf(apiUrl, stockCode)

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 接口调用异常 apiUrl=%s", apiUrl)
		return nil, errReturn
	}
	global.InnerLogger.DebugFormat("GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息【%s】 请求地址为：%s  contentType=%s intervalTime=%d 结果为：%s ", stockCode, apiUrl, contentType, intervalTime, body)
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 ApiGatewayResponse反序列化异常 body=%s apiUrl=%s", body, apiUrl)
		return nil, err
	}
	if apiGatewayResp.RetCode != 0 {
		err = errors.New(apiGatewayResp.RetMsg)
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 RetCode异常 body=%s RetCode=%d apiUrl=%s", body, apiGatewayResp.RetCode, apiUrl)
		return nil, err
	}
	stockResponse := new(GetStockThreeMinuteInfoResponse)
	err = _json.Unmarshal(apiGatewayResp.Message, &stockResponse)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 apiGatewayResp.Message反序列化异常 body=%s  apiUrl=%s", body, apiUrl)
		return nil, err
	}
	if len(stockResponse.Content) == 0 {
		return nil, nil
	}
	stockContent := new(GetStockThreeMinuteInfoContent)
	err = _json.Unmarshal(stockResponse.Content, &stockContent)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockThreeMinuteInfo 根据股票代码获取最新的个股三分钟信息 stockResponse.Content反序列化异常 body=%s  apiUrl=%s", body, apiUrl)
		return nil, err
	}
	if stockContent.BasicInfo == nil || stockContent.MinuteMainInfo == nil {
		return nil, nil
	}
	info = &StockThreeMinuteInfo{
		StockCode:      stockContent.BasicInfo.StockCode,
		StockName:      stockContent.BasicInfo.StockName,
		OverviewValue:  stockContent.MinuteMainInfo.OverviewValue,
		TypeCap:        GetTypeCapDesc(stockContent.BasicInfo.TypeCap),
		TypeStyle:      GetTypeStyleDesc(stockContent.BasicInfo.TypeStyle),
		LifecycleValue: GetLifecycleValueDesc(stockContent.MinuteMainInfo.LifecycleValue),
		ScoreTTM:       GetScoreTTMDesc(stockContent.BasicInfo.ScoreTTM),
		ScoreGrowing:   GetScoreGrowingDesc(stockContent.BasicInfo.ScoreGrowing),
		ScoreProfit:    GetScoreProfitDesc(stockContent.BasicInfo.ScoreProfit),
	}
	return info, err
}
