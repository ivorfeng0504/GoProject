package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type ActivityRepository struct {
	repository.BaseRepository
}

func NewActivityRepository(conf *protected.ServiceConfig) *ActivityRepository {
	repo := &ActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetActivityList 查询所有活动信息 state见ActivityState.go
func (repo *ActivityRepository) GetActivityList(state int) (activityList []*userhome_model.Activity, err error) {
	sql := "SELECT * FROM [UserHome_Activity] WHERE [IsDeleted]=0 ORDER BY [BeginTime] DESC"
	switch state {
	case _const.ActivityState_NotBegin:
		sql = "SELECT * FROM [UserHome_Activity] WHERE [IsDeleted]=0 AND [BeginTime]>GETDATE() ORDER BY [BeginTime] DESC"
		break
	case _const.ActivityState_Beginning:
		sql = "SELECT * FROM [UserHome_Activity] WHERE [IsDeleted]=0 AND [BeginTime]<GETDATE() AND [EndTime]>GETDATE()  ORDER BY [BeginTime] DESC"
		break
	case _const.ActivityState_Finish:
		sql = "SELECT * FROM [UserHome_Activity] WHERE [IsDeleted]=0 AND [EndTime]<GETDATE() ORDER BY [BeginTime] DESC"
		break
	}
	err = repo.FindList(&activityList, sql)
	return activityList, err
}

// GetActivityById 根据活动Id获取活动信息
func (repo *ActivityRepository) GetActivityById(activityId int64) (activity *userhome_model.Activity, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_Activity] WHERE [IsDeleted]=0 AND [ActivityId]=?"
	activity = new(userhome_model.Activity)
	err = repo.FindOne(activity, sql, activityId)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return activity, err
}
