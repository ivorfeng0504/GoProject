package user

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"github.com/devfeel/dotlog"
)

type UserRepository struct {
	repository.BaseRepository
}

var (
	// shareUserLogger 共享的Logger实例
	shareUserLogger dotlog.Logger
)

func NewUserRepository(conf *protected.ServiceConfig) *UserRepository {
	repo := &UserRepository{}
	repo.Init(conf.FeedbDBConn)
	return repo
}

// QqAndWechat_Reg_For_ZT 第三方登录用户注册智投体验版
/*
ALTER  PROCEDURE [dbo].[QqAndWechat_Reg_For_ZT]
    @openid varchar(50),
    @regType INT , --  2 wechat  3 qq
    @sID INT ,
    @trackID INT ,
    @customerid bigint OUTPUT,
    @password varchar(20) OUTPUT,
    @msg VARCHAR(50) OUTPUT

返回结果集
*/
func (repo *UserRepository) QqAndWechat_Reg_For_ZT(openid string, regtype int, sid int, tid int) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("QqAndWechat_Reg_For_ZT", openid, regtype, sid, tid)

	return mapRet, err
}

// Mobile_Reg_For_ZT 手机注册获取密码接口
/*
alter proc [dbo].[Mobile_Reg_For_ZT]
	@mobile varchar(20),
	@sID int,
	@trackID int,
	@UserPasswd varchar(6) OUTPUT,
	@customerid bigint OUTPUT,
	@msg varchar(64) output

-1:手机号已注册
返回结果集
*/
func (repo *UserRepository) Mobile_Reg_For_ZT(mobile string, hardwareInfo string,clientVersion int, sid int, tid int) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("Mobile_Reg_For_ZT", mobile, hardwareInfo,clientVersion, sid, tid)

	return mapRet, err
}

// Stock_ChangePasswd_Reset_For_ZT 重置密码
/*
ALTER PROC [dbo].[Stock_ChangePasswd_Reset_For_ZT]
    @UserName VARCHAR(20) ,
    @Newpasswd VARCHAR(6)

返回2:成功
*/
func (repo *UserRepository) Stock_ChangePasswd_Reset_For_ZT(username string, newpwd string) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("Stock_ChangePasswd_Reset_result", username, newpwd)

	return mapRet, err
}

// 查询已绑定账号列表
/*
create proc BoundGroupQryLogin
	@curPID bigint		-- 当前帐号PID(UID), 返回帐号显示名, CustomerID, 帐号类型 0: em帐号; 1: 手机号; 2: 微信帐号; 3: QQ帐号
*/
func (repo *UserRepository) BoundGroupQryLogin(gid int64) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("BoundGroupQryLogin", gid)

	return mapRet, err
}

// BoundGroupAddLogin 添加绑定
/*
create proc BoundGroupAddLogin_ResultSet
@curUserName varchar(64),	-- 当前帐号名  em帐号; 手机号; wx_uniqueID; qq_openID
	@addUserName varchar(64),	-- 加入帐号名
	@addPassword varchar(20),	-- 加入帐号密码/昵称

select @addInCID as addInCID,@addInShowName as addInShowName,@addInType as addInType,@newCurDID as newCurDID,@addInOldPID as addInOldPID,@retMsg as retMsg
*/
func (repo *UserRepository) BoundGroupAddLogin(curUserName string, addUserName string, addPassword string) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("BoundGroupAddLogin_ResultSet", curUserName, addUserName, addPassword)

	return mapRet, err
}

// BoundGroupRmvLogin
/*
create proc BoundGroupRmvLogin_ResultSet
@curUserName varchar(64),	-- 当前帐号名  em帐号; 手机号; wx_uniqueID; qq_openID
@rmvCID bigint,				-- 移出帐号CustomerID

select @rmvNewPID as rmvNewPID,@retMsg as retMsg
*/
func (repo *UserRepository) BoundGroupRmvLogin(curUserName string, rmvCID int64) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("BoundGroupRmvLogin_ResultSet", curUserName, rmvCID)

	return mapRet, err
}

// GetLoginIDByName
/*
create proc GetLoginIDByName_ResultSet
@userName varchar(64),		-- 帐号名  em帐号; 手机号; wx_uniqueID; qq_openID
	@userPasswd varchar(32),	-- 密码/昵称(查询时可空)  em帐号/手机号: 对应密码; wx/qq: 昵称
	@createLogin int,			-- 查询/创建(登录或添加帐号, em号只在物流创建)  <>0: 创建; null/=0: 查询

select @userType as userType ,@showName as showName ,@guidCP as guidCP ,@uniqueID as uniqueID ,@CID as CID ,@DID as DID
*/
func (repo *UserRepository) GetLoginIDByName(userName string, userPasswd string, createLogin int) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("GetLoginIDByName_ResultSet", userName, userPasswd, createLogin)

	return mapRet, err
}

// EMGet_EndDate_New_Result
/*
create proc EMGet_EndDate_New_Result
@UserName varchar(20),
	@Password varchar(6),

SELECT convert(varchar(10),@getenddate,121) AS enddate ,convert(varchar(10),@getenddate_sz,121) AS enddatel2sz,convert(varchar(10),@getenddate_sh,121) AS enddatel2sh
*/
func (repo *UserRepository) EMGet_EndDate_New_Result(userName string) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("EMGet_EndDate_New_Result", userName)

	return mapRet, err
}

// LoginDaysAndProduct 获取用户连登天数以及产品信息
func (repo *UserRepository) LoginDaysAndProduct(cid int64) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("logindays_sel", cid)
	return mapRet, err
}

// CidFindGid 根据Cid获取Gid
func (repo *UserRepository) CidFindGid(cid int64) ([]map[string]interface{}, error) {
	var mapRet []map[string]interface{}
	mapRet, err := repo.ExecProc("Cid_find_Pid", cid)
	return mapRet, err
}
