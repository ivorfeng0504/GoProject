package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/const"
	"errors"
	"encoding/json"
	"git.emoney.cn/softweb/roboadvisor/protected/model/yqq"
)

type ExpertNews_YqqService struct {
	service.BaseService
}

var (
	ExpertNewsYqqLogger dotlog.Logger

)

const (
	CacheKey_Service_YqqAPI = _const.RedisKey_NewsPre + "YqqAPIDataSet"
	RedisCacheTime = 3*60
)

func ConExpertNews_YqqService() *ExpertNews_YqqService {
	expertNewsYqqService := &ExpertNews_YqqService{
	}
	expertNewsYqqService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return expertNewsYqqService
}

//获取直播间地址列表
func (service *ExpertNews_YqqService) GetLiveRoomList()(datastr string,err error)  {

	key := CacheKey_Service_YqqAPI+"_GetLiveRoomList"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetLiveRoomListUrl
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取直播间列表接口地址配置不正确。ExpertNews_YqqService.GetLiveRoomList()")
		return "行情接口地址配置不正确",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetLiveRoomList URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetAllTagLiveRoomInfo", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetAllTagLiveRoomInfo",errors.New("【API接口数据异常】GetAllTagLiveRoomInfo")
	}

	return body,errReturn
}


//获取热门VIP直播间
func (service *ExpertNews_YqqService) GetHotVipLiveRooms()(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetHotVipLiveRooms"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetRecommandLiveRoomListUrl
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取热门和VIP直播间。ExpertNews_YqqService.GetHotVipLiveRooms()")
		return "获取热门和VIP直播间接口地址异常",nil
	}
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetHotVipLiveRooms URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetHotVipLiveRooms", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetHotVipLiveRooms",errors.New("【API接口数据异常】GetHotVipLiveRooms")
	}
	return body,err
}


//获取直播间锦囊
func (service *ExpertNews_YqqService) GetSkilbag(Lid string)(datastr string,err error){

	if Lid=="" {
		return "",errors.New("需要指定Lid")
	}

	key := CacheKey_Service_YqqAPI+"_GetSkilbag_"+Lid
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetSilkBagListUrl + "?lid="+Lid
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取直播间锦囊列表。ExpertNews_YqqService.GetSkilbag()")
		return "获取直播间锦囊列表接口地址异常",nil
	}
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetSkilbag URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetSkilbag", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetSkilbag",errors.New("【API接口数据异常】GetSkilbag")
	}
	return body,err
}


//获取益圈圈直播平台数据统计信息
func (service *ExpertNews_YqqService) GetExpertLiveData()(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetExpertLiveData"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetExpertLiveDataUrl
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取益圈圈直播平台数据统计信息。ExpertNews_YqqService.GetExpertLiveData()")
		return "获取益圈圈直播平台数据统计信息异常",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetExpertLiveData URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetExpertLiveData", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetExpertLiveData",errors.New("【API接口数据异常】GetExpertLiveData")
	}
	return body,err
}


//按标签查询直播列表
func (service *ExpertNews_YqqService) GetTagLiveRoomInfo(Tags string)(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetTagLiveRoomInfo_"+Tags
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetTagLiveRoomInfoUrl

	apiUrl = strings.Replace(apiUrl,"{0}",Tags,-1)

	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "按标签查询直播列表。ExpertNews_YqqService.GetTagLiveRoomInfo()")
		return "按标签查询直播列表异常",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetTagLiveRoomInfo URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetTagLiveRoomInfo", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetTagLiveRoomInfo",errors.New("【API接口数据异常】GetTagLiveRoomInfo")
	}
	return body,err
}

//查询所有带标签直播列表
func (service *ExpertNews_YqqService) GetAllTagLiveRoomInfo()(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetAllTagLiveRoomInfo"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetTagLiveRoomInfoUrl

	apiUrl = strings.Replace(apiUrl,"{0}","",-1)

	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.Warn("查询所有标签直播列表。ExpertNews_YqqService.GetAllTagLiveRoomInfo()")
		return "查询所有标签直播列表异常",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetAllTagLiveRoomInfo URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetAllTagLiveRoomInfo", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetAllTagLiveRoomInfo",errors.New("【API接口数据异常】GetAllTagLiveRoomInfo")
	}
	return body,nil
}


//获取益圈圈最新活动的直播间
func (service *ExpertNews_YqqService) GetLatestRoomInfo()(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetLatestRoomInfo"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetScrollLiveRoomListUrl
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取益圈圈最新活动的直播间。ExpertNews_YqqService.GetLatestRoomInfo()")
		return "获取益圈圈最新活动的直播间异常",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetLatestRoomInfo URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetLatestRoomInfo", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetLatestRoomInfo",errors.New("【API接口数据异常】GetLatestRoomInfo")
	}
	return body,err
}


