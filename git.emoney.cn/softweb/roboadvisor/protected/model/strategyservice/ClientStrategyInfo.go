package strategyservice

import "github.com/devfeel/mapper"

//客户端策略信息
type ClientStrategyInfo struct {
	//主键Id
	ClientStrategyInfoId int64
	//栏目Id
	ColumnInfoId int
	//策略Id
	ClientStrategyId int
	//策略名称
	ClientStrategyName string
	// 是否删除
	IsDeleted bool
	// 创建时间
	CreateTime mapper.JSONTime
	//是否为顶级策略
	IsTop bool

	//以下为非数据库字段
	//栏目名称
	ColumnName string
	//子集
	Children []*ClientStrategyInfo
	//父级Id
	ParentId int
	//父级名称
	ParentName string
}
