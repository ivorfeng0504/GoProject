package live

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type LiveRoomPidMapRepository struct {
	repository.BaseRepository
}

func NewLiveRoomPidMapRepository(conf *protected.ServiceConfig) *LiveRoomPidMapRepository {
	repo := &LiveRoomPidMapRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetMapList 获取所有直播间与PID之间的映射关系
func (repo *LiveRoomPidMapRepository) GetMapList() (mapList []*livemodel.LiveRoomPidMap, err error) {
	sql := "SELECT * FROM [LiveRoomPidMap]"
	err = repo.FindList(&mapList, sql)
	return mapList, err
}
