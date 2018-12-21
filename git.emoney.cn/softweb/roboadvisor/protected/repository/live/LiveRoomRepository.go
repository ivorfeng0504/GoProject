package live

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type LiveRoomRepository struct {
	repository.BaseRepository
}

func NewLiveRoomRepository(conf *protected.ServiceConfig) *LiveRoomRepository {
	repo := &LiveRoomRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}


// GetLiveRoom 根据直播间Id获取直播间信息
func (repo *LiveRoomRepository) GetLiveRoom(roomId int) (*livemodel.LiveRoom, error) {
	sql := "SELECT TOP 1 * FROM [LiveRoom] WITH(NOLOCK) WHERE [Id]=? AND [IsDelete]=0"
	room := new(livemodel.LiveRoom)
	err := repo.FindOne(room, sql, roomId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return room, err
}