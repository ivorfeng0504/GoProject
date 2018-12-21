package user

import (
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/user"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

type SoftDataSyncService struct {
	service.BaseService
	softDataRepo *user.SoftDataSyncRepository
}

var (
	// shareSoftDataRepo UserInfoRepository
	shareSoftDataRepo *user.SoftDataSyncRepository

	// shareSoftDataLogger 共享的Logger实例
	shareSoftDataLogger dotlog.Logger
)

const (
	SoftDataServiceName = "SoftDataServiceLogger"
)

func NewSoftDataSyncService() *SoftDataSyncService {
	softDataService := &SoftDataSyncService{
		softDataRepo: shareSoftDataRepo,
	}
	softDataService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return softDataService
}

// SoftUserDataSync 同步客户端指标等数据失
func (service *SoftDataSyncService) SoftUserDataSync(ToUID int64,FromUID int64) (int,error) {
	var mapRet []map[string]interface{}
	var retCode int64
	mapRet, err := service.softDataRepo.SoftUserDataSync(ToUID, FromUID)
	if err != nil {
		shareSoftDataLogger.ErrorFormat(err, "同步客户端指标等数据失败", mapRet)
		return -1, err
	}
	shareSoftDataLogger.InfoFormat("绑定或解绑同步客户端指标等数据 formUID：%d ToUID: %d", FromUID, ToUID)

	if len(mapRet) > 0 {
		for i, _ := range mapRet {
			if mapRet[i]["errno"] != nil {
				retCode = mapRet[i]["errno"].(int64)
			}
		}
	}
	return int(retCode), nil
}



func init() {
	protected.RegisterServiceLoader(SoftDataServiceName, SoftDataSyncServiceLoader)
}

func SoftDataSyncServiceLoader() {
	shareSoftDataRepo = user.NewSoftDataSyncRepository(protected.DefaultConfig)
	shareSoftDataLogger = dotlog.GetLogger(SoftDataServiceName)
}

