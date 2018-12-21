package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetRecommendStockList 获取推荐策略及股票列表
func GetRecommendStockList() (result []*StrategyInfoResult, err error) {
	req := contract.NewApiRequest()
	var strateKeyList []*StrategyInfo
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "100005",
		StrategyName:    "底部量变",
		StrategyKeyList: []string{"100005"},
		ReadCacheMode:   1,
	})
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "100006",
		StrategyName:    "深跌回弹",
		StrategyKeyList: []string{"100006"},
		ReadCacheMode:   1,
	})
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "100013",
		StrategyName:    "趋势顶底",
		StrategyKeyList: []string{"100013"},
		ReadCacheMode:   1,
	})
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "100015",
		StrategyName:    "龙腾四海",
		StrategyKeyList: []string{"100015"},
		ReadCacheMode:   1,
	})
	strateKeyList = append(strateKeyList, &StrategyInfo{
		StrategyId:      "80019",
		StrategyName:    "高成长",
		StrategyKeyList: []string{"80019"},
		ReadCacheMode:   1,
	})
	req.RequestData = strateKeyList
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/evaluation/getstrategylist", req)
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
