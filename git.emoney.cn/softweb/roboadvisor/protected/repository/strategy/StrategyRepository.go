package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/strategy"
)

type StrategyRepository struct {
	repository.BaseRepository
}

func StrategyInfoRepository(conf *protected.ServiceConfig) *StrategyRepository {
	repo := &StrategyRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetStrategyInfoByStrategyId 根据strategyid获取资讯
func (repo *StrategyRepository) GetStrategyInfoBySId(strategyid int) (*model.StrategyInfo,error){
	sql := "SELECT * FROM [StrategyInfo] WITH(NOLOCK) WHERE StrategyID = ? AND IsDeleted=0"
	strategyinfo := new(model.StrategyInfo)
	err := repo.FindOne(strategyinfo, sql, strategyid)
	if err == sysSql.ErrNoRows {
		return strategyinfo, nil
	}
	return strategyinfo, err
}

// GetStrategyList 获取所有的策略列表
func (repo *StrategyRepository) GetStrategyList(strategylist []*model.StrategyInfo) ( err error) {
	sql := "SELECT * FROM [StrategyInfo] WITH(NOLOCK) WHERE IsDeleted=0"
	err = repo.FindList(&strategylist, sql)
	return  err
}
