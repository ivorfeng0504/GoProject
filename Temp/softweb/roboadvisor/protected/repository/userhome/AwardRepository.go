package repository

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type AwardRepository struct {
	repository.BaseRepository
}

func NewAwardRepository(conf *protected.ServiceConfig) *AwardRepository {
	repo := &AwardRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}
