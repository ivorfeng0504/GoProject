package myoptional

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	myoptional_repo "git.emoney.cn/softweb/roboadvisor/protected/repository/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/stock3minute"
	myoptional_vm "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/strings"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/mapper"
	"github.com/garyburd/redigo/redis"
	"sort"
	"strconv"
	"strings"
	"time"
)

type StockTalkService struct {
	service.BaseService
	stockTalkRepo *myoptional_repo.StockTalkRepository
}

var (
	shareStockTalkRepo   *myoptional_repo.StockTalkRepository
	shareStockTalkLogger dotlog.Logger
)

const (
	stockTalkServiceName                                        = "StockTalkService"
	EMNET_StockTalk_PreCacheKey                                 = "EMoney:MyOptional:StockTalkBll:"
	EMNET_StockTalk_GetStockListByDate_CacheKey                 = EMNET_StockTalk_PreCacheKey + "GetStockListByDate:"
	EMNET_StockTalk_GetStockListLastDate_CacheKey               = EMNET_StockTalk_PreCacheKey + "GetStockListLastDate:"
	EMNET_StockTalk_GetStockCodeList_CacheKey                   = EMNET_StockTalk_PreCacheKey + "GetStockCodeList:"
	EMNET_StockTalk_GetStockTalkNewstCacheKey_CacheKey          = EMNET_StockTalk_PreCacheKey + "GetStockTalkNewst:"
	EMNET_StockTalk_RefreshDetailToHSetFull_CacheKey            = EMNET_StockTalk_PreCacheKey + "EMNET_StockTalk_RefreshDetailToHSetFull"
	EMNET_StockTalk_RefreshDetailToHSetStockTop20_CacheKey      = EMNET_StockTalk_PreCacheKey + "EMNET_StockTalk_RefreshDetailToHSetStockTop20"
	EMNET_StockTalk_RefreshDetailToHSetStockTop50_CacheKey      = EMNET_StockTalk_PreCacheKey + "EMNET_StockTalk_RefreshDetailToHSetStockTop50"
	EMNET_StockTalk_RefreshStockTalkIdToSortedSetStock_CacheKey = EMNET_StockTalk_PreCacheKey + "EMNET_StockTalk_RefreshStockTalkIdToSortedSetStock"
)

func NewStockTalkService() *StockTalkService {
	service := &StockTalkService{
		stockTalkRepo: shareStockTalkRepo,
	}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(stockTalkServiceName, stockTalkServiceLoader)
}

func stockTalkServiceLoader() {
	shareStockTalkRepo = myoptional_repo.NewStockTalkRepository(protected.DefaultConfig)
	shareStockTalkLogger = dotlog.GetLogger(stockTalkServiceName)
}

// InsertStockTalk 插入一个微股吧评论
func (srv *StockTalkService) InsertStockTalk(model *myoptional_model.StockTalk) (err error) {
	if model == nil {
		return errors.New("评论不能为空")
	}
	model.NickName = _strings.Trim(model.NickName)
	model.StockCode = _strings.Trim(model.StockCode)
	model.StockName = _strings.Trim(model.StockName)
	model.Content = _strings.Trim(model.Content)
	if len(model.NickName) == 0 {
		return errors.New("昵称不能为空")
	}
	if len(model.UID) == 0 {
		return errors.New("UID不能为空")
	}
	if len(model.StockCode) == 0 {
		return errors.New("股票代码")
	}
	if len(model.StockName) == 0 {
		return errors.New("股票名称")
	}
	if len(model.Content) == 0 {
		return errors.New("评论内容不能为空")
	}
	err = srv.stockTalkRepo.InsertStockTalk(model)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "InsertStockTalk 插入一个微股吧评论 异常 model=%s", _json.GetJsonString(model))
	}
	return err
}

// InsertStockTalkWithAdmin 管理员插入一个微股吧评论
func (srv *StockTalkService) InsertStockTalkWithAdmin(model *myoptional_model.StockTalk) (stockTalkId int64, err error) {
	return srv.insertStockTalkWithAdmin(model, true)
}

// insertStockTalkWithAdmin 管理员插入一个微股吧评论
func (srv *StockTalkService) insertStockTalkWithAdmin(model *myoptional_model.StockTalk, refreshCache bool) (stockTalkId int64, err error) {
	if model == nil {
		return stockTalkId, errors.New("评论不能为空")
	}
	model.NickName = _strings.Trim(model.NickName)
	model.StockCode = _strings.Trim(model.StockCode)
	model.StockName = _strings.Trim(model.StockName)
	model.Content = _strings.Trim(model.Content)
	if len(model.NickName) == 0 {
		return stockTalkId, errors.New("昵称不能为空")
	}
	if len(model.StockCode) == 0 {
		return stockTalkId, errors.New("股票代码")
	}
	if len(model.StockName) == 0 {
		return stockTalkId, errors.New("股票名称")
	}
	if len(model.Content) == 0 {
		return stockTalkId, errors.New("评论内容不能为空")
	}
	stockTalkId, err = srv.stockTalkRepo.InsertStockTalkWithAdmin(model)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "InsertStockTalkWithAdmin 管理员插入一个微股吧评论 异常 model=%s", _json.GetJsonString(model))
	} else {
		if refreshCache {
			//插入成功 刷新缓存
			srv.RefreshDetailToHSetFull(stockTalkId)
			srv.RefreshDetailToHSetStockTop20(model.StockCode)
			srv.RefreshDetailToHSetStockTop50(model.StockCode)
			srv.RefreshStockTalkIdToSortedSetStock(stockTalkId)
			srv.RefreshStockCodeList()
		}
	}
	return stockTalkId, err
}

