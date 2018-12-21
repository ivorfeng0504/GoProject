package user

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/config"
)

type SoftDataSyncRepository struct {
	repository.MySqlBaseRepository
}

var(
	// shareSoftDataSyncLogger 共享的Logger实例
	shareSoftDataSyncLogger dotlog.Logger
)


func NewSoftDataSyncRepository(conf *protected.ServiceConfig) *SoftDataSyncRepository {
	repo := &SoftDataSyncRepository{}
	repo.Init(config.CurrentConfig.SoftDataSyncMysqlURL)
	return repo
}

// SoftUserDataSync 同步
/*
create proc QKTaste_RegistMobile_Web
@mobile varchar(11),
@hardwareInfo varchar(64),
@sID int,
@trackID int
select retCode,retMsg,passWord
*/
func (repo *SoftDataSyncRepository) SoftUserDataSync(uid int64,newpid int64) ([]map[string]interface{},error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("CopyUserData", uid,newpid)

	return mapRet, err
}