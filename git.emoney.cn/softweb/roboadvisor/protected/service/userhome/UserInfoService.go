package service

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"strconv"
	"encoding/json"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/protected/service/live"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type UserInfoService struct {
	service.BaseService
	userRepo *repository.UserInfoRepository
}

var (
	// shareUserRepo UserInfoRepository
	shareUserRepo *repository.UserInfoRepository

	// shareUserLogger 共享的Logger实例
	shareUserLogger dotlog.Logger
)

const (
	userInfoServiceName = "UserInfoServiceLogger"
	EMNET_UserInfo_CachePreKey  = "EMoney:RoboAdvisor:UserHome_UserInfo"
	EMNET_UserInfo_GetUserInfoByUIDCacheKey = EMNET_UserInfo_CachePreKey+":GetUserInfoByUID"
	EMNET_UserInfo_GetUserInfoByAccountCacheKey = EMNET_UserInfo_CachePreKey+":GetUserInfoByAccount"
)

func NewUserInfoService() *UserInfoService {
	userService := &UserInfoService{
		userRepo: shareUserRepo,
	}
	userService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return userService
}

// GetUserInfoByAccount 根据Account获取用户信息from redis
func (service *UserInfoService) GetUserInfoByAccount(Account string) (*model.UserInfo,error) {
	cacheKey := EMNET_UserInfo_GetUserInfoByAccountCacheKey + Account
	result := new(model.UserInfo)
	//get from redis
	_ = service.RedisCache.GetJsonObj(cacheKey, result)

	if result == nil {
		info, err := service.userRepo.GetUserInfoByAccount(Account)
		if err != nil {
			shareUserLogger.Error(err, "GetUserInfoByAccount 获取用户信息异常 Account="+Account)
			return info, err
		}
		return info, err
	}
	return result, nil
}

// GetUserInfoByUID 根据UID获取用户信息from redis
func (service *UserInfoService) GetUserInfoByUID(UID int64) (*model.UserInfo,error) {
	cacheKey := EMNET_UserInfo_GetUserInfoByUIDCacheKey + strconv.FormatInt(UID, 10)
	result := new(model.UserInfo)
	//get from redis
	_ = service.RedisCache.GetJsonObj(cacheKey, result)

	if result == nil {
		info, err := service.userRepo.GetUserInfoByUID(UID)
		if err != nil {
			shareUserLogger.Error(err, "GetUserInfoByUID 获取用户信息异常 UID="+strconv.FormatInt(UID, 10))
			return info, err
		}
		return info, err
	}
	return result, nil
}

// GetUserInfoByUIDAndAccount 根据UID获取用户信息from redis
func (service *UserInfoService) GetUserInfoByUIDAndAccount(UID int64,Account string) (*model.UserInfo,error) {
	info, err := service.userRepo.GetUserInfoByUIDAndAccount(UID,Account)
	if err != nil {
		shareUserLogger.Error(err, "GetUserInfoByUIDAndAccount 获取用户信息异常 UID="+strconv.FormatInt(UID, 10))
		return info, err
	}
	return info,err
}

// GetUserInfoByUIDAndMobile 根据UID获取用户信息from redis
func (service *UserInfoService) GetUserInfoByUIDAndMobile(UID int64,MobileX string) (*model.UserInfo,error) {
	info, err := service.userRepo.GetUserInfoByUIDAndMobile(UID,MobileX)
	if err != nil {
		shareUserLogger.Error(err, "GetUserInfoByUIDAndMobile 获取用户信息异常 UID="+strconv.FormatInt(UID, 10))
		return info, err
	}
	return info,err
}

// 昵称是否存在数据库
func (service *UserInfoService) IsExistNickname(nickname string) (bool,error) {
	info, err := service.userRepo.HasNickName(nickname)
	if err != nil {
		shareUserLogger.Error(err, "IsExistNickname 查询是否有改昵称异常 nickname="+nickname)
		return false, err
	}

	if info!=nil && info.NickName!=""{
		return true,nil
	}

	return false,nil
}

