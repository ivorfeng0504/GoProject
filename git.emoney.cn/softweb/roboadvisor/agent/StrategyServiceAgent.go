package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	strategyservice_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/cache"
	"strconv"
	"time"
)

// GetClientStrategyInfoList 获取所有可用的策略及栏目信息
func GetClientStrategyInfoList(includeParent bool) (result []*strategyservice_vmmodel.ClientStrategyInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = includeParent
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategylist", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetColumnIdDict 根据策略Id获取栏目Id字典
func GetColumnIdDict() (dict map[int]int, err error) {
	dict = make(map[int]int)
	result, err := GetClientStrategyInfoList(true)
	if err != nil {
		return dict, err
	}
	if result == nil {
		return dict, errors.New("未获取到栏目字典")
	}
	for _, item := range result {
		dict[item.ClientStrategyId] = item.ColumnInfoId
	}
	return dict, nil
}

// GetColumnIdByClientStrategyId 根据策略Id获取栏目Id
func GetColumnIdByClientStrategyId(clientStrategyId int) (columnId int, err error) {
	result, err := GetClientStrategyInfoList(true)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, errors.New("未获取到栏目Id")
	}
	dict := make(map[int]int)
	for _, item := range result {
		dict[item.ClientStrategyId] = item.ColumnInfoId
	}
	columnId, success := dict[clientStrategyId]
	if success == false {
		return 0, errors.New("未获取到栏目Id")
	}
	return columnId, nil
}

// GetStrategyNameByColumnId 根据栏目Id与策略名称的字典
func GetColumnIdStrategyNameDict() (dict map[string]string, err error) {
	result, err := GetClientStrategyInfoList(true)
	if err != nil {
		return dict, err
	}
	if result == nil {
		return dict, errors.New("未获取到栏目信息")
	}
	dict = make(map[string]string)
	for _, item := range result {
		dict[strconv.Itoa(item.ColumnInfoId)] = item.ClientStrategyName
	}

	return dict, nil
}

// GetStrategyLiveRoomId 根据策略Id获取直播间Id
func GetStrategyLiveRoomId(clientStrategyId string) (roomId int, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategyliveroomdict", req)
	if err != nil {
		return roomId, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return roomId, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return roomId, nil
	}
	dict := make(map[string]string)
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &dict)
	if err != nil {
		return roomId, err
	}
	roomIdStr, success := dict[clientStrategyId]
	if success == false {
		return roomId, errors.New("未获取到直播间Id")
	}
	roomId, err = strconv.Atoi(roomIdStr)
	if err != nil {
		return roomId, err
	}
	return roomId, nil
}

