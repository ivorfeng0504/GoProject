package repository

import (
	sysSql "database/sql"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type GuessMarketActivityRepository struct {
	repository.BaseRepository
}

func NewGuessMarketActivityRepository(conf *protected.ServiceConfig) *GuessMarketActivityRepository {
	repo := &GuessMarketActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetGuessMarketActivityList 查询最新的N条开奖信息，top小于等于0则查询所有
func (repo *GuessMarketActivityRepository) GetGuessMarketActivityList(top int) (resultList []*userhome_model.GuessMarketActivity, err error) {
	sql := "SELECT * FROM [UserHome_GuessMarketActivity] WITH(NOLOCK) WHERE [IsPublish]=1 ORDER BY [PublishTime] DESC"
	if top > 0 {
		sql = fmt.Sprintf("SELECT TOP %d * FROM [UserHome_GuessMarketActivity] WITH(NOLOCK) WHERE [IsPublish]=1 ORDER BY [PublishTime] DESC", top)
	}
	err = repo.FindList(&resultList, sql)
	return resultList, err
}

// GetGuessMarketActivityByIssueNumber 根据期号获取开奖信息
func (repo *GuessMarketActivityRepository) GetGuessMarketActivityByIssueNumber(issueNum string) (result *userhome_model.GuessMarketActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_GuessMarketActivity] WITH(NOLOCK) WHERE [IssueNumber]=?"
	result = new(userhome_model.GuessMarketActivity)
	err = repo.FindOne(result, sql, issueNum)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// GetCurrentBeginingActivity 根据当前进行中的最新一期活动
func (repo *GuessMarketActivityRepository) GetCurrentBeginingActivity() (result *userhome_model.GuessMarketActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_GuessMarketActivity] WITH(NOLOCK) WHERE [BeginTime]<GETDATE() AND [EndTime]>GETDATE()"
	result = new(userhome_model.GuessMarketActivity)
	err = repo.FindOne(result, sql)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// GetCurrentFinishedActivity 获取当前已经结束的最后一期活动
func (repo *GuessMarketActivityRepository) GetCurrentFinishedActivity() (result *userhome_model.GuessMarketActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_GuessMarketActivity] WITH(NOLOCK) WHERE [EndTime]<GETDATE() ORDER BY [EndTime] DESC"
	result = new(userhome_model.GuessMarketActivity)
	err = repo.FindOne(result, sql)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// InsertActivity 插入新一期的活动
func (repo *GuessMarketActivityRepository) InsertActivity(issueNumber string, beginTime string, endTime string) (id int64, err error) {
	sql := "INSERT INTO [UserHome_GuessMarketActivity]([IssueNumber],[CreateTime],[BeginTime],[EndTime])VALUES(?,GETDATE(),?,?)"
	id, err = repo.Insert(sql, issueNumber, beginTime, endTime)
	return id, err
}

// PublishActivityResult 公布开奖号码
func (repo *GuessMarketActivityRepository) PublishActivityResult(issueNumber string, luckNum string) (n int64, err error) {
	sql := "UPDATE [UserHome_GuessMarketActivity] SET [Result]=?, [PublishTime]=GETDATE(),[IsPublish]=1 WHERE [IssueNumber]=?"
	n, err = repo.Update(sql, luckNum, issueNumber)
	return n, err
}
