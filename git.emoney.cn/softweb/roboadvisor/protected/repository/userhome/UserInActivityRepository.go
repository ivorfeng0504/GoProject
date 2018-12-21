package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type UserInActivityRepository struct {
	repository.BaseRepository
}

func NewUserInActivityRepository(conf *protected.ServiceConfig) *UserInActivityRepository {
	repo := &UserInActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetUserInActivity 查询用户参与指定活动的记录
func (repo *UserInActivityRepository) GetUserInActivity(userInfoId int, activityId int64) (record *userhome_model.UserInActivity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInActivity] WITH(NOLOCK) WHERE [ActivityId]=? AND [UserInfoId]=?"
	record = new(userhome_model.UserInActivity)
	err = repo.FindOne(&record, sql, activityId, userInfoId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return record, err
}

// GetUserInActivityList 查询用户所有参与活动记录
func (repo *UserInActivityRepository) GetUserInActivityList(userInfoId int) (recordList []*userhome_model.UserInActivity, err error) {
	sql := "SELECT * FROM [UserHome_UserInActivity] WITH(NOLOCK) WHERE [UserInfoId]=?"
	err = repo.FindList(&recordList, sql, userInfoId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return recordList, err
}

// InsertUserInActivity 插入用户活动参与记录
func (repo *UserInActivityRepository) InsertUserInActivity(userInfoId int, activityId int64) (id int64, err error) {
	activity, err := repo.GetUserInActivity(userInfoId, activityId)
	if err != nil {
		return 0, err
	}
	if activity != nil {
		return activity.UserInActivityId, nil
	}
	sql := "INSERT INTO [UserHome_UserInActivity] ([ActivityId],[UserInfoId],[CreateTime])VALUES(?,?,GETDATE())"
	id, err = repo.Insert(sql, activityId, userInfoId)
	return id, err
}