//获取益圈圈首页数据信息（最新的更新直播消息+前三热门+前三VIP直播间）
func (service *ExpertNews_YqqService) GetYQQHomeData()(datastr string,err error){

	key := CacheKey_Service_YqqAPI+"_GetYQQHomeData"
	if redisdata,cErr := service.RedisCache.GetString(key);cErr==nil&&redisdata!="" {
		return redisdata,cErr
	}

	apiUrl := config.CurrentConfig.GetYqqHomeDataUrl
	if len(apiUrl) == 0 {
		ExpertNewsYqqLogger.ErrorFormat(err, "获取益圈圈首页数据。ExpertNews_YqqService.GetYQQHomeData()")
		return "获取益圈圈首页数据",nil
	}

	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		ExpertNewsYqqLogger.ErrorFormat(errReturn, "【API接口请求异常】GetYQQHomeData URL=【%s】DATA=【%s】", apiUrl, body)
		return "【API接口请求异常】GetYQQHomeData", errReturn
	}
	_ = intervalTime
	_ = contentType

	if strings.Contains(body,"success") {
		service.RedisCache.Set(key,body,RedisCacheTime)
	}else{
		return "【API接口数据异常】GetYQQHomeData",errors.New("【API接口数据异常】GetYQQHomeData")
	}
	return body,err
}



//获取专家资讯的策略栏目的首页的所有数据的整体缓存
func (service *ExpertNews_YqqService) GetExpertYqqEntireData_RedisCache()(obj interface{},err error){

	key := CacheKey_Service_YqqAPI+"GetExpertYqqEntireData_RedisCache"
	redisdata := yqq.YqqEntire{}

	if cErr := service.RedisCache.GetJsonObj(key,&redisdata);cErr==nil {
		return redisdata,cErr
	}

	//entiredata,err:= service.GetExpertYqqEntireData_Organize()
	//return entiredata,err
	return "",errors.New("GetExpertYqqEntireData_RedisCache Empty")

}


//从接口分别组织数据
func (service *ExpertNews_YqqService) GetExpertYqqEntireData_Organize()(obj interface{},err error){

	key := CacheKey_Service_YqqAPI+"GetExpertYqqEntireData_RedisCache"
	//获取分类排行
	//获取VIP HOT Latest
	//获取最新统计
	typegroup,err := service.Sub_GetExpertYqqEntireData_Organize_TypeGroup()
	latestlive,err := service.Sub_GetExpertYqqEntireData_Organize_YqqHomeData()
	yqqstatdata,err:= service.Sub_GetExpertYqqEntireData_Organize_YqqStatData()

	entire_obj := struct {
		TypeGroup interface{};YqqHomeData interface{};YqqStat interface{}
	}{
		TypeGroup:typegroup,YqqHomeData:latestlive,YqqStat:yqqstatdata,
	}

	datajson,_ :=json.Marshal(entire_obj)
	service.RedisCache.Set(key,datajson,0)

	return entire_obj,err
}



func (service *ExpertNews_YqqService) Sub_GetExpertYqqEntireData_Organize_TypeGroup()(interface{},error){
	var array_len = 6
	back_data  := make(map[string][]yqq.YqqRoomSimple)

	//1技术交易型 2价值投资型 3波段操作型 4短线激进型
	var simret_1,simret_2,simret_3,simret_4,simretall  yqq.YqqRetSimple

	Jsonstr,err:=service.GetAllTagLiveRoomInfo()
	if err!=nil {
		return err.Error(),errors.New("sub_GetExpertYqqEntireData_Organize_TypeGroup-->GetAllTagLiveRoomInfo() return error")
	}

	err = json.Unmarshal([]byte(Jsonstr),&simretall)

	for _,v := range simretall.Data {
		if strings.Contains(v.TagNameStr,"技术交易型") {
			simret_1.Data = append(simret_1.Data,v)
		}
		if strings.Contains(v.TagNameStr,"价值投资型") {
			simret_2.Data = append(simret_2.Data,v)
		}
		if strings.Contains(v.TagNameStr,"波段操作型") {
			simret_3.Data = append(simret_3.Data,v)
		}
		if strings.Contains(v.TagNameStr,"短线激进型") {
			simret_4.Data = append(simret_4.Data,v)
		}
	}

	back_data["Technology"] = simret_1.Data[0:array_len]
	back_data["Value"] = simret_2.Data[0:array_len]
	back_data["Operate"] = simret_3.Data[0:array_len]
	back_data["ShortTerm"] =simret_4.Data[0:array_len]

	return back_data,nil

}


func (service *ExpertNews_YqqService) Sub_GetExpertYqqEntireData_Organize_YqqHomeData()(interface{},error){
	Jsonstr,err:=service.GetYQQHomeData()

	if err!=nil {
		return err.Error(),errors.New("sub_GetExpertYqqEntireData_Organize_YqqHomeData-->GetYQQHomeData() return error")
	}

	var yqqhomedata yqq.YqqHomeData
	err = json.Unmarshal([]byte(Jsonstr),&yqqhomedata)

	return yqqhomedata.Data,nil
}


func (service *ExpertNews_YqqService) Sub_GetExpertYqqEntireData_Organize_YqqStatData()(interface{},error){
	Jsonstr,err:=service.GetExpertLiveData()

	if err!=nil {
		return err.Error(),errors.New("sub_GetExpertYqqEntireData_Organize_YqqStatData-->GetExpertLiveData() return error")
	}

	var statobj  yqq.YqqStatRet
	err = json.Unmarshal([]byte(Jsonstr),&statobj)

	return statobj.Data,nil

}

