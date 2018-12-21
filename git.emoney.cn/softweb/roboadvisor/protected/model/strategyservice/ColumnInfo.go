package strategyservice

import "github.com/devfeel/mapper"

type ColumnInfo struct {
	ID int
	// 栏目名称
	ColumnName string

	// 栏目描述
	ColumnDesc string

	// 是否删除
	IsDeleted bool

	// 创建时间
	CreateTime mapper.JSONTime

	// 应用id
	AppID int
}
