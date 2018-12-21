package scmapi

type QueryOrderProdListByParamsResponse struct {
	//返回代码(1成功，其它失败)
	Code int

	//返回信息
	Msg string

	//数据
	Data []*OrderProdInfo
}