// RefreshDetailToHSetFull 刷新单个详情
func (srv *StockTalkService) RefreshDetailToHSetFull(stockTalkId int64) (stockTalk *myoptional_model.StockTalk, err error) {
	if stockTalkId <= 0 {
		return
	}
	stockTalk, err = srv.stockTalkRepo.GetStockTalkById(stockTalkId)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshDetailToHSetFull 刷新单个详情 异常 stockTalkId=%d", stockTalkId)
		return
	}
	if stockTalk == nil {
		return
	}
	if stockTalk.IsDeleted == false && stockTalk.IsValid == true {
		err = srv.RedisCache.HSet(EMNET_StockTalk_RefreshDetailToHSetFull_CacheKey, strconv.FormatInt(stockTalkId, 10), _json.GetJsonString(stockTalk))
	} else {
		_, err = srv.RedisCache.HDel(EMNET_StockTalk_RefreshDetailToHSetFull_CacheKey, strconv.FormatInt(stockTalkId, 10))
	}
	return
}

// RefreshDetailToHSetStockTop20 根据股票代码刷新最新的N条评论
func (srv *StockTalkService) RefreshDetailToHSetStockTop20(stockCode string) (stockTalkList []*myoptional_model.StockTalk, err error) {
	if len(stockCode) == 0 || len(stockCode) > 10 {
		return
	}
	stockTalkList, err = srv.stockTalkRepo.GetStockTalkListByStockCode(stockCode)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshDetailToHSetStockTop20 根据股票代码刷新最新的N条评论 异常 stockCode=%s", stockCode)
		return
	}
	if stockTalkList != nil && len(stockTalkList) > 0 {
		err = srv.RedisCache.HSet(EMNET_StockTalk_RefreshDetailToHSetStockTop20_CacheKey, stockCode, _json.GetJsonString(stockTalkList))
	} else {
		_, err = srv.RedisCache.HDel(EMNET_StockTalk_RefreshDetailToHSetStockTop20_CacheKey, stockCode)
	}
	return
}

// RefreshDetailToHSetStockTop50 根据股票代码刷新最新的N条评论
func (srv *StockTalkService) RefreshDetailToHSetStockTop50(stockCode string) (stockTalkList []*myoptional_model.StockTalk, err error) {
	if len(stockCode) == 0 || len(stockCode) > 10 {
		return
	}
	stockTalkList, err = srv.stockTalkRepo.GetStockTalkListByStockCode(stockCode)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshDetailToHSetStockTop50 根据股票代码刷新最新的N条评论 异常 stockCode=%s", stockCode)
		return
	}
	if stockTalkList != nil && len(stockTalkList) > 0 {
		err = srv.RedisCache.HSet(EMNET_StockTalk_RefreshDetailToHSetStockTop50_CacheKey, stockCode, _json.GetJsonString(stockTalkList))
	} else {
		_, err = srv.RedisCache.HDel(EMNET_StockTalk_RefreshDetailToHSetStockTop50_CacheKey, stockCode)
	}
	return
}

// GetStockTalkByStockCode 根据股票代码获取微股吧数据
func (srv *StockTalkService) GetStockTalkByStockCode(stockCode string, top int) (stockTalkList []*myoptional_vm.StockTalkForPCClient, err error) {
	if len(stockCode) == 0 || len(stockCode) > 10 {
		return
	}
	jsonList, err := srv.RedisCache.HMGet(EMNET_StockTalk_RefreshDetailToHSetStockTop50_CacheKey, stockCode)
	if err == redis.ErrNil || len(jsonList) == 0 || len(jsonList[0]) == 0 {
		return stockTalkList, nil
	}
	if err != nil {
		return stockTalkList, err
	}
	jsonStr := jsonList[0]
	err = _json.Unmarshal(jsonStr, &stockTalkList)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "GetStockTalkByStockCode 根据股票代码获取微股吧数据 异常 stockCode=%s top=%d", stockCode, top)
		return
	}
	if len(stockTalkList) > top {
		stockTalkList = stockTalkList[0:top]
	}
	if len(stockTalkList) > 0 {
		for _, item := range stockTalkList {
			item.DetailUrl = fmt.Sprintf(config.CurrentConfig.StockTalkPageUrl, item.StockCode)
		}
	}
	return
}

