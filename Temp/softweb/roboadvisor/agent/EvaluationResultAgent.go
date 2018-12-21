package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/stockapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"mime/multipart"
)

// GetEvaluationResult 查看当前用户的评测结果
func GetEvaluationResult(uid string) (result *myoptional_model.EvaluationResult, err error) {
	req := contract.NewApiRequest()
	req.RequestData = uid
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/evaluation/getresult", req)
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

// SubmitResult 提交评测
func SubmitResult(result myoptional_model.EvaluationResult) (err error) {
	req := contract.NewApiRequest()
	req.RequestData = result
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/evaluation/submitresult", req)
	if err != nil {
		return errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return errors.New(response.RetMsg)
	}
	return nil
}

// AnalyzeStock 上传图片，识别股票
func AnalyzeStock(file multipart.File) (stockList []*userhome_model.StockInfo, err error) {
	//识别图片中的股票代码
	stocks, err := stockapi.DiscernStockCode(file)
	if err != nil {
		return nil, err
	}
	if stocks == nil || len(stocks) == 0 {
		err = errors.New("未识别到股票代码！")
		return nil, err
	}

	//获取码表
	stockDict, err := stockapi.GetStockListDictCache()
	if err != nil || stockDict == nil {
		err = errors.New("获取码表异常！")
		return nil, err
	}

	myStockList := make(map[string]int)
	for _, stockCode := range stocks {
		stockName, success := stockDict[stockCode]
		if success == false {
			//如果未获取到股票名称，则抛弃
			continue
		}
		_, isExist := myStockList[stockCode]
		if isExist {
			//如果股票重复，则抛弃
			continue
		}
		stockList = append(stockList, &userhome_model.StockInfo{
			StockCode: stockCode,
			StockName: stockName,
		})
	}
	return stockList, nil
}

// GetStrategyList 获取推荐策略及股票列表
func GetStrategyList() (result []*StrategyInfo, err error) {
	req := contract.NewApiRequest()
	var strateKeyList []*StrategyInfo
	//高成长+高盈利
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "1",
		StrategyName:    "优质企业",
		StrategyKeyList: []string{"80019", "80021"},
		Top:             2,
		ReadCacheMode:   2,
	})
	//深跌回弹+趋势顶低
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "2",
		StrategyName:    "技术形态好",
		StrategyKeyList: []string{"100006", "100013"},
		Top:             2,
		ReadCacheMode:   2,
	})
	//资金博弈+大单比率
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "3",
		StrategyName:    "大资金进场",
		StrategyKeyList: []string{"100014", "100011"},
		Top:             2,
		ReadCacheMode:   2,
	})
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:    "4",
		StrategyName:  "大咖热议",
		ReadCacheMode: 2,
	})
	req.RequestData = strateKeyList
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/evaluation/getstrategylist", req)
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

type StrategyInfo struct {
	StrategyId      string
	StrategyName    string
	StrategyKeyList []string
	StockList       []*userhome_model.StockInfo
	//每个策略取前几个 0则全部
	Top int
	//0则不读取任何缓存  1 先读取缓存，缓存不存在则读取接口 2 先读取接口，如果接口中无数据，则读取缓存
	ReadCacheMode int
}

type StrategyInfoResult struct {
	StrategyId   string
	StrategyName string
	StockList    []*userhome_model.StockInfo
}
