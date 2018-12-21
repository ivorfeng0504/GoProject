package dataapi

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// PROD_ZYDS_StockBar_Quote 获取股票PE1数据
func PROD_ZYDS_StockBar_Quote(stockCode string) (pe1 string, err error) {
	apiUrl := `http://dataapi.emoney.cn/platformapi/indicator/execcondition?token=421E7810BD348D10FFB6B142DED5A5F3159A7A5C&condition=PROD_ZYDS_StockBar_Quote(0,%s)`
	apiUrl = fmt.Sprintf(apiUrl, stockCode)
	result, err := GetDataApi(apiUrl)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "PROD_ZYDS_StockBar_Quote 获取股票PE1数据 异常 请求地址为%s", apiUrl)
		return pe1, err
	}
	if result == nil || len(result) <= 2 {
		err = errors.New("未获取到PE1数据")
		global.InnerLogger.ErrorFormat(err, "PROD_ZYDS_StockBar_Quote 获取股票PE1数据 异常 未获取到PE1数据 请求地址为%s result=%s", apiUrl, _json.GetJsonString(result))
		return pe1, err
	}
	pe1 = result[2][2]
	return pe1, err
}

// PROD_ZYDS_StockBar_Finance 获取股票财务数据
func PROD_ZYDS_StockBar_Finance(stockCode string) (data *StockFinanceData, err error) {
	apiUrl := `http://dataapi.emoney.cn/platformapi/indicator/execcondition?token=421E7810BD348D10FFB6B142DED5A5F3159A7A5C&condition=PROD_ZYDS_StockBar_Finance(0,%s)`
	apiUrl = fmt.Sprintf(apiUrl, stockCode)
	result, err := GetDataApi(apiUrl)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "PROD_ZYDS_StockBar_Finance 获取股票财务数据 异常 请求地址为%s", apiUrl)
		return nil, err
	}
	if result == nil || len(result) <= 2 || len(result[2]) < 10 {
		err = errors.New("未获取到股票财务数据")
		global.InnerLogger.ErrorFormat(err, "PROD_ZYDS_StockBar_Finance 获取股票财务数据 异常 未获取到股票财务数据 请求地址为%s result=%s", apiUrl, _json.GetJsonString(result))
		return nil, err
	}
	data = &StockFinanceData{
		StockCode:            stockCode,
		StockName:            result[2][1],
		PerformanceType:      result[2][8],
		PerformanceTypeRange: result[2][9],
		OperatingRevenueYoY:  result[2][3],
		NetProfitYoY:         result[2][4],
		ROE:                  result[2][6],
		SalesGrossMargin:     result[2][5],
		ReportPeriod:         result[2][7],
	}
	return data, err
}

// GetStockQuoteAndFinance 获取股票财务数据与PE1数据（聚合PROD_ZYDS_StockBar_Quote接口与PROD_ZYDS_StockBar_Finance接口）
func GetStockQuoteAndFinance(stockCode string) (data *StockFinanceData, err error) {
	data, err = PROD_ZYDS_StockBar_Finance(stockCode)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockQuoteAndFinance 获取股票财务数据与PE1数据 异常 PROD_ZYDS_StockBar_Finance执行异常 stockCode=%s", stockCode)
		return nil, err
	}
	if data == nil {
		err = errors.New("获取到的数据为空")
		global.InnerLogger.ErrorFormat(err, "GetStockQuoteAndFinance 获取股票财务数据与PE1数据 异常 PROD_ZYDS_StockBar_Finance获取到的数据为空 stockCode=%s", stockCode)
		return nil, err
	}
	pe1, err := PROD_ZYDS_StockBar_Quote(stockCode)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "GetStockQuoteAndFinance 获取股票财务数据与PE1数据 异常 PROD_ZYDS_StockBar_Quote执行异常 stockCode=%s", stockCode)
		return nil, err
	}
	data.PETTM = pe1
	return data, err
}
