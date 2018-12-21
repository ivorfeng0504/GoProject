package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
)

type FocusStrategyRepository struct {
	repository.BaseRepository
}

func NewFocusStrategyRepository(conf *protected.ServiceConfig) *FocusStrategyRepository {
	repo := &FocusStrategyRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}



// AddFocusStrategy 用户关注策略记录入库
func (repo *FocusStrategyRepository) AddFocusStrategy(uid int64,strategyId int) (err error){
	sql := "INSERT [ExpertNews_FocusStrategyInfo]([UID],[ExpertStrategyID]) VALUES(?,?)"
	_, err = repo.Update(sql, uid, strategyId)
	return err
}

// RemoveFocusStrategy 取消已关注策略
func (repo *FocusStrategyRepository) RemoveFocusStrategy(uid int64,strategyId int) (err error){
	sql := "UPDATE [ExpertNews_FocusStrategyInfo] SET [IsDeleted]=1 WHERE [UID]=? AND [ExpertStrategyID]=?"
	_, err = repo.Update(sql, uid, strategyId)
	return err
}

// GetFocusStrategyListByUID 获取已关注策略信息列表根据UID
func (repo *FocusStrategyRepository) GetFocusStrategyListByUID(uid int64) (strategyList []*expertnews.ExpertNews_StrategyInfo, err error){
	sql:="SELECT b.* FROM ExpertNews_FocusStrategyInfo a LEFT JOIN ExpertNews_StrategyInfo b ON a.ExpertStrategyID=b.ID WHERE a.UID=?"
	err = repo.FindList(&strategyList, sql)
	return strategyList, err
}


