package repository

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type AwardInActivityRepository struct {
	repository.BaseRepository
}

func NewAwardInActivityRepository(conf *protected.ServiceConfig) *AwardInActivityRepository {
	repo := &AwardInActivityRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

//获取指定活动的奖品列表 activityId活动Id  ignoreExpired是否过滤掉过期的奖品设置
func (repo *AwardInActivityRepository) GetActivityAwardListByActivityId(activityId int64, ignoreExpired bool) (awardList []*userhome_model.ActivityAward, err error) {
	if activityId <= 0 {
		return nil, errors.New("活动Id不能为空")
	}
	sql := `SELECT TB_Act.*,TB_AWA.[AwardName],TB_AWA.[Summary],TB_AWA.[IntroduceVideo],TB_AWA.[AwardImg],TB_AWA.[AwardType],TB_AWA.[AvailableDay],TB_AWA.[QQ],TB_AWA.[JRPTFunc],TB_AWA.[LinkUrl] FROM [UserHome_AwardInActivity] AS TB_Act WITH(NOLOCK) JOIN [UserHome_Award] AS TB_AWA WITH(NOLOCK) ON TB_Act.[AwardId]=TB_AWA.[AwardId] WHERE TB_Act.[IsDeleted]=0 AND TB_AWA.[IsDeleted]=0 AND TB_ACT.[ActivityId]=?`
	if ignoreExpired {
		sql += ` AND ((TB_Act.[BeginTime]<GETDATE() AND TB_Act.[EndTime]>GETDATE()) OR (TB_Act.[EndTime]='1900-01-01 00:00:00.000'))`
	}
	sql += ` ORDER BY TB_Act.[CreateTime] DESC`
	err = repo.FindList(&awardList, sql, activityId)
	return awardList, err
}
