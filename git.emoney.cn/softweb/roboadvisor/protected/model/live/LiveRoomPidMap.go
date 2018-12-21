package live

import "github.com/devfeel/mapper"

//记录RoomId与Pid之间的映射关系
//物流这边只有PID没有RoomId
type LiveRoomPidMap struct {
	//主键Id
	LiveRoomPidMapId int
	//直播间Id
	RoomId int
	//映射的PID
	PID int
	//备注
	Remark string
	//创建时间
	CreateTime mapper.JSONTime
	//是否删除
	IsDelete int
}
