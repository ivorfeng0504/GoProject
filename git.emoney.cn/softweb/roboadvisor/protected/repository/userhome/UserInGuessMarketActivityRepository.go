package repository

import (
	sysSql "database/sql"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strconv"
)

type UserInGuessMarketActivityRepository struct {
	repository.BaseRepository
}

func NewUserInGuessMarketActivityRepository(conf *protected.ServiceConfig) *UserInGuessMarketActivityRepository {
	repo := &UserInGuessMarketActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetUserInGuessMarketActivityList 查询用户最新的参与记录
func (repo *UserInGuessMarketActivityRepository) GetUserInGuessMarketActivityList(top int, userInfoId int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	sql := "SELECT * FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE UserInfoId=? ORDER BY [CreateTime] DESC"
	if top > 0 {
		sql = fmt.Sprintf("SELECT TOP %d * FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE UserInfoId=? ORDER BY [CreateTime] DESC", top)
	}
	err = repo.FindList(&resultList, sql, userInfoId)
	return resultList, err
}

// GetUserJoinCount 查询用户参与总数
func (repo *UserInGuessMarketActivityRepository) GetUserJoinCount(userInfoId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE UserInfoId=?"
	return repo.Count(sql, userInfoId)
}

// GetUserInGuessMarketActivityGuessedList 查询最新的获奖记录列表
func (repo *UserInGuessMarketActivityRepository) GetUserInGuessMarketActivityGuessedList(top int) (resultList []*userhome_model.UserInGuessMarketActivity, err error) {
	sql := "SELECT * FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE [IsGuess]=1 ORDER BY [PublishTime] DESC"
	if top > 0 {
		sql = fmt.Sprintf("SELECT TOP %d * FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE [IsGuess]=1 ORDER BY [PublishTime] DESC", top)
	}
	err = repo.FindList(&resultList, sql)
	return resultList, err
}

// GetUserInGuessMarketActivity 获取用户指定期号的猜数字记录
func (repo *UserInGuessMarketActivityRepository) GetUserInGuessMarketActivity(userInfoId int, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE [IssueNumber]=? AND [UserInfoId]=?"
	result = new(userhome_model.UserInGuessMarketActivity)
	err = repo.FindOne(result, sql, issueNum, userInfoId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// InsertUserInGuessMarketActivity 插入用户领数字记录
func (repo *UserInGuessMarketActivityRepository) InsertUserInGuessMarketActivity(uid int64, userInfoId int, nickName string, issueNum string, result string, guessMarketActivityId int64) (id int64, err error) {
	sql := "INSERT INTO [UserHome_UserInGuessMarketActivity]([UID],[UserInfoId],[CreateTime],[IssueNumber],[Result],[NickName],[GuessMarketActivityId])VALUES(?,?,GETDATE(),?,?,?,?)"
	id, err = repo.Insert(sql, uid, userInfoId, issueNum, result, nickName, guessMarketActivityId)
	return id, err
}

// GetWinner 获取获奖用户
func (repo *UserInGuessMarketActivityRepository) GetWinner(luckNum string, issueNum string) (result *userhome_model.UserInGuessMarketActivity, err error) {
	luckNumInt, err := strconv.Atoi(luckNum)
	if err != nil {
		return nil, err
	}
	sql := "SELECT TOP 1 *,ABS(CAST([Result] AS INT)-?) AS DiffValue FROM [UserHome_UserInGuessMarketActivity] WITH(NOLOCK) WHERE [Result] IS NOT NULL AND [IssueNumber]=? ORDER BY DiffValue ASC,[CreateTime] ASC"
	result = new(userhome_model.UserInGuessMarketActivity)
	err = repo.FindOne(result, sql, luckNumInt, issueNum)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// UpdateUserAward 更新用户奖品
func (repo *UserInGuessMarketActivityRepository) UpdateUserAward(userInGuessMarketActivityId int64, awardId int64, awardName string) (n int64, err error) {
	sql := "UPDATE [UserHome_UserInGuessMarketActivity] SET [IsGuess]=1, [PublishTime]=GETDATE(),[AwardId]=?,[AwardName]=? WHERE [UserInGuessMarketActivityId]=?"
	n, err = repo.Update(sql, awardId, awardName, userInGuessMarketActivityId)
	return n, err
}
