package model

type BoundAccountAdd_Response struct {
	//绑定账号cid
	AddInCID int64
	//绑定账号显示名
	AddInShowName string
	//绑定账号类型
	AddInType int
	//绑定后DataID
	NewCurDID int64
	//老的parentID
	AddInOldPID int64
	//返回消息
	RetMsg string

	RetCode int
}
