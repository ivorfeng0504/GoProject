package service

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type AwardService struct {
	service.BaseService
	awardRepo *userhome_repo.AwardRepository
}

var (
	shareAwardRepo   *userhome_repo.AwardRepository
	shareAwardLogger dotlog.Logger
)

const (
	UserHomeAwardServiceName          = "AwardService"
	EMNET_UserHome_Award_CacheSeconds = 60 * 30
)

func NewAwardService() *AwardService {
	srv := &AwardService{
		awardRepo: shareAwardRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeAwardServiceName, userHomeAwardServiceLoader)
}

func userHomeAwardServiceLoader() {
	shareAwardRepo = userhome_repo.NewAwardRepository(protected.DefaultConfig)
	shareAwardLogger = dotlog.GetLogger(UserHomeAwardServiceName)
}
