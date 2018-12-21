package user

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"fmt"
)

type TyUserRepository struct {
	repository.BaseRepository
}

var(
// shareTyUserLogger 共享的Logger实例
shareTyUserLogger dotlog.Logger
)

func NewTyUserRepository(conf *protected.ServiceConfig) *TyUserRepository {
	repo := &TyUserRepository{}
	repo.Init(conf.EmoneyDBConn)
	return repo
}

// QKTaste_RegistMobile_Web 获取手机密码，未注册注册后返回，已注册直接返回密码
/*
create proc QKTaste_RegistMobile_Web
@mobile varchar(11),
@hardwareInfo varchar(64),
@sID int,
@trackID int
select retCode,retMsg,passWord
*/
func (repo *TyUserRepository) QKTaste_RegistMobile_Web(mobile string,hardinfo string,sid int ,tid int) ([]map[string]interface{},error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("QKTaste_RegistMobile_Web_ResultSet", mobile,"", sid, tid)

	fmt.Println(mapRet)
	return mapRet, err
}