// RefreshStockTalkIdToSortedSetStock 刷新指定评论对应的有序股票评论Id集合
func (srv *StockTalkService) RefreshStockTalkIdToSortedSetStock(stockTalkId int64) (stockTalk *myoptional_model.StockTalk, err error) {
	if stockTalkId <= 0 {
		return
	}
	stockTalk, err = srv.stockTalkRepo.GetStockTalkById(stockTalkId)
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshStockTalkIdToSortedSetStock 刷新指定评论对应的有序股票评论Id集合 异常 stockTalkId=%d", stockTalkId)
		return
	}
	if stockTalk == nil {
		return
	}
	cacheKey := EMNET_StockTalk_RefreshStockTalkIdToSortedSetStock_CacheKey + stockTalk.StockCode
	member := strconv.FormatInt(stockTalkId, 10)
	timestamp := _time.GetTimestamp(time.Time(stockTalk.CreateTime))
	if stockTalk.IsDeleted == false && stockTalk.IsValid == true {
		shareStockTalkLogger.DebugFormat("RefreshStockTalkIdToSortedSetStock 刷新指定评论对应的有序股票评论Id集合 ZAdd cacheKey=%s score=%d member=%s", cacheKey, timestamp, member)
		_, err = srv.RedisCache.ZAdd(cacheKey, timestamp, member)
	} else {
		shareStockTalkLogger.DebugFormat("RefreshStockTalkIdToSortedSetStock 刷新指定评论对应的有序股票评论Id集合 ZRem cacheKey=%s score=%d member=%s", cacheKey, timestamp, member)
		_, err = srv.RedisCache.ZRem(cacheKey, member)
	}
	return
}

// RefreshStockCodeList 刷新股票列表
func (srv *StockTalkService) RefreshStockCodeList() (stockCodeList []*model.StockInfo, err error) {
	cacheKey := EMNET_StockTalk_GetStockCodeList_CacheKey
	stockCodeList, err = srv.stockTalkRepo.GetStockCodeList()
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshStockCodeList 刷新股票列表 异常")
		return
	}
	if stockCodeList != nil && len(stockCodeList) > 0 {
		_, err = srv.RedisCache.SetJsonObj(cacheKey, stockCodeList)
	}
	return
}

// GetStockCodeList 获取股票列表
func (srv *StockTalkService) GetStockCodeList() (stockCodeList []*model.StockInfo, err error) {
	cacheKey := EMNET_StockTalk_GetStockCodeList_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &stockCodeList)
	if err == redis.ErrNil {
		return stockCodeList, nil
	}
	return
}

// GetStockCodeListForPCClient 获取股票列表-给客户端
func (srv *StockTalkService) GetStockCodeListForPCClient() (list []string, err error) {
	stockCodeList, err := srv.GetStockCodeList()
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "GetStockCodeListForPCClient 获取股票列表-给客户端 异常")
		return list, err
	}
	if stockCodeList == nil || len(stockCodeList) == 0 {
		return list, nil
	}
	for _, item := range stockCodeList {
		stockCode := item.StockCode
		if len(stockCode) == 6 && stockCode[0] != '6' {
			stockCode = "1" + stockCode
		}
		list = append(list, stockCode)
	}
	return
}

