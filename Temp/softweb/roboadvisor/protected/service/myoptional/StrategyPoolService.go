package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	myoptional_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
)

type StrategyPoolService struct {
	service.BaseService
	strategyPoolRepo *myoptional_repo.StrategyPoolRepository
}

var (
	shareStrategyPoolRepo   *myoptional_repo.StrategyPoolRepository
	shareStrategyPoolLogger dotlog.Logger
)

const (
	strategyPoolServiceName        = "StrategyPoolService"
	EMNET_StrategyPool_PreCacheKey = "EMoney:MyOptional:StrategyPoolService:"
)

func NewStrategyPoolService() *StrategyPoolService {
	service := &StrategyPoolService{
		strategyPoolRepo: shareStrategyPoolRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(strategyPoolServiceName, strategyPoolServiceLoader)
}

func strategyPoolServiceLoader() {
	shareStrategyPoolRepo = myoptional_repo.NewStrategyPoolRepository(protected.DefaultConfig)
	shareStrategyPoolLogger = dotlog.GetLogger(strategyPoolServiceName)
}

// InsertStrategyPool 插入一条股池信息
func (srv *StrategyPoolService) InsertStrategyPool(model *myoptional_model.StrategyPool) (strategyPoolId int64, err error) {
	strategyPoolId, err = srv.strategyPoolRepo.InsertStrategyPool(model)
	if err != nil {
		shareStrategyPoolLogger.ErrorFormat(err, "InsertStrategyPool 插入一条股池信息 异常 model=%s", _json.GetJsonString(model))
	}
	return strategyPoolId, err
}
