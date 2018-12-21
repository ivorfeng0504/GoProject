package model

import "time"

//股票信息
type StockInfo struct {
	//股票代码
	StockCode string
	//股票名称
	StockName string
	//创建时间
	CreateTime time.Time
}
