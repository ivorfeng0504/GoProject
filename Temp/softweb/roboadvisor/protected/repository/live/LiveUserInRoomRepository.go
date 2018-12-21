package live

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type LiveUserInRoomRepository struct {
	repository.BaseRepository
}

func NewLiveUserInRoomRepository(conf *protected.ServiceConfig) *LiveUserInRoomRepository {
	repo := &LiveUserInRoomRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// AddLiveUserInRoom 添加用户图文直播权限
func (repo *LiveUserInRoomRepository) AddLiveUserInRoom(mobile string, roomId int, day int, orderId string, source string) (id int64, err error) {
	expireTime := time.Now().Add(time.Hour * 24 * time.Duration(day))
	sql := "INSERT INTO [LiveUserInRoom] ([Mobile],[RoomId],[ExpireTime],[CreateTime],[OrderId],[Source]) VALUES(?,?,?,?,?,?)"
	id, err = repo.Insert(sql, mobile, roomId, expireTime, time.Now(), orderId, source)
	return id, err
}

// GetLiveUserInRoomList 根据账号获取用户的直播间权限
func (repo *LiveUserInRoomRepository) GetLiveUserInRoomList(mobile string) (rooms []*livemodel.LiveUserInRoom, err error) {
	sql := "SELECT * FROM [LiveUserInRoom] WHERE [Mobile]=? AND [ExpireTime]>GETDATE() AND [IsDelete]=0"
	err = repo.FindList(&rooms, sql, mobile)
	return rooms, err
}

// DeleteById 删除指定Id的数据（逻辑删除）
func (repo *LiveUserInRoomRepository) DeleteById(liveUserInRoomId int) (err error) {
	sql := "UPDATE [LiveUserInRoom] SET [IsDelete]=1,[DeleteTime]=GETDATE() WHERE [LiveUserInRoomId]=? AND [IsDelete]=0"
	_, err = repo.Update(sql, liveUserInRoomId)
	return err
}

// GetLiveUserInRoomById 根据Id获取用户的直播间权限
func (repo *LiveUserInRoomRepository) GetLiveUserInRoomById(liveUserInRoomId int) (room *livemodel.LiveUserInRoom, err error) {
	sql := "SELECT TOP 1 * FROM [LiveUserInRoom] WHERE [LiveUserInRoomId]=? AND [IsDelete]=0"
	room = &livemodel.LiveUserInRoom{}
	err = repo.FindOne(room, sql, liveUserInRoomId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return room, err
}
