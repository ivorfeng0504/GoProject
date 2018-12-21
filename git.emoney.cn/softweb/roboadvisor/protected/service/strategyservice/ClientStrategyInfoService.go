package strategyservice

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	strategyservice_model "git.emoney.cn/softweb/roboadvisor/protected/model/strategyservice"
	strategyservice_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type ClientStrategyInfoService struct {
	service.BaseService
	clientStrategyInfoRepo *strategyservice_repo.ClientStrategyInfoRepository
}

var (
	shareClientStrategyInfoRepo   *strategyservice_repo.ClientStrategyInfoRepository
	shareClientStrategyInfoLogger dotlog.Logger
)

const (
	clientStrategyInfoServiceServiceName                        = "ClientStrategyInfoService"
	EMNET_ClientStrategyInfo_PreCacheKey                        = "ClientStrategyInfoService:"
	EMNET_ClientStrategyInfo_GetClientStrategyInfoList_CacheKey = EMNET_ClientStrategyInfo_PreCacheKey + "GetClientStrategyInfoList"
)

func NewClientStrategyInfoService() *ClientStrategyInfoService {
	service := &ClientStrategyInfoService{
		clientStrategyInfoRepo: shareClientStrategyInfoRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(clientStrategyInfoServiceServiceName, clientStrategyInfoLoader)
}

func clientStrategyInfoLoader() {
	shareClientStrategyInfoRepo = strategyservice_repo.NewClientStrategyInfoRepository(protected.DefaultConfig)
	shareClientStrategyInfoLogger = dotlog.GetLogger(clientStrategyInfoServiceServiceName)
}

// GetColumnIdByClientStrategyId 根据策略Id获取栏目Id
func (srv *ClientStrategyInfoService) GetColumnIdByClientStrategyId(clientStrategyId int) (columnId int, err error) {
	dict, err := srv.GetClientStrategyInfoDict()
	columnId, success := dict[clientStrategyId]
	if success == false {
		return columnId, errors.New("未获取到栏目Id")
	}
	return columnId, nil
}

// GetClientStrategyInfoDict 获取所有可用的策略及栏目信息字典
func (srv *ClientStrategyInfoService) GetClientStrategyInfoDict() (dict map[int]int, err error) {
	result, err := srv.GetClientStrategyInfoList(true)
	if err != nil {
		return dict, err
	}
	if result == nil {
		return dict, errors.New("未获取到栏目信息")
	}
	dict = make(map[int]int)
	for _, item := range result {
		dict[item.ClientStrategyId] = item.ColumnInfoId
	}
	return dict, err
}

// GetClientStrategyInfoList 获取所有可用的策略及栏目信息 includeParent 是否包含父级
func (srv *ClientStrategyInfoService) GetClientStrategyInfoList(includeParent bool) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	switch config.CurrentConfig.ReadDB_UserHome {
	case config.ReadDB_JustDB:
		return srv.GetClientStrategyInfoListDB(includeParent)
	case config.ReadDB_CacheOrDB_UpdateCache:
		result, err = srv.GetClientStrategyInfoListCache(includeParent)
		if err == nil && result == nil {
			result, err = srv.RefreshClientStrategyInfoList(includeParent)
		}
		return result, err
	case config.ReadDB_RefreshCache:
		result, err = srv.RefreshClientStrategyInfoList(includeParent)
		return result, err
	default:
		return srv.GetClientStrategyInfoListCache(includeParent)
	}
}

// GetClientStrategyInfoListByParentId 根据父级Id获取策略列表
func (srv *ClientStrategyInfoService) GetClientStrategyInfoListByParentId(parentId int) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	if parentId <= 0 {
		return result, nil
	}
	allList, err := srv.GetClientStrategyInfoList(true)
	if err != nil {
		return result, err
	}
	if allList == nil || len(allList) == 0 {
		return result, nil
	}
	for _, item := range allList {
		if item.ParentId == parentId {
			result = append(result, item)
		}
	}
	return result, err
}

// GetClientStrategyInfoListDB 获取所有可用的策略及栏目信息-读取数据库
func (srv *ClientStrategyInfoService) GetClientStrategyInfoListDB(includeParent bool) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	result, err = srv.clientStrategyInfoRepo.GetClientStrategyInfoList(includeParent)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "GetClientStrategyInfoList 获取所有可用的策略及栏目信息 异常")
	}
	return result, err
}

