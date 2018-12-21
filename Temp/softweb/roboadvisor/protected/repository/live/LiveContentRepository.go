package live

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type LiveContentRepository struct {
	repository.BaseRepository
}

func NewLiveContentRepository(conf *protected.ServiceConfig) *LiveContentRepository {
	repo := &LiveContentRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetContentList 根据topicId查询直播内容
func (repo *LiveContentRepository) GetContentList(topicId int) (contentList []*livemodel.LiveContent, err error) {
	sql := "SELECT * FROM [LiveContent] WITH(NOLOCK) WHERE [IsDelete]=0 AND [TId]=? ORDER BY [CreateTime] ASC"
	err = repo.FindList(&contentList, sql, topicId)
	return contentList, err
}

// GetTopContent 根据topicId查询置顶的直播内容
func (repo *LiveContentRepository) GetTopContent(topicId int) (topContent *livemodel.LiveContent, err error) {
	sql := "SELECT TOP 1 * FROM [LiveContent] WITH(NOLOCK) WHERE [IsDelete]=0 AND [TId]=? AND [IsTop]=1 ORDER BY [CreateTime] DESC"
	topContent = new(livemodel.LiveContent)
	err = repo.FindOne(topContent, sql, topicId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return topContent, err
}
