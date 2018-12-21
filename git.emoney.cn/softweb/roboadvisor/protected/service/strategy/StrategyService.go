package service

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/strategy"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected/model/strategy"
	"errors"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type StrategyService struct {
	service.BaseService
	strateRepo *repository.StrategyRepository
}


var (
	// strategyRepo StrategyRepository
	strategyRepo *repository.StrategyRepository

	// strategyLogger 共享的Logger实例
	strategyLogger dotlog.Logger
)


const (
	RedisKey_Strategy = "RoboAdvisor:Strategy:Cache:"
	StrategyServiceName = "StrategyService"
)

func init() {
	protected.RegisterServiceLoader(StrategyServiceName, StrategyServiceLoader)
}

func StrategyServiceLoader() {
	strategyRepo = repository.StrategyInfoRepository(protected.DefaultConfig)
	strategyLogger = dotlog.GetLogger(StrategyServiceName)
}



func StrategyInfoService() *StrategyService {
	newsService := &StrategyService{
		strateRepo: strategyRepo,
	}
	newsService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return newsService
}

// GetStrategyInfoByStreategyId
func (service *StrategyService) GetStrategyInfoById(strategyid int) (*model.StrategyInfo,error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.GetStrategyInfoBySIdCache(strategyid)
		//return service.strateRepo.GetStrategyInfoBySId(strategyid)
	default:
		return service.GetStrategyInfoBySIdCache(strategyid)
	}
}


func  (service *StrategyService) GetStrategyInfoBySIdCache(strategyid int) (*model.StrategyInfo,error){
	if strategyid <= 0 {
		return nil, errors.New("must set staretegyid")
	}
	rediskey := RedisKey_Strategy +"GetStrategyById:" +strconv.Itoa(strategyid)

	strategyinfo :=new(model.StrategyInfo)
	jsonerr := service.RedisCache.GetJsonObj(rediskey,strategyinfo)

	if jsonerr ==nil {
		return strategyinfo,jsonerr
	}

	return strategyinfo,nil

}

// GetStrategyList
func (service *StrategyService) GetStrategyList(stratelist []*model.StrategyInfo) (error) {
	switch config.CurrentConfig.ReadDB {
	case config.ReadDB_JustDB:
		return service.strateRepo.GetStrategyList(stratelist)
	default:
		return service.strateRepo.GetStrategyList(stratelist)
	}
}
