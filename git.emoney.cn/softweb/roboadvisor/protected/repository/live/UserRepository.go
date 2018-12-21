package live

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type UserRepository struct {
	repository.BaseRepository
}

const (
	//智投账号注册
	RoboadvisorUserType = "6"
)

func NewUserRepository(conf *protected.ServiceConfig) *UserRepository {
	repo := &UserRepository{}
	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetUserById 通过用户Id获取用户信息
func (repo *UserRepository) GetUserById(userId int) (user *livemodel.User, err error) {
	sql := "SELECT TOP 1 * FROM [User] WITH(NOLOCK) WHERE [UserId]=? AND [UserType]=?"
	//必须初始化对象 否则会报错
	user = new(livemodel.User)
	err = repo.FindOne(user, sql, userId, RoboadvisorUserType)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByUID 通过用户UID获取用户信息
func (repo *UserRepository) GetUserByUID(uid int64) (user *livemodel.User, err error) {
	if uid <= 0 {
		return nil, nil
	}
	sql := "SELECT TOP 1 * FROM [User] WITH(NOLOCK) WHERE [UID]=? AND [UserType]=?"
	//必须初始化对象 否则会报错
	user = new(livemodel.User)
	err = repo.FindOne(user, sql, uid, RoboadvisorUserType)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByAccount 通过账号获取用户信息
func (repo *UserRepository) GetUserByAccount(account string) (user *livemodel.User, err error) {
	if len(account) == 0 {
		return nil, nil
	}
	sql := "SELECT TOP 1 * FROM [User] WITH(NOLOCK) WHERE [Account]=? AND [UserType]=?"
	//必须初始化对象 否则会报错
	user = new(livemodel.User)
	err = repo.FindOne(user, sql, account, RoboadvisorUserType)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// AddUser 通过账号和UID自动注册用户
func (repo *UserRepository) AddUser(account string, nickName string, uid int64, source string) (userId int, err error) {
	//检查是否已经用account注册过
	user, err := repo.GetUserByAccount(account)
	if err != nil {
		userId = 0
		return userId, err
	}
	if user != nil {
		//如果数据库中的UID不存在则尝试更新
		if user.UID <= 0 && uid > 0 {
			repo.UpdateUIDByUserId(user.UserId, uid)
		}
		return user.UserId, err
	}
	//检查是否已经用uid注册过
	user, err = repo.GetUserByUID(uid)
	if err != nil {
		userId = 0
		return userId, err
	}
	if user != nil {
		//如果数据库中的Account不存在则尝试更新
		if len(user.Account) == 0 && len(account) > 0 {
			repo.UpdateAccountByUserId(user.UserId, account)
		}
		return user.UserId, err
	}
	sql := "INSERT INTO [User] ([Account],[UserName],[CreateTime],[Status],[LastLoginTime],[NickName],[UserType],[UID],[Source]) VALUES(?,?,?,?,?,?,?,?,?)"
	//正常
	userStatus := 0
	userId64, err := repo.Insert(sql, account, nickName, time.Now(), userStatus, time.Now(), nickName, RoboadvisorUserType, uid, source)
	userId = int(userId64)
	return userId, err
}

// UpdateAccountByUserId 更新Account
func (repo *UserRepository) UpdateAccountByUserId(userId int, account string) (err error) {
	sql := "UPDATE [User] SET [Account]=? WHERE [UserId]=? AND [UserType]=? AND ([Account] IS NULL OR [Account]='')"
	_, err = repo.Update(sql, account, userId, RoboadvisorUserType)
	return err
}

// UpdateUIDByUserId 更新UID
func (repo *UserRepository) UpdateUIDByUserId(userId int, uid int64) (err error) {
	sql := "UPDATE [User] SET [UID]=? WHERE [UserId]=? AND [UserType]=? AND [UID]<=0"
	_, err = repo.Update(sql, uid, userId, RoboadvisorUserType)
	return err
}
