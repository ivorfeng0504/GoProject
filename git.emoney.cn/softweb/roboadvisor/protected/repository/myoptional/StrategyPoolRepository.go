package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type StrategyPoolRepository struct {
	repository.BaseRepository
}

func NewStrategyPoolRepository(conf *protected.ServiceConfig) *StrategyPoolRepository {
	repo := &StrategyPoolRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertStrategyPool 插入一条股池信息
func (repo *StrategyPoolRepository) InsertStrategyPool(model *myoptional_model.StrategyPool) (strategyPoolId int64, err error) {
	sql := "INSERT INTO [StrategyPool] ([ClientStrategyId],[ClientStrategyName],[ParentId],[ParentName],[StockCode],[StockName],[IssueNumber],[CreateTime],[IsDeleted]) VALUES(?,?,?,?,?,?,?,GETDATE(),0)"
	strategyPoolId, err = repo.Insert(sql, model.ClientStrategyId, model.ClientStrategyName, model.ParentId, model.ParentName, model.StockCode, model.StockName, model.IssueNumber)
	return strategyPoolId, err
}
