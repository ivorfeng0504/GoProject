package live

import (
	"github.com/devfeel/mapper"
)

// LiveUserInRoom 用户拥有的直播间权限
type LiveUserInRoom struct {

	//自增ID
	LiveUserInRoomId int64

	//手机号
	Mobile string

	//直播间号
	RoomId int

	//权限过期时间
	ExpireTime mapper.JSONTime

	//创建时间
	CreateTime mapper.JSONTime

	//物流订单Id
	OrderId string

	//来源记录appid
	Source string

	//是否删除
	IsDelete int

	//删除时间
	DeleteTime mapper.JSONTime
}
