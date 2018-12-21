package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type SeedUserInfoRepository struct {
	repository.BaseRepository
}

func NewSeedUserInfoRepository(conf *protected.ServiceConfig) *SeedUserInfoRepository {
	repo := &SeedUserInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetSeedUserInfoByCid 根据CID获取种子用户信息
func (repo *SeedUserInfoRepository) GetSeedUserInfoByCid(cid string) (seedUser *userhome_model.SeedUserInfo, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_SeedUserInfo] WITH(NOLOCK) WHERE [IsDeleted]=0 AND [Cid]=?"
	seedUser = new(userhome_model.SeedUserInfo)
	err = repo.FindOne(seedUser, sql, cid)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return seedUser, err
}

// GetSeedUserInfoList 获取所有种子用户列表，包括已经删除的
func (repo *SeedUserInfoRepository) GetSeedUserInfoList() (seedUserList []*userhome_model.SeedUserInfo, err error) {
	sql := "SELECT * FROM [UserHome_SeedUserInfo] WITH(NOLOCK)"
	err = repo.FindList(&seedUserList, sql)
	if err == sysSql.ErrNoRows {
		return seedUserList, nil
	}
	return seedUserList, err
}
