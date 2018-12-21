package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
)

type StrategyInfoRepository struct {
	repository.BaseRepository
}

func NewStrategyInfoRepository(conf *protected.ServiceConfig) *StrategyInfoRepository {
	repo := &StrategyInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}


// UpdateClickNum 更新点击数
func (repo *StrategyInfoRepository) UpdateClickNum(newsId int64, clickNum int64) (err error) {
	sql := "UPDATE [ExpertNews_StrategyInfo] SET [ClickNum]=? WHERE [ID]=? and ([ClickNum] IS NULL OR [ClickNum]<?)"
	_, err = repo.Update(sql, clickNum, newsId, clickNum)
	return err
}

// UpdateVideoPlayNum 更新点播数
func (repo *StrategyInfoRepository) UpdateVideoPlayNum(newsId int64, playNum int64) (err error) {
	sql := "UPDATE [ExpertNews_StrategyInfo] SET [VideoPlayNum]=? WHERE [ID]=? and ([VideoPlayNum] IS NULL OR [VideoPlayNum]<?)"
	_, err = repo.Update(sql, playNum, newsId, playNum)
	return err
}

// 根据直播ID获取对应策略信息
func (repo *StrategyInfoRepository) GetStrategyListByLiveID(liveID int) (strategyList []*expertnews.ExpertNews_StrategyInfo,err error) {
	sql := "SELECT * FROM [ExpertNews_StrategyInfo] where LiveID=?"
	err = repo.FindList(&strategyList, sql, liveID)
	return strategyList, err
}

// GetLatestNewsByExpertStrategyID 获取最新一条专家策略资讯
func (repo *StrategyInfoRepository) GetLatestNewsByExpertStrategyID(strategyID int) (newsList []*model.NewsInfo,err error) {
	sql:="SELECT TOP 1 * FROM NewsInfo WHERE ExpertStrategyID=? AND IsDeleted = 0 ORDER BY LastModifyTime DESC"
	err = repo.FindList(&newsList, sql, strategyID)
	return newsList, err
}


