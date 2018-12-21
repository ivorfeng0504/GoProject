package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type SerialLoginRuleRepository struct {
	repository.BaseRepository
}

func NewSerialLoginRuleRepository(conf *protected.ServiceConfig) *SerialLoginRuleRepository {
	repo := &SerialLoginRuleRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetSerialLoginRule 根据登录天数和用户类型获取相应的奖励规则 TargetUserType为0 全部用户
func (repo *SerialLoginRuleRepository) GetSerialLoginRule(loginDay int, ztUserType int) (rule *userhome_model.SerialLoginRule, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_SerialLoginRule] WHERE [IsDeleted]=0 AND [LoginDay]=? AND ([TargetUserType]=? OR [TargetUserType]=0)"
	rule = new(userhome_model.SerialLoginRule)
	err = repo.FindOne(rule, sql, loginDay, ztUserType)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return rule, err
}
