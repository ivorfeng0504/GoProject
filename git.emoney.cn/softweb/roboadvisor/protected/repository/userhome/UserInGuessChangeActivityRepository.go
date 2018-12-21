package repository

import (
	sysSql "database/sql"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strconv"
	"time"
)

type UserInGuessChangeActivityRepository struct {
	repository.BaseRepository
}

func NewUserInGuessChangeActivityRepository(conf *protected.ServiceConfig) *UserInGuessChangeActivityRepository {
	repo := &UserInGuessChangeActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetUserInGuessChangeActivity 获取用户指定一期的参与记录
func (repo *UserInGuessChangeActivityRepository) GetUserInGuessChangeActivity(userInfoId int, issueNumber string) (userInGuessChangeActivity *userhome_model.UserInGuessChangeActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInGuessChangeActivity] WHERE [UserInfoId]=? AND [IssueNumber]=?"
	userInGuessChangeActivity = new(userhome_model.UserInGuessChangeActivity)
	err = repo.FindOne(userInGuessChangeActivity, sql, userInfoId, issueNumber)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return userInGuessChangeActivity, err
}

// GetUserInGuessChangeActivityList 获取用户最新的N条参与记录 如果TOP小于等于0 则返回全部
func (repo *UserInGuessChangeActivityRepository) GetUserInGuessChangeActivityList(userInfoId int, top int) (userInGuessChangeActivityList []*userhome_model.UserInGuessChangeActivity, err error) {
	sql := "SELECT %s * FROM [UserHome_UserInGuessChangeActivity] WHERE [UserInfoId]=? ORDER BY [IssueNumber] DESC"
	if top > 0 {
		sql = fmt.Sprintf(sql, "TOP "+strconv.Itoa(top))
	}
	err = repo.FindList(&userInGuessChangeActivityList, sql, userInfoId)
	return userInGuessChangeActivityList, err
}