// GetStrategyNewsList 根据栏目id和资讯类型获取资讯列表
func GetStrategyNewsList(colID int, newsType int, pageIndex int64, pageSize int64) (result []*strategyservice_vmmodel.NewsInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsRequest{
		ColumnId:  colID,
		NewsType:  newsType,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategynewslist", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetStrategyNewsDetailList 根据栏目id和资讯类型批量获取资讯列表
func GetStrategyNewsDetailList(newsIDStr []string) (result []*strategyservice_vmmodel.NewsInfo, err error) {
	if newsIDStr == nil || len(newsIDStr) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	req.RequestData = newsIDStr
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategynewsdetaillist", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	//处理来源
	ProcessFrom(result)
	return result, nil
}

// GetIndexStrategyNewsV1 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV1(columnId int) (result IndexStrategyNewsResponse, err error) {
	result = IndexStrategyNewsResponse{}
	strategyNewsList, err := GetStrategyNewsList(columnId, _const.NewsType_News, 1, 1)
	if err != nil {
		return result, err
	}
	if strategyNewsList != nil && len(strategyNewsList) > 0 {
		result.StrategyNews = strategyNewsList[0]
	}
	strategyNewsVideoList, err := GetStrategyNewsList(columnId, _const.NewsType_Video, 1, 1)
	if err != nil {
		return result, err
	}
	if strategyNewsVideoList != nil && len(strategyNewsVideoList) > 0 {
		result.StrategyVideoNews = strategyNewsVideoList[0]
	}
	return result, err
}

// GetIndexStrategyNewsV2 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV2(clientStrategyGroupId int) (result IndexStrategyNewsResponse, err error) {
	if clientStrategyGroupId == 0 {
		return result, nil
	}
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsInfoRequest{
		ClientStrategyGroupId: clientStrategyGroupId,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getindexstrategynews", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return result, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetIndexStrategyNewsV3 获取首页最新的一条策略资讯与策略视频
func GetIndexStrategyNewsV3(clientStrategyId int, clientStrategyGroupId int) (result *IndexStrategyNewsResponse, err error) {
	runCache := cache.GetRuntimeCache()
	cacheKey := "runtime_GetIndexStrategyNewsV3" + strconv.Itoa(clientStrategyId) + ":" + strconv.Itoa(clientStrategyGroupId)
	cacheData, err := runCache.Get(cacheKey)
	if err == nil && cacheData != nil {
		result = cacheData.(*IndexStrategyNewsResponse)
	}
	if result != nil {
		return result, nil
	}
	result = &IndexStrategyNewsResponse{}
	dict, err := GetColumnIdDict()
	if err != nil {
		return result, err
	}
	columnId := dict[clientStrategyId]
	groupColumnId := dict[clientStrategyGroupId]
	strategyNewsList, err := GetStrategyNewsList(columnId, _const.NewsType_News, 1, 1)
	if err != nil {
		return result, err
	}
	if strategyNewsList != nil && len(strategyNewsList) > 0 {
		result.StrategyNews = strategyNewsList[0]
	} else {
		strategyNewsList, err = GetStrategyNewsList(groupColumnId, _const.NewsType_News, 1, 1)
		if err != nil {
			return result, err
		}
		if strategyNewsList != nil && len(strategyNewsList) > 0 {
			result.StrategyNews = strategyNewsList[0]
			result.StrategyNewsIsGroup = true
		}
	}
	strategyNewsVideoList, err := GetStrategyNewsList(columnId, _const.NewsType_Video, 1, 1)
	if err != nil {
		return result, err
	}
	if strategyNewsVideoList != nil && len(strategyNewsVideoList) > 0 {
		result.StrategyVideoNews = strategyNewsVideoList[0]
	} else {
		strategyNewsVideoList, err = GetStrategyNewsList(groupColumnId, _const.NewsType_Video, 1, 1)
		if err != nil {
			return result, err
		}
		if strategyNewsVideoList != nil && len(strategyNewsVideoList) > 0 {
			result.StrategyVideoNews = strategyNewsVideoList[0]
			result.StrategyVideoNewsIsGroup = true
		}
	}
	if result != nil {
		runCache.Set(cacheKey, result, 5*60)
	}
	return result, err
}

// GetHotStrategyNewsInfo 获取热门视频或者资讯
func GetHotStrategyNewsInfo(clientStrategyIdList []string, newsType int) (result []*strategyservice_vmmodel.NewsInfo, err error) {
	if clientStrategyIdList == nil || len(clientStrategyIdList) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsInfoRequest{
		ClientStrategyIdList: clientStrategyIdList,
		NewsType:             newsType,
		Count:                10,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/gethotstrategynewsinfo", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetNewstStrategyNews 获取最新策略资讯或视频
func GetNewstStrategyNews(clientStrategyIdList []string, newsType int) (result []*strategyservice_vmmodel.NewsInfo, err error) {
	if clientStrategyIdList == nil || len(clientStrategyIdList) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsInfoRequest{
		ClientStrategyIdList: clientStrategyIdList,
		NewsType:             newsType,
		Count:                15,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getnewststrategynews", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, nil
}

// GetStrategyNewsList 根据栏目id和资讯类型获取资讯列表
func GetStrategyNewsList_MultiMedia(colID int, newsType int, date *time.Time) (result []*strategyservice_vmmodel.ExpertNews_MultiMedia_List, err error) {
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsRequest{
		ColumnId:   colID,
		NewsType:   newsType,
		FilterDate: date,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategynewslist_multimedia", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	//处理来源
	ProcessFrom_MultiMediaList(result)
	return result, nil
}

// GetStrategyNewsDetailList 根据栏目id和资讯类型批量获取资讯列表
func GetStrategyNewsDetailList_MultiMedia(newsIDStr []string) (result []*strategyservice_vmmodel.ExpertNews_MultiMedia_Details, err error) {
	if newsIDStr == nil || len(newsIDStr) == 0 {
		return nil, nil
	}
	req := contract.NewApiRequest()
	req.RequestData = newsIDStr
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategynewsdetaillist_multimedia", req)
	if err != nil {
		return result, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	//处理来源
	ProcessFrom_MultiMedia(result)
	return result, nil
}

// GetNewsList_ImportantTips 根据栏目id和资讯类型获取资讯列表
func GetNewsList_ImportantTips(columnID int, StrategyID int, TagID int, pageIndex int64, pageSize int64) (result []*strategyservice_vmmodel.NewsInfo, totalCount int, err error) {
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsRequest{
		ColumnId:         columnID,
		ClientStrategyId: StrategyID,
		TagID:            TagID,
		PageIndex:        pageIndex,
		PageSize:         pageSize,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getnewslist_importanttips", req)
	if err != nil {
		return result, 0, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, 0, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, 0, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	return result, response.TotalCount, nil
}

// GetStrategyNewsList_MultiMedia1Month 获取最近一个月实战培训数据
func GetStrategyNewsList_MultiMedia1Month(pageIndex int64, pageSize int64, date *time.Time) (result []*strategyservice_vmmodel.ExpertNews_MultiMedia_TrainList, totalCount int, err error) {
	req := contract.NewApiRequest()
	req.RequestData = StrategyNewsRequest{
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		FilterDate: date,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/strategyservice/getstrategynewslist_multimedia1month", req)
	if err != nil {
		return result, 0, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return result, 0, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, 0, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &result)
	//处理来源
	ProcessFrom_MultiMediaTrainList(result)
	return result, response.TotalCount, nil
}

// ProcessFrom 处理来源显示
func ProcessFrom(result []*strategyservice_vmmodel.NewsInfo) {
	if result == nil || len(result) == 0 {
		return
	}
	dict, err := GetColumnIdStrategyNameDict()
	if err != nil {
		return
	}
	for _, item := range result {
		if len(item.ShowPosition) > 0 {
			var showPositionInfoList []*ShowPositionInfo
			err = _json.Unmarshal(item.ShowPosition, &showPositionInfoList)
			if err != nil || showPositionInfoList == nil || len(showPositionInfoList) == 0 {
				continue
			}
			colId := showPositionInfoList[0].ShowPosition
			item.From = dict[colId]
		}
	}
}

// ProcessFrom 处理来源显示
func ProcessFrom_MultiMedia(result []*strategyservice_vmmodel.ExpertNews_MultiMedia_Details) {
	if result == nil || len(result) == 0 {
		return
	}
	dict, err := GetColumnIdStrategyNameDict()
	if err != nil {
		return
	}
	for _, item := range result {
		if len(item.ShowPosition) > 0 {
			var showPositionInfoList []*ShowPositionInfo
			err = _json.Unmarshal(item.ShowPosition, &showPositionInfoList)
			if err != nil || showPositionInfoList == nil || len(showPositionInfoList) == 0 {
				continue
			}
			colId := showPositionInfoList[0].ShowPosition
			item.From = dict[colId]
		}
	}
}

// ProcessFrom 处理来源显示
func ProcessFrom_MultiMediaList(result []*strategyservice_vmmodel.ExpertNews_MultiMedia_List) {
	if result == nil || len(result) == 0 {
		return
	}
	dict, err := GetColumnIdStrategyNameDict()
	if err != nil {
		return
	}
	for _, item := range result {
		if len(item.ShowPosition) > 0 {
			var showPositionInfoList []*ShowPositionInfo
			err = _json.Unmarshal(item.ShowPosition, &showPositionInfoList)
			if err != nil || showPositionInfoList == nil || len(showPositionInfoList) == 0 {
				continue
			}
			colId := showPositionInfoList[0].ShowPosition
			item.From = dict[colId]
		}
	}
}

// ProcessFrom 处理来源显示
func ProcessFrom_MultiMediaTrainList(result []*strategyservice_vmmodel.ExpertNews_MultiMedia_TrainList) {
	if result == nil || len(result) == 0 {
		return
	}
	dict, err := GetColumnIdStrategyNameDict()
	if err != nil {
		return
	}
	for _, item := range result {
		if len(item.ShowPosition) > 0 {
			var showPositionInfoList []*ShowPositionInfo
			err = _json.Unmarshal(item.ShowPosition, &showPositionInfoList)
			if err != nil || showPositionInfoList == nil || len(showPositionInfoList) == 0 {
				continue
			}
			colId := showPositionInfoList[0].ShowPosition
			item.From = dict[colId]
		}
	}
}

type StrategyNewsRequest struct {
	//栏目Id
	ColumnId int
	//客户端策略Id
	ClientStrategyId int
	//0资讯 1视频 NewsType.go
	NewsType int
	//分页索引
	PageIndex int64
	//分页大小
	PageSize int64
	//是否包含子集策略的数据
	ContainChild bool
	//标签ID
	TagID int
	//过滤时间 2018-11-20
	FilterDate *time.Time
}

type IndexStrategyNewsResponse struct {
	StrategyNews             *strategyservice_vmmodel.NewsInfo
	StrategyVideoNews        *strategyservice_vmmodel.NewsInfo
	StrategyNewsIsGroup      bool
	StrategyVideoNewsIsGroup bool
}

type StrategyNewsInfoRequest struct {
	//客户端策略ID集合
	ClientStrategyIdList []string
	//资讯类型
	NewsType int
	//获取的数量
	Count int

	//组Id
	ClientStrategyGroupId int
}

type ShowPositionInfo struct {
	ShowPosition string
	IsTop        string
}
