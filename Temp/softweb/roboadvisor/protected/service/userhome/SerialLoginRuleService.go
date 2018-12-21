package service

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userloginlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type SerialLoginRuleService struct {
	service.BaseService
	serialLoginRuleRepo *userhome_repo.SerialLoginRuleRepository
}

var (
	shareSerialLoginRuleRepo   *userhome_repo.SerialLoginRuleRepository
	shareSerialLoginRuleLogger dotlog.Logger
)

const (
	UserHomeSerialLoginRuleServiceName                = "SerialLoginRuleService"
	EMNET_SerialLoginRule_GetSerialLoginRule_CacheKey = UserHomeSerialLoginRuleServiceName + ":GetSerialLoginRule:"
	EMNET_SerialLoginRule_CacheSeconds                = 60 * 30
)

func NewSerialLoginRuleService() *SerialLoginRuleService {
	srv := &SerialLoginRuleService{
		serialLoginRuleRepo: shareSerialLoginRuleRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeSerialLoginRuleServiceName, userHomeSerialLoginRuleServiceLoader)
}

func userHomeSerialLoginRuleServiceLoader() {
	shareSerialLoginRuleRepo = userhome_repo.NewSerialLoginRuleRepository(protected.DefaultConfig)
	shareSerialLoginRuleLogger = dotlog.GetLogger(UserHomeSerialLoginRuleServiceName)
}

// GetUserLevelByLoginDay 根据用户累计登录天数计算相应的等级
func (srv *SerialLoginRuleService) GetUserLevelByLoginDay(loginDay int) int {
	userLevel := 1
	if loginDay <= 3 {
		userLevel = 1
	} else if loginDay <= 7 {
		userLevel = 2
	} else if loginDay <= 15 {
		userLevel = 3
	} else if loginDay <= 30 {
		userLevel = 4
	} else if loginDay <= 45 {
		userLevel = 5
	} else if loginDay <= 60 {
		userLevel = 6
	} else if loginDay <= 90 {
		userLevel = 7
	} else if loginDay <= 180 {
		userLevel = 8
	} else if loginDay <= 270 {
		userLevel = 9
	} else if loginDay <= 365 {
		userLevel = 10
	} else {
		userLevel = (loginDay / 365) + 10
	}
	return userLevel
}

// GetSerialLoginRule 根据登录天数和用户类型获取相应的奖励规则
func (srv *SerialLoginRuleService) GetSerialLoginRule(loginDay int, ztUserType int) (rule *userhome_model.SerialLoginRule, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.getSerialLoginRuleDB(loginDay, ztUserType)
	case config.ReadDB_CacheOrDB_UpdateCache:
		rule, err = srv.getSerialLoginRuleCache(loginDay, ztUserType)
		if err == nil && rule == nil {
			rule, err = srv.refreshSerialLoginRule(loginDay, ztUserType)
		}
		return rule, err
	case config.ReadDB_RefreshCache:
		rule, err = srv.refreshSerialLoginRule(loginDay, ztUserType)
		return rule, err
	default:
		return srv.getSerialLoginRuleCache(loginDay, ztUserType)
	}
}

// getSerialLoginRuleDB 根据登录天数和用户类型获取相应的奖励规则-读取数据库
func (srv *SerialLoginRuleService) getSerialLoginRuleDB(loginDay int, ztUserType int) (rule *userhome_model.SerialLoginRule, err error) {
	rule, err = srv.serialLoginRuleRepo.GetSerialLoginRule(loginDay, ztUserType)
	if err != nil {
		shareSerialLoginRuleLogger.ErrorFormat(err, "根据登录天数和用户类型获取相应的奖励规则异常 GetSerialLoginRule loginDay=%d ztUserType=%d", loginDay, ztUserType)
	}
	return rule, err
}

// getSerialLoginRuleCache 根据登录天数和用户类型获取相应的奖励规则-读取缓存
func (srv *SerialLoginRuleService) getSerialLoginRuleCache(loginDay int, ztUserType int) (rule *userhome_model.SerialLoginRule, err error) {
	cacheKey := EMNET_SerialLoginRule_GetSerialLoginRule_CacheKey + strconv.Itoa(loginDay) + ":" + strconv.Itoa(ztUserType)
	err = srv.RedisCache.GetJsonObj(cacheKey, &rule)
	if err == redis.ErrNil {
		return nil, nil
	}
	return rule, err
}

// refreshSerialLoginRule 根据登录天数和用户类型获取相应的奖励规则-刷新缓存
func (srv *SerialLoginRuleService) refreshSerialLoginRule(loginDay int, ztUserType int) (rule *userhome_model.SerialLoginRule, err error) {
	cacheKey := EMNET_SerialLoginRule_GetSerialLoginRule_CacheKey + strconv.Itoa(loginDay) + ":" + strconv.Itoa(ztUserType)
	rule, err = srv.getSerialLoginRuleDB(loginDay, ztUserType)
	if err != nil {
		return nil, err
	}
	if rule != nil {
		srv.RedisCache.Set(cacheKey, jsonutil.GetJsonString(rule), EMNET_SerialLoginRule_CacheSeconds)
	}
	return rule, err
}

// GetUserSerialLoginRule 获取用户连登奖励
func (srv *SerialLoginRuleService) GetUserSerialLoginRule(userInfo userhome_model.UserInfo) (rule *userhome_model.SerialLoginRule, err error) {
	shareSerialLoginRuleLogger.DebugFormat("GetUserSerialLoginRule 获取用户连登奖励 userInfo=%s", _json.GetJsonString(userInfo))
	loginLog, err := userloginlog.GetUserLoginCountSerial(userInfo.UID)
	if err != nil {
		shareSerialLoginRuleLogger.ErrorFormat(err, "GetUserSerialLoginRule->GetUserLoginCountSerial 获取用户连登奖励 异常 userInfo=%s loginLog=%s", _json.GetJsonString(userInfo), _json.GetJsonString(loginLog))
		return nil, err
	} else {
		shareSerialLoginRuleLogger.DebugFormat("GetUserSerialLoginRule->GetUserLoginCountSerial  获取用户连登奖励 userInfo=%s loginLog=%s", _json.GetJsonString(userInfo), _json.GetJsonString(loginLog))
	}
	rule, err = srv.GetSerialLoginRule(loginLog.LoginCountSerial, userInfo.PIDType)
	if err != nil {
		shareSerialLoginRuleLogger.ErrorFormat(err, "GetUserSerialLoginRule->GetSerialLoginRule 获取用户连登奖励 异常 userInfo=%s loginLog=%s rule=%s", _json.GetJsonString(userInfo), _json.GetJsonString(loginLog), _json.GetJsonString(rule))
	} else {
		shareSerialLoginRuleLogger.DebugFormat("GetUserSerialLoginRule->GetSerialLoginRule 获取用户连登奖励 userInfo=%s loginLog=%s  rule=%s", _json.GetJsonString(userInfo), _json.GetJsonString(loginLog), _json.GetJsonString(rule))
	}
	return rule, err
}
