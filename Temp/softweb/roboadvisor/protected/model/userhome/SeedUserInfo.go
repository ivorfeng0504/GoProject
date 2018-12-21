package model

import "github.com/devfeel/mapper"

// 种子用户信息记录表
type SeedUserInfo struct {
	//主键Id
	SeedUserInfoId int64
	//cid
	Cid string

	//账号
	Account string

	//创建时间
	CreateTime mapper.JSONTime

	//修改时间
	ModifyTime mapper.JSONTime

	//是否删除
	IsDeleted bool
}