// GetUserInGuessChangeActivityListCurrentWeek 获取用户指定一周的参与记录
func (repo *UserInGuessChangeActivityRepository) GetUserInGuessChangeActivityListCurrentWeek(userInfoId int, date time.Time) (userInGuessChangeActivityList []*userhome_model.UserInGuessChangeActivity, err error) {
	weekDay := int(date.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	startTime := date.AddDate(0, 0, -weekDay+1).Format("20060102")
	endTime := date.AddDate(0, 0, 7-weekDay).Format("20060102")

	sql := "SELECT * FROM [UserHome_UserInGuessChangeActivity] WHERE [UserInfoId]=? AND [IssueNumber]>=? AND [IssueNumber]<=? ORDER BY [IssueNumber] DESC"
	err = repo.FindList(&userInGuessChangeActivityList, sql, userInfoId, startTime, endTime)
	return userInGuessChangeActivityList, err
}

// GetUserInGuessChangeActivityAwardList 获取用户最新的N条奖品记录 如果TOP小于等于0 则返回全部
func (repo *UserInGuessChangeActivityRepository) GetUserInGuessChangeActivityAwardList(userInfoId int, top int) (userInGuessChangeActivityList []*userhome_model.UserInGuessChangeActivity, err error) {
	sql := "SELECT %s * FROM [UserHome_UserInGuessChangeActivity] WHERE [UserInfoId]=? AND [AwardId]>0 ORDER BY [IssueNumber] DESC"
	if top > 0 {
		sql = fmt.Sprintf(sql, "TOP "+strconv.Itoa(top))
	}
	err = repo.FindList(&userInGuessChangeActivityList, sql, userInfoId)
	return userInGuessChangeActivityList, err
}

// GetGuessTotal 获取指定一期的竞猜统计
func (repo *UserInGuessChangeActivityRepository) GetGuessTotal(issueNumber string) (upCount int64, downCount int64, err error) {
	if len(issueNumber) == 0 {
		return
	}
	sql := "SELECT COUNT(1) AS TCount FROM [UserHome_UserInGuessChangeActivity] WHERE [IssueNumber]=? AND [Result]=1"
	sql += " UNION ALL "
	sql += "SELECT COUNT(1) AS TCount FROM [UserHome_UserInGuessChangeActivity] WHERE [IssueNumber]=? AND [Result]=-1"
	var result []*TotalCount
	err = repo.FindList(&result, sql, issueNumber, issueNumber)
	if err != nil || len(result) != 2 {
		return upCount, downCount, err
	}
	upCount = result[0].TCount
	downCount = result[1].TCount
	return upCount, downCount, err
}

// InsertUserInGuessChangeActivity 新增用户竞猜记录
func (repo *UserInGuessChangeActivityRepository) InsertUserInGuessChangeActivity(issueNumber string, userInfoId int, uid int64, nickName string, result int, activityCycle string) (err error) {
	record, err := repo.GetUserInGuessChangeActivity(userInfoId, issueNumber)
	if err != nil {
		return err
	}
	if record != nil {
		return nil
	}
	sql := "INSERT INTO [UserHome_UserInGuessChangeActivity] ([UID],[UserInfoId],[NickName],[IssueNumber],[Result],[ResultDesc],[ActivityCycle],[CreateTime])VALUES(?,?,?,?,?,?,?,GETDATE())"
	resultDesc := ""
	if result == 1 {
		resultDesc = "看涨"
	} else if result == -1 {
		resultDesc = "看跌"
	}
	_, err = repo.Insert(sql, uid, userInfoId, nickName, issueNumber, result, resultDesc, activityCycle)
	return err
}

// PublishUserGuessChangeResult 公布用户竞猜结果
func (repo *UserInGuessChangeActivityRepository) PublishUserGuessChangeResult(issueNumber string, result int) (count int64, err error) {
	sql := "UPDATE [UserHome_UserInGuessChangeActivity] SET [IsGuessed]=(CASE WHEN [Result]=? THEN 1 ELSE 0 END),[IsPublish]=1,[PublishTime]=GETDATE() WHERE [IssueNumber]=? "
	count, err = repo.Update(sql, result, issueNumber)
	return count, err
}

// UpdateUserGuessChangeAward 发放奖品
func (repo *UserInGuessChangeActivityRepository) UpdateUserGuessChangeAward(startIssueNumber string, endIssueNumber string, awardId int64, awardName string, reportUrl string) (count int64, err error) {
	sql := `UPDATE [UserHome_UserInGuessChangeActivity]  SET [AwardId]=?,[AwardName]=?,[ReportUrl]=? WHERE ([ReportUrl] IS NULL OR [ReportUrl]='') AND [UserInGuessChangeActivityId] IN (
SELECT [UserInGuessChangeActivityId] FROM [UserHome_UserInGuessChangeActivity] AS TB1 WHERE [UserInGuessChangeActivityId] IN(
SELECT TOP 1 [UserInGuessChangeActivityId] FROM [UserHome_UserInGuessChangeActivity] WHERE [IssueNumber]>=? AND [IssueNumber]<=? AND [IsGuessed]=1 AND [IsPublish]=1 AND [UserInfoId]=TB1.[UserInfoId]
AND [UserInfoId] IN(
SELECT [UserInfoId] FROM
(SELECT [UserInfoId],COUNT(1) AS [GuessedCount] FROM [UserHome_UserInGuessChangeActivity] WHERE [IssueNumber]>=? AND [IssueNumber]<=? AND [IsGuessed]=1 AND [IsPublish]=1 GROUP BY [UserInfoId]) AS TB2
WHERE [GuessedCount]>=3
)
ORDER BY [IssueNumber] DESC
)
)`
	count, err = repo.Update(sql, awardId, awardName, reportUrl, startIssueNumber, endIssueNumber, startIssueNumber, endIssueNumber)
	return count, err
}

// GetUserJoinCount 获取用户参与次数
func (repo *UserInGuessChangeActivityRepository) GetUserJoinCount(userInfoId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [UserHome_UserInGuessChangeActivity] WHERE [UserInfoId]=?"
	return repo.Count(sql, userInfoId)
}

type TotalCount struct {
	TCount int64
}
