package live

import (
	"github.com/devfeel/mapper"
)

type MaskWord struct{

	// 唯一标识
	Id int

	//屏蔽字名称
	MaskName string

	//发布时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime

	//是否删除
	IsDelete int

	//删除时间
	DeleteTime mapper.JSONTime
}