package live

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	liverepo "git.emoney.cn/softweb/roboadvisor/protected/repository/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type UserService struct {
	service.BaseService
	userRepo *liverepo.UserRepository
}

var (
	shareUserRepo   *liverepo.UserRepository
	shareUserLogger dotlog.Logger
)

const (
	userServiceName          = "UserService"
	GetUserByIdCacheKey      = "live.UserService.GetUserById:"
	GetUserByAccountCacheKey = "live.UserService.GetUserByAccount:"
	GetUserByUIDCacheKey     = "live.UserService.GetUserByUID:"
	AddUserLockCacheKey      = "live.UserService.Locker:"
)

func NewUserService() *UserService {
	userService := &UserService{
		userRepo: shareUserRepo,
	}
	userService.RedisCache = _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	return userService
}

// GetUserById 通过用户Id获取用户信息
func (service *UserService) GetUserById(userId int) (user *livemodel.User, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getUserByIdDB(userId)
	case config.ReadDB_CacheOrDB_UpdateCache:
		user, err = service.getUserByIdCache(userId)
		//缓存不存在则读取数据库
		if err == nil && user == nil {
			user = service.refreshUserByUserId(userId)
		}
		return user, err
	default:
		return service.getUserByIdCache(userId)
	}
}

// getUserByIdDB 通过用户Id获取用户信息
func (service *UserService) getUserByIdDB(userId int) (user *livemodel.User, err error) {
	user, err = service.userRepo.GetUserById(userId)
	if err != nil {
		shareUserLogger.Error(err, "通过用户Id获取用户信息异常 getUserByIdDB userId="+strconv.Itoa(userId))
	}
	return user, err
}

// getUserByIdCache 通过用户Id获取用户信息
func (service *UserService) getUserByIdCache(userId int) (user *livemodel.User, err error) {
	cacheKey := GetUserByIdCacheKey + strconv.Itoa(userId)
	err = service.RedisCache.GetJsonObj(cacheKey, &user)
	if err == redis.ErrNil {
		return user, nil
	}
	if err != nil {
		shareUserLogger.Error(err, "读取缓存-通过用户Id获取用户信息 getUserByIdCache 异常")
	}
	return user, err
}

// GetUserByAccount 通过账号获取用户信息
func (service *UserService) GetUserByAccount(account string) (user *livemodel.User, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getUserByAccountDB(account)
	case config.ReadDB_CacheOrDB_UpdateCache:
		user, err = service.getUserByAccountCache(account)
		//缓存不存在则读取数据库
		if err == nil && user == nil {
			user = service.refreshUserByAccount(account)
		}
		return user, err
	default:
		return service.getUserByAccountCache(account)
	}
}

// getUserByAccountDB 通过账号获取用户信息
func (service *UserService) getUserByAccountDB(account string) (user *livemodel.User, err error) {
	user, err = service.userRepo.GetUserByAccount(account)
	if err != nil {
		shareUserLogger.Error(err, "通过账号获取用户信息 getUserByAccountDB account="+account)
	}
	return user, err
}

// getUserByAccountCache 通过账号获取用户信息
func (service *UserService) getUserByAccountCache(account string) (user *livemodel.User, err error) {
	cacheKey := GetUserByAccountCacheKey + account
	err = service.RedisCache.GetJsonObj(cacheKey, &user)
	if err == redis.ErrNil {
		return user, nil
	}
	if err != nil {
		shareUserLogger.Error(err, "读取缓存-通过账号获取用户信息 getUserByAccountCache 异常")
	}
	return user, err
}

// GetUserByUID 通过UID账号获取用户信息
func (service *UserService) GetUserByUID(uid int64) (user *livemodel.User, err error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.getUserByUIDDB(uid)
	case config.ReadDB_CacheOrDB_UpdateCache:
		user, err = service.getUserByUIDCache(uid)
		//缓存不存在则读取数据库
		if err == nil && user == nil {
			user = service.refreshUserByUID(uid)
		}
		return user, err
	default:
		return service.getUserByUIDCache(uid)
	}
}

// getUserByUIDDB 通过UID账号获取用户信息
func (service *UserService) getUserByUIDDB(uid int64) (user *livemodel.User, err error) {
	user, err = service.userRepo.GetUserByUID(uid)
	if err != nil {
		shareUserLogger.Error(err, "通过账号获取用户信息异常 getUserByUIDDB uid="+strconv.FormatInt(uid, 10))
	}
	return user, err
}

