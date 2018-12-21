package train

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/train"
	trainModel "git.emoney.cn/softweb/roboadvisor/protected/model/train"
	"git.emoney.cn/softweb/roboadvisor/protected/service/resapi"
	"git.emoney.cn/softweb/roboadvisor/protected/service/train"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers"
	"github.com/devfeel/dotweb"
	"sort"
	"strconv"
	"time"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/const"
	strategyservice_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/strategyservice"
)

//获取runtimcache
var rtcache_resconfig = handlers.NewWebRuntimeCache("train_", 60*60)
var rtcache = handlers.NewWebRuntimeCache("train_", 60)
var redisKey_nowTrainList = _const.RedisKey_NewsPre + "GetTrainListNow"

func TrainGuidor(ctx dotweb.Context) error {
	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)
	ctx.ViewData().Set("StaticResEnv", config.CurrentConfig.StaticResEnv)
	ctx.ViewData().Set("ServerVirtualPath", config.CurrentConfig.ServerVirtualPath)
	return ctx.View("userguidor.html")
}

// 用户培训 首页
//func TrainIndex(ctx dotweb.Context) error {
//
//	ctx.ViewData().Set("StaticServerHost", config.CurrentConfig.StaticServerHost)
//	ctx.ViewData().Set("ResourceVersion", config.CurrentConfig.ResourceVersion)
//	ctx.ViewData().Set("StaticResEnv", config.CurrentConfig.StaticResEnv)
//	ctx.ViewData().Set("ServerVirtualPath", config.CurrentConfig.ServerVirtualPath)
//
//	response := contract.NewResonseInfo()
//	user := contract_train.TrainUserInfo(ctx)
//	service := train.NewTrainService()
//
//	//判断用户是否首次使用用户培训（首次跳转宣传页）
//	cookiekey_IsFirst := "Emoney.ZYTrain"
//	rediskey_IsFirst := _const.RedisKey_NewsPre + "Train"
//
//	//cookie 判断是否首次使用
//	isfirstValue, err := ctx.ReadCookieValue(cookiekey_IsFirst)
//	if isfirstValue == "" || err != nil {
//		ctx.SetCookieValue(cookiekey_IsFirst, strconv.FormatInt(user.UID, 10), 0)
//
//		//redis 判断是否首次使用
//		isfirstValue, err = service.RedisCache.HGet(rediskey_IsFirst, strconv.FormatInt(user.UID, 10))
//		if err != nil {
//			service.GetTrainLogger().ErrorFormat(err, "TrainIndex 培训首页-是否首次点击 获取redis异常 uid:%d", user.UID)
//		}
//		if isfirstValue == "" {
//			service.RedisCache.HSet(rediskey_IsFirst, strconv.FormatInt(user.UID, 10), strconv.FormatInt(user.UID, 10))
//
//			//首次使用跳转到宣传页
//			return ctx.View("userguidor.html")
//		}
//	}
//
//	response.RetCode = 0
//	response.RetMsg = ""
//	response.Message = nil
//	ctx.ViewData().Set("TrainList", response)
//
//	trainList, err := service.GetTrainListByDate(strconv.Itoa(user.PID), time.Now())
//
//	if err != nil {
//		response.RetCode = -1
//		response.RetMsg = ""
//		response.Message = nil
//		ctx.ViewData().Set("TrainList", response)
//	}
//	sort.Sort(NetworkMeetingInfoData(trainList))
//
//	response.RetCode = 0
//	response.RetMsg = ""
//	response.Message = trainList
//
//	str, err := _json.Marshal(response)
//
//	ctx.ViewData().Set("TrainList", str)
//	ctx.ViewData().Set("NowTime", time.Now().Format("2006-01-02 15:04:05"))
//
//	return ctx.View("index.html")
//}

// 用户培训-新改版
func Index(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "user_training.html")
}

