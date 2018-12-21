package live

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type LiveTopicRepository struct {
	repository.BaseRepository
}

func NewLiveTopicRepository(conf *protected.ServiceConfig) *LiveTopicRepository {
	repo := &LiveTopicRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetTopic 根据直播间Id获取主题信息
func (repo *LiveTopicRepository) GetTopic(roomId int, date time.Time) (*livemodel.LiveTopic, error) {
	sql := "SELECT TOP 1 * FROM [LiveTopic] WITH(NOLOCK) WHERE [LId]=? AND DATEDIFF(dd,[CreateTime],?)=0 ORDER BY [CreateTime] DESC"
	topic := new(livemodel.LiveTopic)
	err := repo.FindOne(topic, sql, roomId, date.Format("2006-01-02"))
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return topic, err
}

// GetNewestTopic 根据直播间Id获取最新主题信息
func (repo *LiveTopicRepository) GetNewestTopic(roomId int) (*livemodel.LiveTopic, error) {
	sql := "SELECT TOP 1 * FROM [LiveTopic] WITH(NOLOCK) WHERE [LId]=? ORDER BY [CreateTime] DESC"
	topic := new(livemodel.LiveTopic)
	err := repo.FindOne(topic, sql, roomId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return topic, err
}