// GetClientStrategyInfoTreeDB 获取所有可用的策略及栏目-树形结构-读取数据库
func (srv *ClientStrategyInfoService) GetClientStrategyInfoTreeDB() (result []*strategyservice_model.ClientStrategyInfo, err error) {
	allList, err := srv.clientStrategyInfoRepo.GetClientStrategyInfoList(true)
	dict := make(map[int]*strategyservice_model.ClientStrategyInfo)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "GetClientStrategyInfoTreeDB 获取所有可用的策略及栏目-树形结构-读取数据库 异常")
	}
	if allList == nil || len(allList) == 0 {
		return result, err
	}
	for _, item := range allList {
		if item.ParentId == 0 {
			result = append(result, item)
			dict[item.ClientStrategyId] = item
		}
	}
	for _, item := range allList {
		if item.ParentId != 0 {
			parentStrategy := dict[item.ParentId]
			if parentStrategy == nil {
				continue
			}
			parentStrategy.Children = append(parentStrategy.Children, item)
		}
	}
	return result, err
}

// getClientStrategyInfoListCache 获取所有可用的策略及栏目信息-读取缓存
func (srv *ClientStrategyInfoService) GetClientStrategyInfoListCache(includeParent bool) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	cacheKey := EMNET_ClientStrategyInfo_GetClientStrategyInfoList_CacheKey + strconv.FormatBool(includeParent)
	err = srv.RedisCache.GetJsonObj(cacheKey, &result)
	if err == redis.ErrNil {
		return nil, nil
	}
	return result, err
}

// refreshClientStrategyInfoList 获取所有可用的策略及栏目信息-刷新缓存
func (srv *ClientStrategyInfoService) RefreshClientStrategyInfoList(includeParent bool) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	cacheKey := EMNET_ClientStrategyInfo_GetClientStrategyInfoList_CacheKey + strconv.FormatBool(includeParent)
	result, err = srv.GetClientStrategyInfoListDB(includeParent)
	if err != nil {
		return nil, err
	}
	if result != nil {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, result)
	}
	return result, err
}

// GetClientStrategyInfo 根据策略Id获取策略及栏目信息
func (srv *ClientStrategyInfoService) GetClientStrategyInfo(clientStrategyId int) (strategyInfo *strategyservice_model.ClientStrategyInfo, err error) {
	strategyInfo, err = srv.clientStrategyInfoRepo.GetClientStrategyInfo(clientStrategyId)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "GetClientStrategyInfo 根据策略Id获取策略及栏目信息 异常 clientStrategyId=%d", clientStrategyId)
	}
	return strategyInfo, err
}

// InsertClientStrategyInfo 插入一个新的策略关系
func (srv *ClientStrategyInfoService) InsertClientStrategyInfo(strategyInfo strategyservice_model.ClientStrategyInfo) (err error) {
	err = srv.clientStrategyInfoRepo.InsertClientStrategyInfo(strategyInfo)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "InsertClientStrategyInfo 插入一个新的策略关系 异常 strategyInfo=%s", _json.GetJsonString(strategyInfo))
	}
	return err
}

// DeleteClientStrategyInfoNotInList 删除不在列表中的策略关联关系
func (srv *ClientStrategyInfoService) DeleteClientStrategyInfoNotInList(clientStrategyIdList []int) (err error) {
	err = srv.clientStrategyInfoRepo.DeleteClientStrategyInfoNotInList(clientStrategyIdList)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "DeleteClientStrategyInfoNotInList 删除不在列表中的策略关联关系异常 clientStrategyIdList=%s", _json.GetJsonString(clientStrategyIdList))
	}
	return err
}