// 用户培训 当日最新（用户培训最近1个月数据+实战培训最近一个月数据）
func GetTrainListByNew(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)

	currPageStr := ctx.QueryString("currpage")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	cacheKey := fmt.Sprintf("GetTrainListByNew_%d_%d_%d", currPage, pageSize, user.PID)
	cacheKeyCount := fmt.Sprintf("GetTrainListByNewCount_%d_%d_%d", currPage, pageSize, user.PID)
	cacheObj, exist := rtcache.GetCache(cacheKey)
	cacheObjCount, _ := rtcache.GetCache(cacheKeyCount)
	if !exist || cacheObj == nil {
		service := train.NewTrainService()
		trainList, totalCount, err := service.GetTrainListByNew(strconv.Itoa(user.PID), currPage, pageSize)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "获取培训列表失败"
			response.Message = err.Error()
			return ctx.WriteJson(response)
		}
		cacheObj = trainList
		cacheObjCount = totalCount
		rtcache.SetCache(cacheKey, trainList)
		rtcache.SetCache(cacheKeyCount, totalCount)
	}
	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = cacheObj
	response.TotalCount = cacheObjCount.(int)
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// 根据标签获取课程列表
func GetTrainListByTag(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)
	trainTagStr := ctx.QueryString("trainTag")
	trainTag, err := strconv.Atoi(trainTagStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "培训标签不正确"
		return ctx.WriteJson(response)
	}

	currPageStr := ctx.QueryString("currpage")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	cacheKey := fmt.Sprintf("GetTrainListByTag_%d_%d_%d_%d", currPage, pageSize, trainTag, user.PID)
	cacheKeyCount := fmt.Sprintf("GetTrainListByTagCount_%d_%d_%d_%d", currPage, pageSize, trainTag, user.PID)
	cacheObj, exist := rtcache.GetCache(cacheKey)
	cacheObjCount, _ := rtcache.GetCache(cacheKeyCount)
	if !exist || cacheObj == nil {
		service := train.NewTrainService()
		trainList, totalCount, err := service.GetTrainListByTag(strconv.Itoa(user.PID), trainTag, currPage, pageSize)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "获取培训列表失败"
			response.Message = err.Error()
			return ctx.WriteJson(response)
		}
		cacheObj = trainList
		cacheObjCount = totalCount
		rtcache.SetCache(cacheKey, trainList)
		rtcache.SetCache(cacheKeyCount, totalCount)
	}
	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = cacheObj
	response.TotalCount = cacheObjCount.(int)
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// 根据地区获取课程列表 （老版用户培训）
func GetTrainListByDateAndArea(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)
	area := ctx.QueryString("area")
	if area == "" {
		response.RetCode = -1
		response.RetMsg = "地区不能为空"
		return ctx.WriteJson(response)
	}
	req_date := ctx.QueryString("date")
	if req_date == "" {
		response.RetCode = -1
		response.RetMsg = "日期不正确"
		return ctx.WriteJson(response)
	}
	date, err := time.Parse("2006-01-02", req_date)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "日期格式不正确"
		response.Message = err.Error()
		return ctx.WriteJson(response)
	}

	if area == "默认" {
		area = "全国"
	}

	service := train.NewTrainService()
	trainList, err := service.GetTrainListByDateAndArea(strconv.Itoa(user.PID), date, area)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	sort.Sort(NetworkMeetingInfoData(trainList))

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	return ctx.WriteJson(response)
}

// 获取用户培训 实战策略显示的策略信息
func GetTrainClientInfo(ctx dotweb.Context) error {
	cacheKey := "GetTrainClientInfo"
	cacheObj, exist := rtcache_resconfig.GetCache(cacheKey)
	if !exist {
		ret, err := resapi.GetTrainClientInfoConfig()
		if err != nil {
			return ctx.WriteJson(nil)
		}
		var clientList []*trainModel.TrainClientInfo
		err = _json.Unmarshal(ret, &clientList)
		if err != nil {
			return ctx.WriteJson(nil)
		}
		cacheObj = clientList
		rtcache_resconfig.SetCache(cacheKey, clientList)
	}

	return ctx.WriteJson(cacheObj)
}

// 获取用户培训显示的标签信息
func GetTrainTagInfo(ctx dotweb.Context) error {
	tagtype := ctx.QueryString("type")
	cacheKey := "GetTrainTagInfo1:" + tagtype
	cacheObj, exist := rtcache_resconfig.GetCache(cacheKey)
	if !exist {
		ret, err := resapi.GetTrainTagInfoConfig()
		if err != nil {
			return ctx.WriteJson(nil)
		}
		var tagList []*trainModel.TrainTagInfo
		var retList []*trainModel.TrainTagInfo
		err = _json.Unmarshal(ret, &tagList)
		if err != nil {
			return ctx.WriteJson(nil)
		}

		for _, v := range tagList {
			if v.Type == tagtype {
				retList = append(retList, v)
			}
		}

		cacheObj = retList
		rtcache_resconfig.SetCache(cacheKey, retList)
	}

	return ctx.WriteJson(cacheObj)
}



