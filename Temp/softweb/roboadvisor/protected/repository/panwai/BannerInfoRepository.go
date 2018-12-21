package repository

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/panwai"
)

type BannerInfoRepository struct {
	repository.BaseRepository
}

func NewBannerInfoRepository(conf *protected.ServiceConfig) *BannerInfoRepository {
	repo := &BannerInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetNewsInfoByID 根据newsID获取资讯
func (repo *BannerInfoRepository) GetBannerInfoList() (bannerinfo []*model.BannerInfo,err error){
	sql := "SELECT * FROM [BannerInfo] WITH(NOLOCK) WHERE IsDeleted=0"
	err = repo.FindList(&bannerinfo, sql)
	return bannerinfo, err
}