package resapi

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
)

const (
	resourceConfigAPIServiceName                 = "ResourceConfigAPI"
	EMNET_GetStrategyLiveRoomConfig_CacheKey     = "ResourceConfigAPI:GetStrategyLiveRoomConfig"
	EMNET_GetClientExpertStrategyConfig_CacheKey = "ResourceConfigAPI:GetClientExpertStrategyConfig"
	EMNET_GetZYPID_CacheKey                      = "ResourceConfigAPI:GetZYPIDConfig"
	EMNET_GetTrainClientInfoConfig_CacheKey      = "ResourceConfigAPI:GetTrainClientInfoConfig"
	EMNET_GetTrainTagInfoConfig_CacheKey 		 = "ResourceConfigAPI:GetTrainTagInfoConfig1"
	EMNET_ResourceConfigAPI_CacheSeconds         = 60 * 60
)

var (
	shareResourceConfigAPILogger dotlog.Logger
	shareResourceConfigAPIRedis  cache.RedisCache
)

func init() {
	protected.RegisterServiceLoader(resourceConfigAPIServiceName, func() {
		shareResourceConfigAPILogger = dotlog.GetLogger(resourceConfigAPIServiceName)
		shareResourceConfigAPIRedis = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	})
}

// GetStrategyLiveRoomConfig 获取所有策略对应的直播间信息
func GetStrategyLiveRoomConfig() (result []*StrategyLiveRoomConfig, err error) {
	err = shareResourceConfigAPIRedis.GetJsonObj(EMNET_GetStrategyLiveRoomConfig_CacheKey, &result)
	if err == nil && result != nil && len(result) > 0 {
		return result, nil
	}
	err = nil
	apiUrl := config.CurrentConfig.ResourceConfigAPI
	if len(apiUrl) == 0 {
		err = errors.New("获取所有策略对应的直播间信息 接口地址配置不正确")
		shareResourceConfigAPILogger.ErrorFormat(err, "获取所有策略对应的直播间信息 接口地址配置不正确 configkey=ResourceConfigAPI")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, config.CurrentConfig.ConfigKey_StrategyLiveRoom)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareResourceConfigAPILogger.DebugFormat("GetStrategyLiveRoomConfig 获取所有策略对应的直播间信息 请求地址为：%s  结果为：%s", apiUrl, body)
	apiresult := ApiResult{}
	err = _json.Unmarshal(body, &apiresult)
	if err != nil {
		return result, err
	}
	if apiresult.RetCode != "0" {
		return result, errors.New(apiresult.RetMsg)
	}

	if len(apiresult.Message.ConfigContent) == 0 {
		err = errors.New("配置内容为空")
		return result, err
	}
	err = _json.Unmarshal(apiresult.Message.ConfigContent, &result)
	if err == nil && result != nil && len(result) > 0 {
		err = shareResourceConfigAPIRedis.Set(EMNET_GetStrategyLiveRoomConfig_CacheKey, _json.GetJsonString(result), EMNET_ResourceConfigAPI_CacheSeconds)
	}
	return result, err
}

// GetStrategyLiveRoomConfig 获取客户端策略对应的专家策略信息
func GetClientExpertStrategyConfig() (result string, err error) {
	err = shareResourceConfigAPIRedis.GetJsonObj(EMNET_GetClientExpertStrategyConfig_CacheKey, &result)
	if err == nil && result != "" && len(result) > 0 {
		return result, nil
	}
	err = nil
	apiUrl := config.CurrentConfig.ResourceConfigAPI
	if len(apiUrl) == 0 {
		err = errors.New("获取客户端策略对应的专家策略信息 接口地址配置不正确")
		shareResourceConfigAPILogger.ErrorFormat(err, "获取客户端策略对应的专家策略信息 接口地址配置不正确 configkey=ResourceConfigAPI")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, config.CurrentConfig.ConfigKey_ClientExpertStrategyRelation)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareResourceConfigAPILogger.DebugFormat("GetClientExpertStrategyConfig 获取客户端策略对应的专家策略信息 请求地址为：%s  结果为：%s", apiUrl, body)
	apiresult := ApiResult{}
	err = _json.Unmarshal(body, &apiresult)
	if err != nil {
		return result, err
	}
	if apiresult.RetCode != "0" {
		return result, errors.New(apiresult.RetMsg)
	}

	if len(apiresult.Message.ConfigContent) == 0 {
		err = errors.New("配置内容为空")
		return result, err
	}
	result = apiresult.Message.ConfigContent
	if err == nil && result != "" && len(result) > 0 {
		err = shareResourceConfigAPIRedis.Set(EMNET_GetClientExpertStrategyConfig_CacheKey, _json.GetJsonString(result), EMNET_ResourceConfigAPI_CacheSeconds)
	}
	return result, err
}

