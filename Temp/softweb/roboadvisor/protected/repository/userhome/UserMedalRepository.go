package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type UserMedalRepository struct {
	repository.BaseRepository
}

func NewUserMedalRepository(conf *protected.ServiceConfig) *UserMedalRepository {
	repo := &UserMedalRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertUserMedal插入一个新的勋章
func (repo *UserMedalRepository) InsertUserMedal(medal *userhome_model.UserMedal) (id int64, err error) {
	sql := "INSERT INTO [UserHome_UserMedal]([UserInfoId],[UID],[MedalType],[MedalName],[MedalLevel],[CreateTime],[LastModifyTime]) VALUES(?,?,?,?,?,GETDATE(),GETDATE())"
	id, err = repo.Insert(sql, medal.UserInfoId, medal.UID, medal.MedalType, medal.MedalName, medal.MedalLevel)
	return id, err
}

// GetUserMedal 查询某个用户的某个勋章
func (repo *UserMedalRepository) GetUserMedal(uid int64, medalType int) (userMedal *userhome_model.UserMedal, err error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserMedal] WHERE [UID]=? AND [MedalType]=? ORDER BY [MedalLevel] DESC"
	userMedal = new(userhome_model.UserMedal)
	err = repo.FindOne(userMedal, sql, uid, medalType)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return userMedal, err
}
