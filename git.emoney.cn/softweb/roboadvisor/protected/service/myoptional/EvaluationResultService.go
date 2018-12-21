package myoptional

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	myoptional_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type EvaluationResultService struct {
	service.BaseService
	evaluationResultRepo *myoptional_repo.EvaluationResultRepository
}

var (
	shareEvaluationResultRepo   *myoptional_repo.EvaluationResultRepository
	shareEvaluationResultLogger dotlog.Logger
)

const (
	evaluationResultServiceName                         = "EvaluationResultService"
	EMNET_EvaluationResult_PreCacheKey                  = "EvaluationResultService:"
	EMNET_EvaluationResult_GetEvaluationResult_CacheKey = EMNET_EvaluationResult_PreCacheKey + "GetEvaluationResult:"
	EMNET_EvaluationResult_CacheSeconds                 = 60 * 30
)

func NewEvaluationResultService() *EvaluationResultService {
	service := &EvaluationResultService{
		evaluationResultRepo: shareEvaluationResultRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(evaluationResultServiceName, evaluationResultServiceLoader)
}

func evaluationResultServiceLoader() {
	shareEvaluationResultRepo = myoptional_repo.NewEvaluationResultRepository(protected.DefaultConfig)
	shareEvaluationResultLogger = dotlog.GetLogger(evaluationResultServiceName)
}

// InsertEvaluationResult 新增一个评测记录
func (srv *EvaluationResultService) InsertEvaluationResult(result myoptional_model.EvaluationResult) (err error) {
	if len(result.UID) == 0 {
		return errors.New("UID不正确")
	}
	if len(result.Result) == 0 || len(result.BuyStyle) == 0 || len(result.HoldStockTime) == 0 || len(result.ChooseStockReason) == 0 || len(result.InvestTarget) == 0 {
		return errors.New("评测内容不完整")
	}
	err = srv.evaluationResultRepo.InsertEvaluationResult(result)
	if err != nil {
		shareEvaluationResultLogger.ErrorFormat(err, "InsertEvaluationResult 新增一个评测记录 异常 result=%s", _json.GetJsonString(result))
	} else {
		srv.refreshEvaluationResult(result.UID)
	}
	return err
}

// UpdateEvaluationResult 更新评测记录
func (srv *EvaluationResultService) UpdateEvaluationResult(result myoptional_model.EvaluationResult) (err error) {
	if len(result.UID) == 0 {
		return errors.New("UID不正确")
	}
	if len(result.Result) == 0 || len(result.BuyStyle) == 0 || len(result.HoldStockTime) == 0 || len(result.ChooseStockReason) == 0 || len(result.InvestTarget) == 0 {
		return errors.New("评测内容不完整")
	}
	err = srv.evaluationResultRepo.UpdateEvaluationResult(result)
	if err != nil {
		shareEvaluationResultLogger.ErrorFormat(err, "UpdateEvaluationResult 更新评测记录 异常 result=%s", _json.GetJsonString(result))
	} else {
		srv.refreshEvaluationResult(result.UID)
	}
	return err
}

// GetEvaluationResult 获取评测结果
func (srv *EvaluationResultService) GetEvaluationResult(uid string) (result *myoptional_model.EvaluationResult, err error) {
	switch config.CurrentConfig.ReadDB_MyOptional {
	case config.ReadDB_JustDB:
		return srv.getEvaluationResultDB(uid)
	case config.ReadDB_CacheOrDB_UpdateCache:
		result, err = srv.getEvaluationResultCache(uid)
		if err == nil && result == nil {
			result, err = srv.refreshEvaluationResult(uid)
		}
		return result, err
	case config.ReadDB_RefreshCache:
		result, err = srv.refreshEvaluationResult(uid)
		return result, err
	default:
		return srv.getEvaluationResultCache(uid)
	}
}

// getEvaluationResultDB 获取评测结果-读取数据库
func (srv *EvaluationResultService) getEvaluationResultDB(uid string) (result *myoptional_model.EvaluationResult, err error) {
	result, err = srv.evaluationResultRepo.GetEvaluationResult(uid)
	if err != nil {
		shareEvaluationResultLogger.ErrorFormat(err, "GetEvaluationResult 获取评测结果 异常 uid=%s", uid)
	}
	return result, err
}

// getEvaluationResultCache 获取评测结果-读取缓存
func (srv *EvaluationResultService) getEvaluationResultCache(uid string) (result *myoptional_model.EvaluationResult, err error) {
	cacheKey := EMNET_EvaluationResult_GetEvaluationResult_CacheKey + uid
	err = srv.RedisCache.GetJsonObj(cacheKey, &result)
	if err == redis.ErrNil {
		return nil, nil
	}
	return result, err
}

// refreshEvaluationResult 获取评测结果-刷新缓存
func (srv *EvaluationResultService) refreshEvaluationResult(uid string) (result *myoptional_model.EvaluationResult, err error) {
	cacheKey := EMNET_EvaluationResult_GetEvaluationResult_CacheKey + uid
	result, err = srv.getEvaluationResultDB(uid)
	if err != nil {
		return result, err
	}
	if result != nil {
		srv.RedisCache.Set(cacheKey, _json.GetJsonString(result), EMNET_EvaluationResult_CacheSeconds)
	}
	return result, err
}
