package repository

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type UserInfoRepository struct {
	repository.BaseRepository
}

func NewUserInfoRepository(conf *protected.ServiceConfig) *UserInfoRepository {
	repo := &UserInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetUserInfoByAccount 根据Account获取用户信息
func (repo *UserInfoRepository) GetUserInfoByAccount(account string) (*model.UserInfo, error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInfo] WITH(NOLOCK) WHERE Account = ? AND IsDeleted=0 order by CreateTime asc"

	//必须初始化对象 否则会报错
	newsinfo := new(model.UserInfo)
	err := repo.FindOne(newsinfo, sql, account)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}



// GetUserInfoByUID 根据UID获取用户信息
func (repo *UserInfoRepository) GetUserInfoByUID(UID int64) (*model.UserInfo, error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInfo] WITH(NOLOCK) WHERE UID = ? AND IsDeleted=0 order by CreateTime asc"

	//必须初始化对象 否则会报错
	newsinfo := new(model.UserInfo)
	err := repo.FindOne(newsinfo, sql, UID)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}

// GetUserInfoByUIDAndAccount 根据UID和非手机号账号获取用户信息
func (repo *UserInfoRepository) GetUserInfoByUIDAndAccount(UID int64, account string) (*model.UserInfo, error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInfo] WITH(NOLOCK) WHERE UID = ? AND Account = ? AND IsDeleted=0"

	//必须初始化对象 否则会报错
	newsinfo := new(model.UserInfo)
	err := repo.FindOne(newsinfo, sql, UID,account)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}

// GetUserInfoByUIDAndMobile 根据UID和手机账号获取用户信息
func (repo *UserInfoRepository) GetUserInfoByUIDAndMobile(UID int64, mobilex string) (*model.UserInfo, error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInfo] WITH(NOLOCK) WHERE UID = ? AND MobileX = ? AND IsDeleted=0"

	//必须初始化对象 否则会报错
	newsinfo := new(model.UserInfo)
	err := repo.FindOne(newsinfo, sql, UID,mobilex)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}

// 是否存在昵称
func (repo *UserInfoRepository) HasNickName(nickname string) (*model.UserInfo, error) {
	sql := "SELECT TOP 1 * FROM [UserHome_UserInfo] WITH(NOLOCK) WHERE NickName = ? AND IsDeleted=0"

	//必须初始化对象 否则会报错
	newsinfo := new(model.UserInfo)
	err := repo.FindOne(newsinfo, sql, nickname)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return newsinfo, err
}

// AddUserInfo 用户信息注册
func (repo *UserInfoRepository) AddUserInfo(info *model.UserInfo) (id int64, err error) {
	defalutUserLevel := 1
	sql := "INSERT INTO [UserHome_UserInfo] ([UID],[UserType],[Account],[MobileX],[MobileMask],[NickName],[Headportrait],[OpenID_WeChat],[OpenID_QQ],[UserLevel]) VALUES(?,?,?,?,?,?,?,?,?,?) "
	id64, err := repo.Insert(sql, info.UID, info.UserType, info.Account, info.MobileX, info.MobileMask, info.NickName, info.Headportrait, info.OpenID_WeChat, info.OpenID_QQ, defalutUserLevel)
	id = id64
	return id, err
}

// ModifyLastLoginTime 更新最后登录时间（根据UID）
func (repo *UserInfoRepository) ModifyLastLoginTime(UID int64) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET LastLoginTime=getdate() WHERE UID=?"
	id64, err := repo.Update(sql, UID)
	id = int(id64)
	return id, err
}

// ModifyLastLoginTimeByAccount 更新最后登录时间
func (repo *UserInfoRepository) ModifyLastLoginTimeByAccount(Account string) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET LastLoginTime=getdate() WHERE Account=?"
	id64, err := repo.Update(sql, Account)
	id = int(id64)
	return id, err
}

// ModifyNickName 修改昵称(根据UID）
func (repo *UserInfoRepository) ModifyNickName(UID int64, NickName string) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET NickName=? WHERE UID=?"
	id64, err := repo.Update(sql, NickName, UID)
	id = int(id64)
	return id, err
}

// ModifyNickNameByAccount 修改昵称(根据Account）
func (repo *UserInfoRepository) ModifyNickNameByAccount(Account string, NickName string) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET NickName=? WHERE Account=?"
	id64, err := repo.Update(sql, NickName, Account)
	id = int(id64)
	return id, err
}

// ModifyHeadportrait 修改头像（根据UID）
func (repo *UserInfoRepository) ModifyHeadportrait(UID int64, Headportrait string) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET Headportrait=? WHERE UID=?"
	id64, err := repo.Update(sql, Headportrait, UID)
	id = int(id64)
	return id, err
}

// ModifyHeadportraitByAccount 修改头像（根据Account）
func (repo *UserInfoRepository) ModifyHeadportraitByAccount(Account string, Headportrait string) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET Headportrait=? WHERE Account=?"
	id64, err := repo.Update(sql, Headportrait, Account)
	id = int(id64)
	return id, err
}

// ModifyUserLevel 修改用户等级（根据UID）
func (repo *UserInfoRepository) ModifyUserLevel(UID int64, userLevel int) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET [UserLevel]=? WHERE UID=?"
	id64, err := repo.Update(sql, userLevel, UID)
	id = int(id64)
	return id, err
}

// ModifyUserLevel 修改用户等级（根据Account）
func (repo *UserInfoRepository) ModifyUserLevelByAccount(Account string, userLevel int) (id int, err error) {
	sql := "UPDATE [UserHome_UserInfo] SET [UserLevel]=? WHERE Account=?"
	id64, err := repo.Update(sql, userLevel, Account)
	id = int(id64)
	return id, err
}

// ModifyMobile 根据UID修改用户手机（绑定手机号成功后调用）
func (repo *UserInfoRepository) ModifyMobile(UID int64,mobilemask string ,mobilex string)(id int,err error) {
	sql := "UPDATE [UserHome_UserInfo] SET [MobileMask]=?,[MobileX]=? WHERE UID=?"
	id64, err := repo.Update(sql, mobilemask, mobilex, UID)
	id = int(id64)
	return id, err
}