// GetStockTalkListByDate 获取指定日期的微股吧评论
func (srv *StockTalkService) GetStockTalkListByDate(date time.Time) (stockTalkList []*myoptional_model.StockTalk, err error) {
	cacheKey := EMNET_StockTalk_GetStockListByDate_CacheKey + date.Format("2006-01-02")
	err = srv.RedisCache.GetJsonObj(cacheKey, &stockTalkList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return stockTalkList, err
}

// GetStockTalkListLastDate 获取最新一天的微股吧评论
func (srv *StockTalkService) GetStockTalkListLastDate() (stockTalkList []*myoptional_model.StockTalk, err error) {
	cacheKey := EMNET_StockTalk_GetStockListLastDate_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &stockTalkList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return stockTalkList, err
}

// GetStockTalkNewst 刷新最新的60条评论
func (srv *StockTalkService) GetStockTalkNewst() (stockTalkList []*myoptional_model.StockTalk, err error) {
	cacheKey := EMNET_StockTalk_GetStockTalkNewstCacheKey_CacheKey
	err = srv.RedisCache.GetJsonObj(cacheKey, &stockTalkList)
	if err == redis.ErrNil {
		return nil, nil
	}
	return stockTalkList, err
}

// GetStockTalkList 根据股票列表获取评论数据 pageIndex分页索引 索引从0开始  pageSize分页大小
func (srv *StockTalkService) GetStockTalkListByPage(stockCodeList []string, pageIndex int, pageSize int) (stockTalkList []*myoptional_model.StockTalk, err error) {
	if len(stockCodeList) == 0 {
		return nil, nil
	}
	var stockTalkSortList myoptional_vm.StockTalkLevelList
	//记录最早的一个股票评论的时间戳
	lastTalkDict := make(map[string]time.Time)
	if stockCodeList == nil || len(stockCodeList) == 0 {
		return stockTalkList, err
	}
	var stockTalkListParams []interface{}
	stockDistinctDict := make(map[string]bool)
	for _, stockCode := range stockCodeList {
		//如果是7代码去除第一位
		if len(stockCode) > 6 {
			stockCode = stockCode[1:]
		} else if len(stockCode) == 0 {
			continue
		}
		_, isExist := stockDistinctDict[stockCode]
		if isExist == false {
			stockTalkListParams = append(stockTalkListParams, stockCode)
			stockDistinctDict[stockCode] = true
		}
	}

	if len(stockTalkListParams) == 0 {
		return nil, nil
	}

	//如果大于1000个股票 则取前1000个
	if len(stockTalkListParams) > 1000 {
		stockTalkListParams = stockTalkListParams[:1000]
	}

	//1、先从热点数据集中获取数据
	stockTalkListJsonList, err := srv.RedisCache.HMGet(EMNET_StockTalk_RefreshDetailToHSetStockTop20_CacheKey, stockTalkListParams...)
	if err != nil {
		return stockTalkList, err
	}

	if stockTalkListJsonList != nil && len(stockTalkListJsonList) > 0 {
		for _, stockTalkJson := range stockTalkListJsonList {
			if len(stockTalkJson) == 0 {
				continue
			}
			var stockTalkListTmp []*myoptional_model.StockTalk
			jsonErr := _json.Unmarshal(stockTalkJson, &stockTalkListTmp)
			if jsonErr != nil || stockTalkListTmp == nil || len(stockTalkListTmp) == 0 {
				continue
			}
			stockTalkSortList = append(stockTalkSortList, stockTalkListTmp...)
			stockTalkListSortTmp := myoptional_vm.StockTalkList(stockTalkListTmp)
			sort.Sort(stockTalkListSortTmp)
			lastTalk := stockTalkListSortTmp[len(stockTalkListSortTmp)-1]
			lastTalkDict[lastTalk.StockCode] = time.Time(lastTalk.CreateTime)
		}
	}
	//2、排序后分页
	sort.Sort(stockTalkSortList)
	length := len(stockTalkSortList)
	startIndex := pageIndex * pageSize
	endIndex := (pageIndex + 1) * pageSize
	if startIndex < length {
		if endIndex >= length {
			endIndex = length
		}
		stockTalkSortList = stockTalkSortList[startIndex:endIndex]
		return stockTalkSortList, err
	}
	secondPageIndex := length / pageSize
	if length%pageSize > 0 {
		secondPageIndex++
	}

	//3、如果当前页已经没有数据，则从分别从股票评论排序集合中取出指定期限的Id集合
	var stockTalkIdList []interface{}
	var stockTalkSimpleList myoptional_vm.StockTalkSimpleList
	for _, stockCode := range stockTalkListParams {
		start := "(" + strconv.FormatInt(time.Now().UTC().AddDate(0, -3, 0).UnixNano()/1000000, 10)
		end := "(" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
		startTime, isExist := lastTalkDict[stockCode.(string)]
		if isExist {
			end = "(" + strconv.FormatInt(startTime.UnixNano()/1000000, 10)
			_ = startTime
		}
		cacheKey := EMNET_StockTalk_RefreshStockTalkIdToSortedSetStock_CacheKey + stockCode.(string)
		result, err := srv.RedisCache.ZRangeByScore(cacheKey, start, end, true)
		if err != nil || result == nil || len(result) < 2 {
			continue
		}
		simpleLen := len(result) / 2
		for i := 0; i < simpleLen; i++ {
			stockTalkSimpleList = append(stockTalkSimpleList, &myoptional_vm.StockTalkSimple{
				StockTalkId: result[i*2],
				CreateTime:  result[i*2+1],
			})
		}

	}
	//4、排序后分页
	sort.Sort(stockTalkSimpleList)
	length2 := len(stockTalkSimpleList)
	startIndex2 := (pageIndex - secondPageIndex) * pageSize
	endIndex2 := (pageIndex - secondPageIndex + 1) * pageSize
	if startIndex2 < length2 {
		if endIndex2 >= length2 {
			endIndex2 = length2
		}
	} else {
		return nil, nil
	}
	stockTalkSimpleList = stockTalkSimpleList[startIndex2:endIndex2]
	if len(stockTalkSimpleList) == 0 {
		return nil, nil
	}
	for _, item := range stockTalkSimpleList {
		stockTalkIdList = append(stockTalkIdList, item.StockTalkId)
	}
	//5、批量获取详情信息
	stockJsonList, err := srv.RedisCache.HMGet(EMNET_StockTalk_RefreshDetailToHSetFull_CacheKey, stockTalkIdList...)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err == nil {
		//反序列化
		jsonStr := "[" + strings.Join(stockJsonList, ",") + "]"
		err = _json.Unmarshal(jsonStr, &stockTalkList)
	}
	return stockTalkList, err
}

// AutoInsertStrategyStockTalk 自动同步策略股池信息 发送到微股吧消息中
func (srv *StockTalkService) AutoInsertStrategyStockTalk(strategyStockTalkMap map[string]*myoptional_model.StrategyStockTalk) (err error) {
	if strategyStockTalkMap == nil || len(strategyStockTalkMap) == 0 {
		return
	}
	for _, strategyStockTalk := range strategyStockTalkMap {
		//去除空字符
		strategyStockTalk.NickName = _strings.Trim(strategyStockTalk.NickName)
		strategyStockTalk.StockCode = _strings.Trim(strategyStockTalk.StockCode)
		strategyStockTalk.StockName = _strings.Trim(strategyStockTalk.StockName)
		//检查数据完整性
		if len(strategyStockTalk.NickName) == 0 || len(strategyStockTalk.StockCode) == 0 || len(strategyStockTalk.StockName) == 0 || strategyStockTalk.StrategyDescList == nil || len(strategyStockTalk.StrategyDescList) == 0 {
			continue
		}
		content := "入选" + strings.Join(strategyStockTalk.StrategyDescList, "、") + "策略"
		stockTalk := &myoptional_model.StockTalk{
			NickName:      strategyStockTalk.NickName,
			StockCode:     strategyStockTalk.StockCode,
			StockName:     strategyStockTalk.StockName,
			Content:       content,
			CreateTime:    mapper.JSONTime(time.Now()),
			IsValid:       true,
			IsDeleted:     false,
			ModifyTime:    mapper.JSONTime(time.Now()),
			ModifyUser:    "智盈Task",
			TalkLevel:     0,
			StockTalkType: _const.StockTalkType_StrategyStockPool,
		}
		stockTalkId, err := srv.InsertStockTalkWithAdmin(stockTalk)
		if err != nil {
			shareStockTalkLogger.ErrorFormat(err, "AutoInsertStrategyStockTalk 自动同步策略股池信息 发送到微股吧消息中 异常 stockTalk=%s", _json.GetJsonString(stockTalk))
			continue
		} else {
			shareStockTalkLogger.DebugFormat("AutoInsertStrategyStockTalk 自动同步策略股池信息 发送到微股吧消息中 成功 stockTalkId=%d stockTalk=%s", stockTalkId, _json.GetJsonString(stockTalk))
		}
	}
	return err
}

// ProcessRepeatStockTalk 处理重复的微股吧数据  逻辑删除
func (srv *StockTalkService) ProcessRepeatStockTalk() (err error) {
	repeatList, err := srv.stockTalkRepo.GetRepeatStockTalk()
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "ProcessRepeatStockTalk 处理重复的微股吧数据  逻辑删除 异常")
		return err
	}
	shareStockTalkLogger.DebugFormat("ProcessRepeatStockTalk 处理重复的微股吧数据  逻辑删除 repeatList=%s", _json.GetJsonString(repeatList))

	if repeatList == nil || len(repeatList) == 0 {
		return nil
	}
	for _, item := range repeatList {
		if item != nil && item.RepeatCount > 1 {
			n, deleteErr := srv.stockTalkRepo.DeleteRepeatStockTalk(item)
			if deleteErr != nil {
				shareStockTalkLogger.ErrorFormat(deleteErr, "ProcessRepeatStockTalk->DeleteRepeatStockTalk 处理重复的微股吧数据  逻辑删除 异常 item=%s", _json.GetJsonString(item))
				continue
			}
			shareStockTalkLogger.DebugFormat("ProcessRepeatStockTalk->DeleteRepeatStockTalk 处理重复的微股吧数据  逻辑删除成功  item=%s n=%d", _json.GetJsonString(item), n)
		}
	}
	srv.RefreshStockCodeList()
	//刷新缓存 从缓存移除被标记为IsDeleted=true的数据
	_, _, err = srv.RefreshDeletedStockTalkListToday()
	return err
}

