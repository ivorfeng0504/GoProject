package repository

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type UserBehaviorStatRepository struct {
	repository.BaseRepository
}

func NewUserBehaviorStatRepository(conf *protected.ServiceConfig) *UserBehaviorStatRepository {
	repo := &UserBehaviorStatRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}