// InsertClientStrategyInfoAndColumnInfo 新增策略和栏目
func (srv *ClientStrategyInfoService) InsertClientStrategyInfoAndColumnInfo(strategyId int, strategyName string, columnDesc string, isTop bool) (err error) {
	appId, err := strconv.Atoi(config.CurrentConfig.StrategyServiceAppId)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->无效的StrategyServiceAppId   config.CurrentConfig.StrategyServiceAppId=%s", config.CurrentConfig.StrategyServiceAppId)
		return err
	}
	columnSrv := NewColumnInfoService()
	newColumn := strategyservice_model.ColumnInfo{
		ColumnName: strategyName,
		ColumnDesc: columnDesc,
		AppID:      appId,
	}

	columnInfoId, err := columnSrv.InsertColumnInfo(newColumn)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->InsertColumnInfo newColumn=%s", _json.GetJsonString(newColumn))
		return err
	}
	if columnInfoId <= 0 {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->InsertColumnInfo 无效的columnInfoId =%d", columnInfoId)
		return err
	}
	newStrategyInfo := strategyservice_model.ClientStrategyInfo{
		ColumnInfoId:       columnInfoId,
		ClientStrategyId:   strategyId,
		ClientStrategyName: strategyName,
		IsTop:              isTop,
	}
	err = srv.InsertClientStrategyInfo(newStrategyInfo)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->InsertClientStrategyInfo newStrategyInfo=%s", _json.GetJsonString(newStrategyInfo))
		return err
	}
	return err
}

// UpdateClientStrategyInfoAndColumnInfo 更新策略和栏目
func (srv *ClientStrategyInfoService) UpdateClientStrategyInfoAndColumnInfo(current *strategyservice_model.ClientStrategyInfo, strategyId int, strategyName string, columnDesc string, isTop bool) (err error) {
	columnSrv := NewColumnInfoService()
	current.ClientStrategyName = strategyName
	current.ColumnName = strategyName
	current.IsTop = isTop
	err = srv.UpdateClientStrategyInfo(*current)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->UpdateClientStrategyInfo info=%s", _json.GetJsonString(current))
		return err
	}
	err = columnSrv.UpdateColumnInfo(current.ColumnName, columnDesc, current.ColumnInfoId)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->UpdateClientStrategyInfo info=%s", _json.GetJsonString(current))
		return err
	}
	return err
}