// 当日最新（支持日期筛选）
func GetTrainListByDate(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)
	nowFlag := false
	dateStr := _http.GetRequestValue(ctx, "date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		date = &now
		err = nil
		nowFlag = true
	} else {
		date = &dateTmp
	}
	trainList, err := agent.GetTrainListByDate(strconv.Itoa(user.PID), date)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	trainService := train.NewTrainService()

	if nowFlag && len(trainList) > 0 {
		obj, _ := rtcache_resconfig.GetCache(redisKey_nowTrainList)
		if obj == nil ||obj.(*time.Time).Format("2006-01-02")!=date.Format("2006-01-02"){
			trainService.RedisCache.SetJsonObj(redisKey_nowTrainList, date)
			rtcache_resconfig.SetCache(redisKey_nowTrainList, date)
		}
	}
	if nowFlag && len(trainList) == 0 {
		trainService.RedisCache.GetJsonObj(redisKey_nowTrainList, &date)
		trainList, err = agent.GetTrainListByDate(strconv.Itoa(user.PID), date)
		if err != nil {
			response.RetCode = -3
			response.RetMsg = "获取培训列表失败"
			return ctx.WriteJson(response)
		}
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// 解盘晨会和指标学习（支持日期和标签筛选）
func GetTrainListByDateAndTag(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)

	dateStr := _http.GetRequestValue(ctx, "date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = nil
		err = nil
	} else {
		date = &dateTmp
	}
	trainTagStr := ctx.QueryString("trainTag")
	trainTag, err := strconv.Atoi(trainTagStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "培训标签不正确"
		return ctx.WriteJson(response)
	}
	trainType := ctx.QueryString("trainType")
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "培训类型不正确"
		return ctx.WriteJson(response)
	}
	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	trainList, totalCount, err := agent.GetTrainListByDateAndTag(strconv.Itoa(user.PID), date, trainTag, trainType, pageSize, currpage)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取培训列表失败"
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取培训列表成功"
	response.Message = trainList
	response.TotalCount = totalCount
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// 获取最近一个月的实战培训综合课表
func GetStrategyNewsList_MultiMedia1Month(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	dateStr := _http.GetRequestValue(ctx, "date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = nil
		err = nil
	} else {
		date = &dateTmp
	}
	currpageStr := ctx.QueryString("currpage")
	currpage, err := strconv.Atoi(currpageStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "当前页码不正确"
		return ctx.WriteJson(response)
	}

	pageSizeStr := ctx.QueryString("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "pageSize不正确"
		return ctx.WriteJson(response)
	}

	newsList, totalCount, err := agent.GetStrategyNewsList_MultiMedia1Month(int64(currpage), int64(pageSize), date)
	var retNewsList []*strategyservice_vmmodel.ExpertNews_MultiMedia_TrainList

	//过滤中继回踩策略
	for _, v := range newsList {
		if v.From != _const.TrainFilterCL_ZJHC {
			retNewsList = append(retNewsList, v)
		}
	}
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取策略实战综合列表失败"
		response.Message = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "获取策略实战综合列表成功"
	response.Message = retNewsList
	response.TotalCount = totalCount
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// GetStrategyNewsList_MultiMedia 获取策略指定的实战培训（多媒体类型资讯）
func GetStrategyNewsList_MultiMedia(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	newsType := _const.NewsType_MultiMedia //多媒体资讯类型

	clientStrategyIdStr := _http.GetRequestValue(ctx, "ClientStrategyId")
	dateStr := _http.GetRequestValue(ctx, "date")
	var date *time.Time
	dateTmp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = nil
		err = nil
	} else {
		date = &dateTmp
	}
	clientStrategyId, err := strconv.Atoi(clientStrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略Id不正确"
		return ctx.WriteJson(response)
	}
	colID, err := agent.GetColumnIdByClientStrategyId(clientStrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	newsList, err := agent.GetStrategyNewsList_MultiMedia(colID, newsType, date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	response.SystemTime = time.Now()
	return ctx.WriteJson(response)
}

// 培训-开课时间排序
type NetworkMeetingInfoData []*trainModel.NetworkMeetingInfo

func (s NetworkMeetingInfoData) Len() int {
	return len(s)
}
func (s NetworkMeetingInfoData) Less(i, j int) bool {
	return s[i].Class_date.Before(s[j].Class_date)
}
func (s NetworkMeetingInfoData) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}