// getUserByUIDCache 通过UID账号获取用户信息
func (service *UserService) getUserByUIDCache(uid int64) (user *livemodel.User, err error) {
	cacheKey := GetUserByUIDCacheKey + strconv.FormatInt(uid, 10)
	err = service.RedisCache.GetJsonObj(cacheKey, &user)
	if err == redis.ErrNil {
		return user, nil
	}
	if err != nil {
		shareUserLogger.Error(err, "读取缓存-通过UID账号获取用户信息 getUserByUIDCache 异常")
	}
	return user, err
}

// AddUser 通过账号和UID自动注册用户
func (service *UserService) AddUser(account string, nickName string, uid int64, source string) (userId int, err error) {
	lockKey := AddUserLockCacheKey + account
	counter, err := service.RedisCache.Incr(lockKey)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "1.获取用户锁失败key=%s", lockKey)
		return 0, err
	}
	shareUserLogger.DebugFormat("1.尝试获取用户锁key=%s counter=%s", lockKey, counter)
	//每100毫秒循环尝试获取锁
	maxTry := 20
	for counter != 1 && maxTry > 0 {
		time.Sleep(time.Millisecond * 100)
		counter, err = service.RedisCache.Incr(lockKey)
		if err != nil {
			shareUserLogger.ErrorFormat(err, "2.获取用户锁失败key=%s 计数=%s", lockKey, maxTry)
			return 0, err
		}
		shareUserLogger.DebugFormat("2.尝试获取用户锁key=%s 计数=%s counter=%s", lockKey, maxTry, counter)
		maxTry = maxTry - 1
	}

	//多次获取锁失败后自动放行 防止因为程序中断导致永远都不释放锁
	if counter != 1 {
		err = errors.New("服务器繁忙，请稍后重试！")
		shareUserLogger.ErrorFormat(err, "服务器繁忙 多次获取用户锁失败后自动放行  key=%s 计数=%s", lockKey, maxTry)
		//return 0, err
	}
	userId, err = service.userRepo.AddUser(account, nickName, uid, source)
	if err != nil {
		shareUserLogger.ErrorFormat(err, "通过账号自动注册用户异常 AddUser account=%s  uid=%s  nickname=%s source=%s", account, strconv.FormatInt(uid, 10), nickName, source)
	} else {
		//刷新缓存
		service.refreshUserByUserId(userId)
	}
	//删除键 释放锁
	service.RedisCache.Delete(lockKey)
	return userId, err
}

// refreshUserByUserId 刷新用户缓存
func (service *UserService) refreshUserByUserId(userId int) (userInfo *livemodel.User) {
	userInfo, err := service.userRepo.GetUserById(userId)
	if err != nil || userInfo == nil {
		return nil
	}
	cacheKey := GetUserByIdCacheKey + strconv.Itoa(userId)
	service.RedisCache.SetJsonObj(cacheKey, userInfo)
	if userInfo.Account != "" {
		cacheKey = GetUserByAccountCacheKey + userInfo.Account
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	if userInfo.UID > 0 {
		cacheKey = GetUserByUIDCacheKey + strconv.FormatInt(userInfo.UID, 10)
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	return userInfo
}

// refreshUserByAccount 刷新用户缓存
func (service *UserService) refreshUserByAccount(account string) (userInfo *livemodel.User) {
	userInfo, err := service.userRepo.GetUserByAccount(account)
	if err != nil || userInfo == nil {
		return nil
	}
	cacheKey := GetUserByIdCacheKey + strconv.Itoa(userInfo.UserId)
	service.RedisCache.SetJsonObj(cacheKey, userInfo)
	if userInfo.Account != "" {
		cacheKey = GetUserByAccountCacheKey + userInfo.Account
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	if userInfo.UID > 0 {
		cacheKey = GetUserByUIDCacheKey + strconv.FormatInt(userInfo.UID, 10)
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	return userInfo
}

// refreshUserByUID 刷新用户缓存
func (service *UserService) refreshUserByUID(uid int64) (userInfo *livemodel.User) {
	userInfo, err := service.userRepo.GetUserByUID(uid)
	if err != nil || userInfo == nil {
		return nil
	}
	cacheKey := GetUserByIdCacheKey + strconv.Itoa(userInfo.UserId)
	service.RedisCache.SetJsonObj(cacheKey, userInfo)
	if userInfo.Account != "" {
		cacheKey = GetUserByAccountCacheKey + userInfo.Account
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	if userInfo.UID > 0 {
		cacheKey = GetUserByUIDCacheKey + strconv.FormatInt(userInfo.UID, 10)
		service.RedisCache.SetJsonObj(cacheKey, userInfo)
	}
	return userInfo
}

func init() {
	protected.RegisterServiceLoader(userServiceName, userServiceLoader)
}

func userServiceLoader() {
	shareUserRepo = liverepo.NewUserRepository(protected.DefaultConfig)
	shareUserLogger = dotlog.GetLogger(userServiceName)
}