// AddUserInfo 用户信息录入
func (service *UserInfoService) AddUserInfo(info *model.UserInfo) (int64,error) {
	Userjsonstr, err := json.Marshal(info)
	id, err := service.userRepo.AddUserInfo(info)
	if err != nil {
		shareUserLogger.Error(err, "AddUserInfo 用户信息注册 异常user="+string(Userjsonstr))
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByUID(info.UID)
	return id, err
}

// ModifyLastLoginTime 更新用户最后登录时间
func (service *UserInfoService) ModifyLastLoginTime(UID int64) (int,error) {
	id, err := service.userRepo.ModifyLastLoginTime(UID)
	if err != nil {
		shareUserLogger.Error(err, "ModifyLastLoginTime 更新用户最后登录时间 异常UID="+strconv.FormatInt(UID, 10))
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByUID(UID)
	return id, err
}

// ModifyLastLoginTimeByAccount 更新用户最后登录时间
func (service *UserInfoService) ModifyLastLoginTimeByAccount(Account string) (int,error) {
	id, err := service.userRepo.ModifyLastLoginTimeByAccount(Account)
	if err != nil {
		shareUserLogger.Error(err, "ModifyLastLoginTimeByAccount 更新用户最后登录时间 异常Account="+Account)
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByAccount(Account)
	return id, err
}

// ModifyNickName 更新用户昵称
func (service *UserInfoService) ModifyNickName(UID int64,NickName string) (int,string,error) {
	//屏蔽关键字
	maskService := live.NewMaskWordService()
	NickName, err := maskService.ProcessMaskWord(NickName)

	id, err := service.userRepo.ModifyNickName(UID,NickName)
	if err != nil {
		shareUserLogger.Error(err, "ModifyNickName 更新用户昵称 异常UID="+strconv.FormatInt(UID, 10))
		return id,NickName, err
	}

	//更新缓存
	service.freshCacheUserInfoByUID(UID)
	return id,NickName, err
}

// ModifyNickNameByAccount 更新用户昵称
func (service *UserInfoService) ModifyNickNameByAccount(Account string,NickName string) (int,string,error) {
	//屏蔽关键字
	maskService := live.NewMaskWordService()
	NickName, err := maskService.ProcessMaskWord(NickName)

	id, err := service.userRepo.ModifyNickNameByAccount(Account, NickName)
	if err != nil {
		shareUserLogger.Error(err, "ModifyNickNameByAccount 更新用户昵称 异常Account="+Account)
		return id, NickName, err
	}

	//更新缓存
	service.freshCacheUserInfoByAccount(Account)
	return id, NickName, err
}

// ModifyHeadportrait 更新用户头像
func (service *UserInfoService) ModifyHeadportrait(UID int64,Headportrait string) (int,error) {
	id, err := service.userRepo.ModifyHeadportrait(UID,Headportrait)
	if err != nil {
		shareUserLogger.Error(err, "ModifyHeadportrait 更新用户昵称 异常UID="+strconv.FormatInt(UID, 10))
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByUID(UID)
	return id, err
}

// ModifyHeadportraitByAccount 更新用户头像
func (service *UserInfoService) ModifyHeadportraitByAccount(Account string,Headportrait string) (int,error) {
	id, err := service.userRepo.ModifyHeadportraitByAccount(Account,Headportrait)
	if err != nil {
		shareUserLogger.Error(err, "ModifyHeadportraitByAccount 更新用户昵称 异常Account="+Account)
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByAccount(Account)
	return id, err
}

// ModifyMobile 更新用户头像（绑定手机成功调用）
func (service *UserInfoService) ModifyMobile(UID int64,mobilemask string ,mobilex string) (int,error) {
	id, err := service.userRepo.ModifyMobile(UID, mobilemask, mobilex)
	if err != nil {
		shareUserLogger.Error(err, "ModifyMobile 更新用户手机 异常UID="+strconv.FormatInt(UID, 10))
		return id, err
	}

	//更新缓存
	service.freshCacheUserInfoByUID(UID)
	return id, err
}

// freshCacheUserInfoByUID 持久化用户信息-缓存刷新
func (service *UserInfoService) freshCacheUserInfoByUID(UID int64) (info *model.UserInfo, err error) {
	info, err = service.userRepo.GetUserInfoByUID(UID)
	if err != nil {
		shareUserLogger.Error(err, "刷新缓存异常 freshCacheUserInfoByUID UID="+strconv.FormatInt(UID, 10))
		return info, err
	}
	cacheKey := EMNET_UserInfo_GetUserInfoByUIDCacheKey + strconv.FormatInt(UID, 10)
	fmt.Println(info)
	_, err = service.RedisCache.SetJsonObj(cacheKey, info)
	return info, err
}

// freshCacheUserInfoByAccount 持久化用户信息-缓存刷新
func (service *UserInfoService) freshCacheUserInfoByAccount(account string) (info *model.UserInfo, err error) {
	info, err = service.userRepo.GetUserInfoByAccount(account)
	if err != nil {
		shareUserLogger.Error(err, "刷新缓存异常 freshCacheUserInfoByAccount account="+account)
		return info, err
	}
	cacheKey := EMNET_UserInfo_GetUserInfoByAccountCacheKey + account
	fmt.Println(info)
	_, err = service.RedisCache.SetJsonObj(cacheKey, info)
	return info, err
}


func init() {
	protected.RegisterServiceLoader(userInfoServiceName, userServiceLoader)
}

func userServiceLoader() {
	shareUserRepo = repository.NewUserInfoRepository(protected.DefaultConfig)
	shareUserLogger = dotlog.GetLogger(userInfoServiceName)
}

