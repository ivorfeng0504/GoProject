package model

import (
	"github.com/devfeel/mapper"
)

type BannerInfo struct{

	//唯一标识
	ID int

	//标题
	Title string

	//banner指向（0:内部文章/视频  1:外部链接）
	BannerPointType int

	NewsID int

	//外部链接
	OutLinkUrl string

	//排序位置
	SortIndex int

	//banner图片
	BannerImg string

	//策略权限ID（全部：0   1,2,3其他权限）
	StrategyID string

	//是否删除
	IsDeleted bool

	//创建用户
	CreateUser string

	//创建时间
	CreateTime mapper.JSONTime

	//最后更新时间
	LastModifyTime mapper.JSONTime
}
