package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type GuessChangeActivityRepository struct {
	repository.BaseRepository
}

func NewGuessChangeActivityRepository(conf *protected.ServiceConfig) *GuessChangeActivityRepository {
	repo := &GuessChangeActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetNewstGuessChangeActivity 获取最新一期的猜涨跌活动
func (repo *GuessChangeActivityRepository) GetNewstGuessChangeActivity() (activity *userhome_model.GuessChangeActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_GuessChangeActivity] WHERE [BeginTime]<=GETDATE() AND EndTime>=GETDATE() ORDER BY [IssueNumber] DESC"
	activity = new(userhome_model.GuessChangeActivity)
	err = repo.FindOne(activity, sql)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return activity, err
}

// GetGuessChangeActivityByIssueNumber 获取指定一期的猜涨跌活动
func (repo *GuessChangeActivityRepository) GetGuessChangeActivityByIssueNumber(issueNumber string) (activity *userhome_model.GuessChangeActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_GuessChangeActivity] WHERE [IssueNumber]=?"
	activity = new(userhome_model.GuessChangeActivity)
	err = repo.FindOne(activity, sql, issueNumber)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return activity, err
}

// InsertGuessChangeActivity 新增一个猜涨跌活动
func (repo *GuessChangeActivityRepository) InsertGuessChangeActivity(issueNumber string, beginTime time.Time, endTime time.Time, activityCycle string) (err error) {
	sql := "INSERT INTO [UserHome_GuessChangeActivity]([IssueNumber],[BeginTime],[EndTime],[ActivityCycle],[CreateTime])VALUES(?,?,?,?,GETDATE())"
	_, err = repo.Insert(sql, issueNumber, beginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), activityCycle)
	return err
}

// PublishGuessChangeActivity 公布竞猜结果
func (repo *GuessChangeActivityRepository) PublishGuessChangeActivity(issueNumber string, result int) (err error) {
	sql := "UPDATE [UserHome_GuessChangeActivity] SET [Result]=? , [IsPublish]=1,[PublishTime]=GETDATE() WHERE [IssueNumber]=? AND [IsPublish]=0"
	_, err = repo.Update(sql, result, issueNumber)
	return err
}