// SyncStockTalkForStockBasicInfo 同步最新的基本面信息到微股吧
func (srv *StockTalkService) SyncStockTalkForStockBasicInfo() (err error) {
	//获取码表
	stockList, err := stock3minute.GetStockListInfo()
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfo 同步最新的基本面信息到微股吧 异常 获取码表失败")
		return err
	}
	return srv.SyncStockTalkForStockBasicInfoByStockList(stockList)
}

// SyncStockTalkForStockBasicInfoByStockList 根据指定的股票列表同步最新的基本面信息到微股吧
func (srv *StockTalkService) SyncStockTalkForStockBasicInfoByStockList(stockList []*model.StockInfo) (err error) {
	if stockList == nil || len(stockList) == 0 {
		err = errors.New("码表为空")
		shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 码表为空")
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 开始执行 开始时间:【%s】 StockList=【%s】", now, _json.GetJsonString(stockList))

	//同步成功的股票代码集合
	var successStockCodeList []string
	//新增的微股吧消息主键Id集合
	var successStockTalkIdList []int64
	//同步失败的股票代码集合
	var failedStockCodeList []string
	//遍历所有股票代码 获取基本面信息
	for _, stock := range stockList {
		if len(stock.StockCode) < 0 {
			continue
		}
		// 1.获取个股三分钟基本面信息
		stockThreeMinuteInfo, err := stock3minute.GetStockThreeMinuteInfo(stock.StockCode)
		if err != nil {
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 获取个股三分钟基本面信息异常 StockCode=%s", stock.StockCode)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}
		if stockThreeMinuteInfo == nil || len(stockThreeMinuteInfo.OverviewValue) == 0 {
			err = errors.New("基本面信息获取失败")
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 基本面信息获取失败 StockCode=%s", stock.StockCode)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}
		// 2.获取财务信息
		financeData, err := dataapi.GetStockQuoteAndFinance(stock.StockCode)
		if err != nil {
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 获取财务信息异常 StockCode=%s", stock.StockCode)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}
		if financeData == nil {
			err = errors.New("财务信息为空")
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 财务信息为空 StockCode=%s", stock.StockCode)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}

		// 3.重新组合数据
		stockTalkList, err := srv.ComposeStockTalk(stockThreeMinuteInfo, financeData)
		if err != nil {
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 ComposeStockTalk 组建微股吧信息异常 StockCode=%s stockThreeMinuteInfo=%s  financeData=%s", stock.StockCode, stockThreeMinuteInfo, financeData)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}
		if stockTalkList == nil || len(stockTalkList) == 0 {
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 ComposeStockTalk 重组结果为空 StockCode=%s stockThreeMinuteInfo=%s  financeData=%s", stock.StockCode, stockThreeMinuteInfo, financeData)
			continue
		}

		// 4.数据准备完毕之后 删除原有数据
		_, err = srv.stockTalkRepo.DeleteStockTalk(stock.StockCode, _const.StockTalkType_StockBasicInfo)
		if err != nil {
			shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 删除旧数据异常 StockCode=%s", stock.StockCode)
			failedStockCodeList = append(failedStockCodeList, stock.StockCode)
			continue
		}

		//防止重复记录股票代码
		isSuccessCodeRecord := false

		// 5.循环插入微股吧数据
		for _, stockInfo := range stockTalkList {
			shareStockTalkLogger.InfoFormat("构建好的微股吧数据：%s", _json.GetJsonString(stockInfo))
			//插入新数据 不刷新缓存 后面一次性刷新
			stockTalkId, err := srv.insertStockTalkWithAdmin(stockInfo, false)
			if err != nil {
				shareStockTalkLogger.ErrorFormat(err, "SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 异常 插入微股吧数据异常 StockCode=%s stockInfo=%s", stock.StockCode, _json.GetJsonString(stockInfo))
				failedStockCodeList = append(failedStockCodeList, stock.StockCode)
				continue
			}
			if isSuccessCodeRecord == false {
				successStockCodeList = append(successStockCodeList, stock.StockCode)
			}
			successStockTalkIdList = append(successStockTalkIdList, stockTalkId)
			shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 【成功】 StockCode=【%s】 StockTalk=【%s】", stock.StockCode, _json.GetJsonString(stockInfo))
			isSuccessCodeRecord = true
		}
	}

	//记录所有成功和失败的列表
	shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 【执行完毕】 执行成功的股票代码=【%s】", _json.GetJsonString(successStockCodeList))
	shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 【执行完毕】 执行失败的股票代码=【%s】", _json.GetJsonString(failedStockCodeList))
	shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 同步最新的基本面信息到微股吧 【执行完毕】 执行成功的微股吧数据集=【%s】", _json.GetJsonString(successStockTalkIdList))

	// 6.任务处理完毕后 清理数据 刷新缓存

	//刷新缓存-码表
	srv.RefreshStockCodeList()
	//刷新缓存-从缓存移除被标记为IsDeleted=true的数据
	deleteStockCodeDict, _, err := srv.RefreshDeletedStockTalkListToday()
	//刷新缓存-所有股票对应的消息
	for _, successStockCode := range stockList {
		//RefreshDeletedStockTalkListToday中已经刷新过的不用刷新
		if deleteStockCodeDict[successStockCode.StockCode] == false {
			srv.RefreshDetailToHSetStockTop20(successStockCode.StockCode)
			srv.RefreshDetailToHSetStockTop50(successStockCode.StockCode)
		}
	}

	//刷新缓存-所有新增的消息
	for _, successStockTalkId := range successStockTalkIdList {
		srv.RefreshDetailToHSetFull(successStockTalkId)
		srv.RefreshStockTalkIdToSortedSetStock(successStockTalkId)
	}
	shareStockTalkLogger.DebugFormat("SyncStockTalkForStockBasicInfoByStockList 执行完毕 开始时间:【%s】 StockList=【%s】", now, _json.GetJsonString(stockList))
	return err
}

// RefreshDeletedStockTalkListToday 刷新缓存 从缓存移除被标记为IsDeleted=true的数据
func (srv *StockTalkService) RefreshDeletedStockTalkListToday() (stockCodeDict map[string]bool, deletedStockTalkIdDict map[int64]bool, err error) {
	stockCodeDict = make(map[string]bool)
	deletedStockTalkIdDict = make(map[int64]bool)
	deletedList, err := srv.stockTalkRepo.GetDeletedStockTalkListToday()
	if err != nil {
		shareStockTalkLogger.ErrorFormat(err, "RefreshDeletedStockTalkListToday 刷新缓存 从缓存移除被标记为IsDeleted=true的数据 ->执行GetDeletedStockTalkListToday异常")
		return stockCodeDict, deletedStockTalkIdDict, err
	}
	if deletedList != nil && len(deletedList) > 0 {
		for _, item := range deletedList {
			srv.RefreshDetailToHSetFull(item.StockTalkId)
			//已经刷新过则忽略 减少刷新次数
			if stockCodeDict[item.StockCode] == false {
				srv.RefreshDetailToHSetStockTop20(item.StockCode)
				srv.RefreshDetailToHSetStockTop50(item.StockCode)
			}
			srv.RefreshStockTalkIdToSortedSetStock(item.StockTalkId)
			stockCodeDict[item.StockCode] = true
			deletedStockTalkIdDict[item.StockTalkId] = true
		}
	}
	return stockCodeDict, deletedStockTalkIdDict, err
}

// ComposeStockTalk 组建微股吧信息
func (srv *StockTalkService) ComposeStockTalk(stockBasicInfo *stock3minute.StockThreeMinuteInfo, financeData *dataapi.StockFinanceData) (stockTalkList []*myoptional_model.StockTalk, err error) {
	temp := "{OverviewValue} {TypeCap}{TypeStyle}股，处于{LifecycleValue}。估值{ScoreTTM}，{RenderPETTM}成长性{ScoreGrowing}，" +
		"{RenderOperatingRevenueYoY}{RenderNetProfitYoY}盈利能力{ScoreProfit}，{RenderROE}{RenderSalesGrossMargin}，{RenderPerformanceTypeAndRange}"
	stockBasicInfo.OverviewValue = _strings.Trim(stockBasicInfo.OverviewValue)
	if len(stockBasicInfo.OverviewValue) > 0 && _strings.LastString(stockBasicInfo.OverviewValue) != "。" {
		stockBasicInfo.OverviewValue = stockBasicInfo.OverviewValue + "。"
	}
	financeData.StockCode = _strings.Trim(financeData.StockCode)
	financeData.StockName = _strings.Trim(financeData.StockName)

	if _strings.StartWith(stockBasicInfo.OverviewValue, "公司是一家") {
		stockBasicInfo.OverviewValue = _strings.SubString(stockBasicInfo.OverviewValue, 5, -1)
	}
	if _strings.StartWith(stockBasicInfo.OverviewValue, "公司是") {
		stockBasicInfo.OverviewValue = _strings.SubString(stockBasicInfo.OverviewValue, 3, -1)
	}
	if _strings.StartWith(stockBasicInfo.OverviewValue, "公司") {
		stockBasicInfo.OverviewValue = _strings.SubString(stockBasicInfo.OverviewValue, 2, -1)
	}

	temp = _strings.Replace(temp, "{OverviewValue}", stockBasicInfo.OverviewValue)
	temp = _strings.Replace(temp, "{TypeCap}", stockBasicInfo.TypeCap)
	temp = _strings.Replace(temp, "{TypeStyle}", stockBasicInfo.TypeStyle)
	temp = _strings.Replace(temp, "{LifecycleValue}", stockBasicInfo.LifecycleValue)
	temp = _strings.Replace(temp, "{ScoreTTM}", stockBasicInfo.ScoreTTM)
	temp = _strings.Replace(temp, "{RenderPETTM}", srv.RenderPETTM(financeData.PETTM))
	temp = _strings.Replace(temp, "{ScoreGrowing}", stockBasicInfo.ScoreGrowing)
	temp = _strings.Replace(temp, "{RenderPerformanceTypeAndRange}", srv.RenderPerformanceTypeAndRange(financeData.ReportPeriod, financeData.PerformanceType, financeData.PerformanceTypeRange))
	temp = _strings.Replace(temp, "{RenderOperatingRevenueYoY}", srv.RenderOperatingRevenueYoY(financeData.OperatingRevenueYoY))
	temp = _strings.Replace(temp, "{RenderNetProfitYoY}", srv.RenderNetProfitYoY(financeData.NetProfitYoY))
	temp = _strings.Replace(temp, "{ScoreProfit}", stockBasicInfo.ScoreProfit)
	temp = _strings.Replace(temp, "{RenderROE}", srv.RenderROE(financeData.ROE))
	temp = _strings.Replace(temp, "{RenderSalesGrossMargin}", srv.RenderSalesGrossMargin(financeData.SalesGrossMargin))
	lastWord := _strings.LastString(temp)
	if lastWord == "，" || lastWord == "。" {
		temp = _strings.SubString(temp, 0, _strings.StringLen(temp)-1)
	}
	//再剔除一遍 句尾可能同时存在，和。
	lastWord = _strings.LastString(temp)
	if lastWord == "，" || lastWord == "。" {
		temp = _strings.SubString(temp, 0, _strings.StringLen(temp)-1)
	}
	temp += "。"
	temp = _strings.Replace(temp, "\n", "")
	temp = _strings.Replace(temp, "\r", "")
	lines := strings.Split(temp, "{Cut}")
	//if len(lines) < 2 {
	//	err = errors.New("行数不足")
	//	return nil, err
	//}
	//倒序存储在集合中
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if len(line) == 0 {
			continue
		}
		stockTalk := &myoptional_model.StockTalk{
			NickName:      "个股3分钟君",
			CreateTime:    mapper.JSONTime(time.Now()),
			IsValid:       true,
			ModifyTime:    mapper.JSONTime(time.Now()),
			ModifyUser:    "MainTask",
			StockCode:     financeData.StockCode,
			StockName:     financeData.StockName,
			TalkLevel:     9, //置顶
			Content:       line,
			StockTalkType: _const.StockTalkType_StockBasicInfo,
		}
		stockTalkList = append(stockTalkList, stockTalk)
	}
	return stockTalkList, err
}

// RenderPETTM RenderPETTM
func (srv *StockTalkService) RenderPETTM(pettm string) string {
	str := "PE（TTM）：{PETTM}；"
	pettm = _strings.Trim(pettm)
	if len(pettm) == 0 {
		return ""
	}
	str = _strings.Replace(str, "{PETTM}", _strings.RoundStrWithNoError(pettm, 0))
	return str
}

// RenderPerformanceTypeAndRange Render最新业绩预告
func (srv *StockTalkService) RenderPerformanceTypeAndRange(reportPeriod, performanceType, performanceTypeRange string) string {
	reportPeriod = _strings.Trim(reportPeriod)
	performanceType = _strings.Trim(performanceType)
	performanceTypeRange = _strings.Trim(performanceTypeRange)
	if len(reportPeriod) == 0 || len(performanceType) == 0 || len(performanceTypeRange) == 0 {
		return ""
	}
	date, err := _time.ParseTime(reportPeriod)
	if err != nil {
		return ""
	}
	if date.Month() <= 3 {
		reportPeriod = "一季度"
	} else if date.Month() <= 6 {
		reportPeriod = "上半年"
	} else if date.Month() <= 9 {
		reportPeriod = "前三季度"
	} else {
		reportPeriod = "本年度"
	}
	str := "预计{ReportPeriod}业绩{PerformanceType}{PerformanceTypeRange}，"
	str = _strings.Replace(str, "{PerformanceType}", performanceType)
	str = _strings.Replace(str, "{PerformanceTypeRange}", performanceTypeRange)
	str = _strings.Replace(str, "{ReportPeriod}", reportPeriod)
	return str
}

// RenderOperatingRevenueYoY Render本期营业收入收增长率
func (srv *StockTalkService) RenderOperatingRevenueYoY(operatingRevenueYoY string) string {
	str := "当期营收增长：{OperatingRevenueYoY}，"
	operatingRevenueYoY = _strings.Trim(operatingRevenueYoY)
	if len(operatingRevenueYoY) == 0 {
		return ""
	}
	str = _strings.Replace(str, "{OperatingRevenueYoY}", _strings.RoundStrWithNoError(operatingRevenueYoY, 1)+"%")
	return str
}

// RenderNetProfitYoY Render本期净利润增长率
func (srv *StockTalkService) RenderNetProfitYoY(netProfitYoY string) string {
	str := "当期利润增长：{NetProfitYoY}；"
	netProfitYoY = _strings.Trim(netProfitYoY)
	if len(netProfitYoY) == 0 {
		return ""
	}
	str = _strings.Replace(str, "{NetProfitYoY}", _strings.RoundStrWithNoError(netProfitYoY, 1)+"%")
	return str
}

// RenderROE Render本期净资产收益率
func (srv *StockTalkService) RenderROE(roe string) string {
	str := "净资产收益率：{ROE}，"
	roe = _strings.Trim(roe)
	if len(roe) == 0 {
		return ""
	}
	str = _strings.Replace(str, "{ROE}", _strings.RoundStrWithNoError(roe, 1)+"%")
	return str
}

// RenderSalesGrossMargin Render本期毛利率
func (srv *StockTalkService) RenderSalesGrossMargin(salesGrossMargin string) string {
	str := "毛利率：{SalesGrossMargin}"
	salesGrossMargin = _strings.Trim(salesGrossMargin)
	if len(salesGrossMargin) == 0 {
		return ""
	}
	str = _strings.Replace(str, "{SalesGrossMargin}", _strings.RoundStrWithNoError(salesGrossMargin, 1)+"%")
	return str
}
