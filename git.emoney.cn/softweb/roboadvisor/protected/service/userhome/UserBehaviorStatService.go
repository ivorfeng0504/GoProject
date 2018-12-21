package service

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type UserBehaviorStatService struct {
	service.BaseService
	userBehaviorStatRepo *userhome_repo.UserBehaviorStatRepository
}

var (
	shareUserBehaviorStatRepo *userhome_repo.UserBehaviorStatRepository
	shareUserBehaviorLogger   dotlog.Logger
)

const (
	UserHomeUserBehaviorStatServiceName = "UserBehaviorStatService"
)

func NewUserBehaviorStatService() *UserBehaviorStatService {
	srv := &UserBehaviorStatService{
		userBehaviorStatRepo: shareUserBehaviorStatRepo,
	}
	srv.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeStatRedisConfig)
	return srv
}

func init() {
	protected.RegisterServiceLoader(UserHomeUserBehaviorStatServiceName, userHomeUserBehaviorStatServiceLoader)
}

func userHomeUserBehaviorStatServiceLoader() {
	shareUserBehaviorStatRepo = userhome_repo.NewUserBehaviorStatRepository(protected.DefaultConfig)
	shareUserBehaviorLogger = dotlog.GetLogger(UserHomeUserBehaviorStatServiceName)
}
