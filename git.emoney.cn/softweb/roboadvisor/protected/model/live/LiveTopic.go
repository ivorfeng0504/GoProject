package live

import (
	"github.com/devfeel/mapper"
)

//直播主题
type LiveTopic struct {
	// 唯一标识
	Id int

	//直播室ID
	LId int

	//话题内容
	Topic string

	//发表内容
	Content string

	//用户ID
	UserId int

	//用户ID
	AdminUserId int64

	//通行证ID
	PassportId int64

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime
}
