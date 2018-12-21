package dataapi

type StockFinanceData struct {
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//最新业绩预告-类型
	PerformanceType string
	//最新业绩预告-范围
	PerformanceTypeRange string
	//本期营业收入收增长率
	OperatingRevenueYoY string
	//本期净利润增长率
	NetProfitYoY string
	//本期净资产收益率
	ROE string
	//本期毛利率
	SalesGrossMargin string
	//市盈率TTM
	PETTM string
	//最新业绩预告的报告期
	ReportPeriod string
}