// GetZYPIDConfig 获取智盈产品PID
func GetZYPIDConfig() (result string, err error) {
	err = shareResourceConfigAPIRedis.GetJsonObj(EMNET_GetZYPID_CacheKey, &result)
	if err == nil && result != "" && len(result) > 0 {
		return result, nil
	}
	err = nil
	apiUrl := config.CurrentConfig.ResourceConfigAPI
	if len(apiUrl) == 0 {
		err = errors.New("获取智盈产品PID 接口地址配置不正确")
		shareResourceConfigAPILogger.ErrorFormat(err, "获取智盈产品PID 接口地址配置不正确 configkey=ResourceConfigAPI")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, config.CurrentConfig.ConfigKey_ZYPID)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareResourceConfigAPILogger.DebugFormat("GetZYPIDConfig 获取智盈产品PID 请求地址为：%s  结果为：%s", apiUrl, body)
	apiresult := ApiResult{}
	err = _json.Unmarshal(body, &apiresult)
	if err != nil {
		return result, err
	}
	if apiresult.RetCode != "0" {
		return result, errors.New(apiresult.RetMsg)
	}

	if len(apiresult.Message.ConfigContent) == 0 {
		err = errors.New("配置内容为空")
		return result, err
	}
	result = apiresult.Message.ConfigContent
	if err == nil && result != "" && len(result) > 0 {
		err = shareResourceConfigAPIRedis.Set(EMNET_GetZYPID_CacheKey, _json.GetJsonString(result), EMNET_ResourceConfigAPI_CacheSeconds)
	}
	return result, err
}

// GetTrainClientInfoConfig 获取用户培训显示策略
func GetTrainClientInfoConfig() (result string, err error) {
	err = shareResourceConfigAPIRedis.GetJsonObj(EMNET_GetTrainClientInfoConfig_CacheKey, &result)
	if err == nil && result != "" && len(result) > 0 {
		return result, nil
	}
	err = nil
	apiUrl := config.CurrentConfig.ResourceConfigAPI
	if len(apiUrl) == 0 {
		err = errors.New("获取用户培训显示策略 接口地址配置不正确")
		shareResourceConfigAPILogger.ErrorFormat(err, "获取用户培训显示策略 接口地址配置不正确 configkey=ResourceConfigAPI")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, config.CurrentConfig.ConfigKey_TrainClient)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareResourceConfigAPILogger.DebugFormat("EMNET_GetTrainClientInfoConfig_CacheKey 获取用户培训显示策略 请求地址为：%s  结果为：%s", apiUrl, body)
	apiResult := ApiResult{}
	err = _json.Unmarshal(body, &apiResult)
	if err != nil {
		return result, err
	}
	if apiResult.RetCode != "0" {
		return result, errors.New(apiResult.RetMsg)
	}

	if len(apiResult.Message.ConfigContent) == 0 {
		err = errors.New("配置内容为空")
		return result, err
	}
	result = apiResult.Message.ConfigContent
	if err == nil && result != "" && len(result) > 0 {
		err = shareResourceConfigAPIRedis.Set(EMNET_GetTrainClientInfoConfig_CacheKey, _json.GetJsonString(result), EMNET_ResourceConfigAPI_CacheSeconds)
	}
	return result, err
}

// GetTrainTagInfoConfig 获取用户培训显示标签
func GetTrainTagInfoConfig() (result string, err error) {
	err = shareResourceConfigAPIRedis.GetJsonObj(EMNET_GetTrainTagInfoConfig_CacheKey, &result)
	if err == nil && result != "" && len(result) > 0 {
		return result, nil
	}
	err = nil
	apiUrl := config.CurrentConfig.ResourceConfigAPI
	if len(apiUrl) == 0 {
		err = errors.New("获取用户培训显示标签 接口地址配置不正确")
		shareResourceConfigAPILogger.ErrorFormat(err, "获取用户培训显示标签 接口地址配置不正确 configkey=ResourceConfigAPI")
		return result, err
	}
	apiUrl = fmt.Sprintf(apiUrl, config.CurrentConfig.ConfigKey_TrainTag)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return result, errReturn
	}
	_ = contentType
	_ = intervalTime
	shareResourceConfigAPILogger.DebugFormat("EMNET_GetTrainTagInfoConfig_CacheKey 获取用户培训显示标签 请求地址为：%s  结果为：%s", apiUrl, body)
	apiResult := ApiResult{}
	err = _json.Unmarshal(body, &apiResult)
	if err != nil {
		return result, err
	}
	if apiResult.RetCode != "0" {
		return result, errors.New(apiResult.RetMsg)
	}

	if len(apiResult.Message.ConfigContent) == 0 {
		err = errors.New("配置内容为空")
		return result, err
	}
	result = apiResult.Message.ConfigContent
	if err == nil && result != "" && len(result) > 0 {
		err = shareResourceConfigAPIRedis.Set(EMNET_GetTrainTagInfoConfig_CacheKey, _json.GetJsonString(result), EMNET_ResourceConfigAPI_CacheSeconds)
	}
	return result, err
}

// GetStrategyLiveRoomConfig获取所有策略对应的直播间信息 key为策略Id，value为直播间号
func GetStrategyLiveRoomConfigDict() (result map[string]string, err error) {
	result = make(map[string]string)
	config, err := GetStrategyLiveRoomConfig()
	if err != nil || config == nil {
		return result, err
	}
	for _, item := range config {
		result[item.ClientStrategyId] = item.LiveRoom
	}
	return result, err
}

type StrategyLiveRoomConfig struct {
	//策略Id
	ClientStrategyId string
	//直播间号
	LiveRoom string
}

type ApiResult struct {
	RetCode    string
	RetMsg     string
	IsSucess   bool
	TipMessage string
	Message    ProductVersionConfig
}

type ProductVersionConfig struct {
	ProductVersionConfigId int64
	//配置唯一KEY
	ConfigKey string
	//配置内容
	ConfigContent string
	//配置的格式 1.JSON 2.XML 3.Text
	ConfigFormat int
	//配置名称
	ConfigName string
	//配置描述
	ConfigDesc string
}
