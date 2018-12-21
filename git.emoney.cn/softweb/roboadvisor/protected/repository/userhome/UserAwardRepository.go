package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type UserAwardRepository struct {
	repository.BaseRepository
}

func NewUserAwardRepository(conf *protected.ServiceConfig) *UserAwardRepository {
	repo := &UserAwardRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetUserAwardList 获取用户的奖品列表 state大于0 则根据状态筛选
func (repo *UserAwardRepository) GetUserAwardList(userInfoId int, state int) (userAwardList []*userhome_model.UserAward, err error) {
	sql := "SELECT * FROM [UserHome_UserAward] WITH(NOLOCK) WHERE [UserInfoId]=?"
	if state > 0 {
		sql = "SELECT * FROM [UserHome_UserAward] WITH(NOLOCK) WHERE [UserInfoId]=? AND [State]=?"
		err = repo.FindList(&userAwardList, sql, userInfoId, state)
	} else {
		err = repo.FindList(&userAwardList, sql, userInfoId)
	}
	return userAwardList, err
}

// GetUserAwardById 根据Id获取用户奖品记录
func (repo *UserAwardRepository) GetUserAwardById(userAwardId int64) (userAward *userhome_model.UserAward, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserAward] WITH(NOLOCK) WHERE [UserAwardId]=?"
	err = repo.FindOne(&userAward, sql, userAwardId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return userAward, err
}

// InsertUserAward 插入用户奖品记录
func (repo *UserAwardRepository) InsertUserAward(userInfoId int, awardName string, awardId int64, introduceVideo string, activityId int64, activityName string, state int, awardImg string, avaDay int, awardType int) (id int64, err error) {
	sql := "INSERT INTO [UserHome_UserAward]([UserInfoId],[AwardName],[AwardId],[CreateTime],[IntroduceVideo],[ActivityId],[ActivityName],[State],[AwardImg],[AvailableDay],[AwardType])VALUES(?,?,?,GETDATE(),?,?,?,?,?,?,?)"
	id, err = repo.Insert(sql, userInfoId, awardName, awardId, introduceVideo, activityId, activityName, state, awardImg, avaDay, awardType)
	return id, err
}