// SyncClientStrategyInfo 定时同步策略信息
func (srv *ClientStrategyInfoService) SyncClientStrategyInfo() (err error) {
	strategySrv := NewRedisStrategyInfoService()
	strategyGroupListP := strategySrv.GetStrategyData()
	if strategyGroupListP == nil {
		err = errors.New("未查询到任何策略信息")
		return err
	}
	strategyGroupList := *strategyGroupListP
	if strategyGroupList == nil || len(strategyGroupList) == 0 {
		err = errors.New("未查询到任何策略信息")
		return err
	}
	shareClientStrategyInfoLogger.DebugFormat("SyncClientStrategyInfo 【同步到的原始策略信息为】【%s】", _json.GetJsonString(strategyGroupList))
	//记录策略Id集合
	var clientStrategyIdList []int
	//循环所有策略Id
	for _, strategyGroup := range strategyGroupList {
		clientStrategyIdList = append(clientStrategyIdList, strategyGroup.StrategyGroupId)
		if strategyGroup.StrategyList == nil || len(strategyGroup.StrategyList) == 0 {
			continue
		}
		//查询父级策略是否存在
		parentStrategy, err := srv.GetClientStrategyInfo(strategyGroup.StrategyGroupId)
		if err != nil {
			shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 查询父级策略是否存在异常->GetClientStrategyInfo ClientStrategyId=%d", strategyGroup.StrategyGroupId)
			return err
		}
		//如果父级不存在则插入
		if parentStrategy == nil {
			err = srv.InsertClientStrategyInfoAndColumnInfo(strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName, strategyGroup.StrategyGroupName, true)
			if err != nil {
				shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->InsertClientStrategyInfoAndColumnInfo 父级 StrategyId=%d StrategyName=%s", strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName)
				return err
			}
		} else {
			//如果父级已存在 则更新相关名称
			err = srv.UpdateClientStrategyInfoAndColumnInfo(parentStrategy, strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName, strategyGroup.StrategyGroupName, true)
			if err != nil {
				shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->UpdateClientStrategyInfoAndColumnInfo 父级 current=%s StrategyId=%d StrategyName=%s", _json.GetJsonString(parentStrategy), strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName)
				return err
			}
		}

		for _, strategy := range strategyGroup.StrategyList {
			clientStrategyIdList = append(clientStrategyIdList, strategy.StrategyId)
			//如果策略Id不存在  则生成对应的栏目，并记录关联关系
			//如果策略已存在 则更新策略和栏目的名称
			info, err := srv.GetClientStrategyInfo(strategy.StrategyId)
			if err != nil {
				shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->GetClientStrategyInfo ClientStrategyId=%d", strategy.StrategyId)
				return err
			}

			columnDesc := strategy.StrategyName
			if len(strategyGroup.StrategyGroupName) > 0 {
				columnDesc = strategyGroup.StrategyGroupName + "-" + strategy.StrategyName
			}

			if info == nil {
				//如果策略不存在 则新增策略和栏目，包括父级
				err = srv.InsertClientStrategyInfoAndColumnInfo(strategy.StrategyId, strategy.StrategyName, columnDesc, false)
				if err != nil {
					shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->InsertClientStrategyInfoAndColumnInfo 子级 StrategyId=%d StrategyName=%s StrategyGroupId=%d StrategyGroupName=%s", strategy.StrategyId, strategy.StrategyName, strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName)
					return err
				}
			} else {
				err = srv.UpdateClientStrategyInfoAndColumnInfo(info, strategy.StrategyId, strategy.StrategyName, columnDesc, false)
				if err != nil {
					shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->UpdateClientStrategyInfoAndColumnInfo 子级 current=%s StrategyId=%d StrategyName=%s StrategyGroupId=%d StrategyGroupName=%s", _json.GetJsonString(info), strategy.StrategyId, strategy.StrategyName, strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName)
					return err
				}
			}
		}
	}
	//删除原有的策略栏目关联关系
	err = srv.DeleteAllClientStrategyInfoRelation()
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 删除原有的策略栏目关联关系 异常->DeleteAllClientStrategyInfoRelation")
		return err
	}
	//建立新的关联关系
	for _, strategyGroup := range strategyGroupList {
		if strategyGroup.StrategyList == nil || len(strategyGroup.StrategyList) == 0 {
			continue
		}
		err = srv.InsertClientStrategyInfoRelation(strategyGroup.StrategyGroupId, 0, "")
		if err != nil {
			shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 建立新的关联关系 父级 异常->InsertClientStrategyInfoRelation strategy=%s", _json.GetJsonString(strategyGroup))
			continue
		}
		for _, strategy := range strategyGroup.StrategyList {
			err = srv.InsertClientStrategyInfoRelation(strategy.StrategyId, strategyGroup.StrategyGroupId, strategyGroup.StrategyGroupName)
			if err != nil {
				shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 建立新的关联关系 异常->InsertClientStrategyInfoRelation strategy=%s", _json.GetJsonString(strategy))
				continue
			}
		}
	}

	//删除策略集合中不存在的 策略与栏目的关联关系
	err = srv.DeleteClientStrategyInfoNotInList(clientStrategyIdList)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "SyncClientStrategyInfo 定时同步策略信息 异常->DeleteClientStrategyInfoNotInList clientStrategyIdList=%s", _json.GetJsonString(clientStrategyIdList))
	} else {
		//刷新缓存
		srv.RefreshClientStrategyInfoList(true)
		srv.RefreshClientStrategyInfoList(false)
	}
	return err
}

// UpdateClientStrategyInfo 更新策略信息
func (srv *ClientStrategyInfoService) UpdateClientStrategyInfo(strategyInfo strategyservice_model.ClientStrategyInfo) (err error) {
	err = srv.clientStrategyInfoRepo.UpdateClientStrategyInfo(strategyInfo)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "UpdateClientStrategyInfo 更新策略信息 异常 strategyInfo=%s", _json.GetJsonString(strategyInfo))
	}
	return err
}

// InsertClientStrategyInfoRelation 插入一条新的关联关系
func (srv *ClientStrategyInfoService) InsertClientStrategyInfoRelation(clientStrategyId int, parentId int, parentName string) (err error) {
	err = srv.clientStrategyInfoRepo.InsertClientStrategyInfoRelation(clientStrategyId, parentId, parentName)
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "InsertClientStrategyInfoRelation 更新策略关系表 异常 clientStrategyId=%d parentId=%d", clientStrategyId, parentId)
	}
	return err
}

// DeleteAllClientStrategyInfoRelation 删除所有原有的关联关系
func (srv *ClientStrategyInfoService) DeleteAllClientStrategyInfoRelation() (err error) {
	err = srv.clientStrategyInfoRepo.DeleteAllClientStrategyInfoRelation()
	if err != nil {
		shareClientStrategyInfoLogger.ErrorFormat(err, "DeleteAllClientStrategyInfoRelation 删除所有原有的关联关系 异常")
	}
	return err
}
