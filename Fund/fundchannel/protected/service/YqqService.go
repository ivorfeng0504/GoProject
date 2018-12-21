package service

import (
	"time"
	"fmt"
	"github.com/devfeel/dotlog"
	"emoney.cn/fundchannel/util/cache"
	"emoney.cn/fundchannel/protected"
	"emoney.cn/fundchannel/util/http"
	"emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/contract"
	"emoney.cn/fundchannel/util/json"
	"emoney.cn/fundchannel/protected/model/yqq"
)

type YqqService struct {
	BaseService
}

var (
	fundChannelYqqLogger dotlog.Logger
)

const (
	YqqServiceName    = FundChannelKey + "YqqService"
	YqqNewLiveKey     = FundChannelKey + "Yqq.NewLive"
	YqqAllLiveKey     = FundChannelKey + "Yqq.AllLive_Date[%s]"
	YqqNewQuestionKey = FundChannelKey + "Yqq.NewQuestion"
	YqqLiveRoomKey    = FundChannelKey + "Yqq.LiveRoom"
	YqqAllQuestionKey = FundChannelKey + "Yqq.AllQuestion"
)

func NewYqqService() *YqqService {
	service := &YqqService{}
	service.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	return service
}

func init() {
	protected.RegisterServiceLoader(YqqServiceName, yqqServiceLoader)
}

func yqqServiceLoader() {
	fundChannelYqqLogger = dotlog.GetLogger("YqqService")
}

func (service *YqqService) GetYqqLiveLatestContent() (arr []string, err error) {
	var flg string
	flg, _, _, err = _http.HttpGet(config.CurrentConfig.YqqApiNewLiveInfoUrl)
	if err != nil {
		fundChannelYqqLogger.Error(err, " GetYqqLiveLatestContent 获取益圈圈直播最新内容 异常")
	}
	fundChannelYqqLogger.Debug(" GetYqqLiveLatestContent 获取益圈圈直播最新内容：" + flg)

	live, _ := service.RedisCache.GetString(YqqNewLiveKey)
	arr = append(arr, live)

	question, _ := service.RedisCache.GetString(YqqNewQuestionKey)
	arr = append(arr, question)

	room, _ := service.RedisCache.GetString(YqqLiveRoomKey)
	arr = append(arr, room)

	return arr, err
}

func (service *YqqService) GetLiveRoomInfo() (roomInfo string, err error) {
	roomInfo, err = service.RedisCache.GetString(YqqLiveRoomKey)
	return roomInfo, err
}

func (service *YqqService) GetYqqLiveAllContent(date string, msgId string) (response string, err error) {
	var flg string
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(YqqAllLiveKey, date)

	if today != date {
		response, err = service.RedisCache.GetString(key)
		if len(response) > 0 {
			return response, err
		}
	}

	url := fmt.Sprintf(config.CurrentConfig.YqqApiAllLiveInfoUrl+"&date=%s&messageId=%s", date, msgId)
	flg, _, _, err = _http.HttpGet(url)
	if err != nil {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqLiveAllContent 获取益圈圈直播全部内容 异常:url=%s", url)
	}
	fundChannelYqqLogger.Debug(" GetYqqLiveAllContent 获取益圈圈直播全部内容：" + flg)

	response, err = service.RedisCache.GetString(key)

	return response, err
}

func (service *YqqService) GetYqqAllQuestion(msgId string) (response string, err error) {
	var flg string
	url := fmt.Sprintf(config.CurrentConfig.YqqApiAllQuestionUrl+"&messageId=%s", msgId)
	flg, _, _, err = _http.HttpGet(url)
	if err != nil {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqLiveAllContent 获取益圈圈直播全部内容 异常:url=%s", url)
	}
	fundChannelYqqLogger.Debug(" GetYqqAllQuestion 获取益圈圈全部问答内容：" + flg)

	response, err = service.RedisCache.GetString(YqqAllQuestionKey)

	return response, err
}

func (service *YqqService) GetYqqMyQuestion(ssoQuery string) (response contract.ResponseInfo, err error) {
	url := config.CurrentConfig.YqqApiMyQuestionUrl + "&page=0&" + ssoQuery
	var json string
	json, _, _, err = _http.HttpGet(url)
	if err != nil {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqMyQuestion 获取益圈圈我的问答内容 异常:url=%s", url)
		response.RetMsg = "获取益圈圈我的问答内容 异常"
		response.RetCode = -2
		return response, err
	}
	if len(json) == 0 {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqMyQuestion 获取益圈圈我的问答内容 返回消息为空:url=%s", url)
		response.RetMsg = "获取益圈圈我的问答内容 返回消息为空"
		response.RetCode = -1
		return response, err
	}
	fundChannelYqqLogger.DebugFormat(" GetYqqMyQuestion 获取益圈圈我的问答内容:url=%s；json=%s", url, json)

	yqqMyQuestionResult := new(yqq.YqqMyQuestionResult)
	fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD1")

	err = _json.Unmarshal(json, yqqMyQuestionResult)
	if err != nil {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqMyQuestion 获取益圈圈我的问答内容 json序列化失败:json=%s", json)
		response.RetMsg = "获取益圈圈我的问答内容 json序列化失败"
		response.RetCode = -1
		return response, err
	}
	fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD2")
	fundChannelYqqLogger.Debug(_json.GetJsonString(yqqMyQuestionResult))
	fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD2.1")

	//response.RetCode, err = strconv.Atoi(yqqMyQuestionResult.RetCode)
	if err != nil {
		fundChannelYqqLogger.ErrorFormat(err, " GetYqqMyQuestion 获取益圈圈我的问答内容 RetCode序列化失败:json=%s", json)
		response.RetMsg = "获取益圈圈我的问答内容 RetCode序列化失败"
		response.RetCode = -1
		return response, err
	}
	fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD3")

	if yqqMyQuestionResult !=nil {
		response.RetCode = 99
		response.RetMsg = yqqMyQuestionResult.RetMsg
		response.Message = yqqMyQuestionResult.Data
		fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD4")
	}

	fundChannelYqqLogger.Debug(" GetYqqMyQuestion MDGD5")

	return response, err
}
