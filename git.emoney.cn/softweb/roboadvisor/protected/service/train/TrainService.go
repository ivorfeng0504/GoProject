package train

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	trainmodel "git.emoney.cn/softweb/roboadvisor/protected/model/train"
	"git.emoney.cn/softweb/roboadvisor/protected/repository/train"
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"github.com/garyburd/redigo/redis"
	"sort"
	"time"
	"git.emoney.cn/softweb/roboadvisor/protected/service/resapi"
	"strconv"
)

type TrainService struct {
	service.BaseService
	trainRepo   *train.TrainRepository
	trainLogger dotlog.Logger
}

var (
	// shareTrainRepo TrainRepository
	shareTrainRepo *train.TrainRepository

	// shareTrainLogger 共享的Logger实例
	shareTrainLogger dotlog.Logger
)

const (
	trainServiceName          = "trainServiceLogger"
	CacheKey_TrainListByPID   = _const.RedisKey_NewsPre + "Train.GetTrainListByPID"
	CacheKey_ServiceAgentName = _const.RedisKey_NewsPre + "Train.GetServiceAgentName"
	CacheKey_TrainListByPIDAndDate = _const.RedisKey_NewsPre + "Train.GetTrainListByPIDAndDate"
	CacheKey_TrainListByPIDAndDateAndTag = _const.RedisKey_NewsPre + "Train.GetTrainListByPIDAndDateAndTag"
	trainBeginTime            = "2018-11-15"
	baseTimeFormat            = "2006-01-02"
	//培训数据缓存时间 10分钟
	TrainRedisTTL = 60*10
)

func NewTrainService() *TrainService {
	trainService := &TrainService{
		trainRepo: shareTrainRepo,
	}
	trainService.RedisCache = _cache.GetRedisCacheProvider(protected.DefaultRedisConfig)
	trainService.trainLogger = shareTrainLogger
	return trainService
}

func (service *TrainService) GetTrainLogger() dotlog.Logger {
	return service.trainLogger
}

func (service *TrainService) RedisDuration(key string) time.Duration {
	service.RedisCache.Set(key, "1", 60*1000)

	t1 := time.Now()
	service.RedisCache.Get(key)
	elapsed := time.Since(t1)

	return elapsed
}

// 获取培训课程存入redis - task任务定时获取
func (service *TrainService) RefreshTrainListToRedis(pid string) error {
	trainlist, err := service.trainRepo.GetTrainListByPID(pid)

	jsonstr, err := _json.Marshal(trainlist)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "RefreshTrainListToRedis 刷新智盈大师培训课程数据 异常 pid=%s", pid, _json.GetJsonString(trainlist))
		return err
	}
	service.RedisCache.HSet(CacheKey_TrainListByPID, pid, jsonstr)

	return err
}

// GetTrainListFromRedis 获取智盈培训课程（近一个月所有课程）
func (service *TrainService) GetTrainListFromRedis(pid string) ([]*trainmodel.NetworkMeetingInfo, error) {
	var trainList []*trainmodel.NetworkMeetingInfo
	jsonStr, err := service.RedisCache.HGet(CacheKey_TrainListByPID, pid)
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetTrainListFromRedis 获取培训课程列表失败 pid=%s", pid)
		return nil, err
	}
	if jsonStr != "" {
		err = _json.Unmarshal(jsonStr, &trainList)
		if err != nil {
			shareTrainLogger.ErrorFormat(err, "GetTrainListFromRedis 获取培训课程反序列化失败 pid=%s jsonStr=%s", pid, jsonStr)
			return nil, err
		}
		return trainList, err
	}
	return nil, nil
}

// GetTrainListByTag 根据标签获取培训课程
func (service *TrainService) GetTrainListByTag(pid string, tag int, currpage int, pageSize int) ([]*trainmodel.NetworkMeetingInfo, int, error) {
	trainList, err := service.GetTrainListFromRedis(pid)
	var allList []*trainmodel.NetworkMeetingInfo
	var retList []*trainmodel.NetworkMeetingInfo
	for _, v := range trainList {
		//取2018-11-15号之后的数据显示
		begintime, _ := time.Parse(baseTimeFormat, trainBeginTime)
		if v.Class_date.After(begintime) {
			//标签=1，读取1和2标签的数据（解盘分析+益晨会）。
			if tag == 1 {
				if v.TrainTag == 1 || v.TrainTag == 2 {
					v.TrainTagName = GetTrainTagName(v.TrainTag)
					allList = append(allList, v)
				}
			} else {
				//根据标签过滤培训课程
				if v.TrainTag == tag {
					v.TrainTagName = GetTrainTagName(v.TrainTag)
					allList = append(allList, v)
				}
			}
		}
	}
	retList, totalCount := service.GetPage(allList, pageSize, currpage)
	return retList, totalCount, err
}

func (service *TrainService) GetPage(allList []*trainmodel.NetworkMeetingInfo, pageSize int, currPage int) ([]*trainmodel.NetworkMeetingInfo, int) {
	var retList []*trainmodel.NetworkMeetingInfo
	//分页
	totalCount := len(allList)

	if currPage <= 0 {
		currPage = 1
	}
	if pageSize <= 0 || totalCount <= 0 {
		return nil, totalCount
	}

	if pageSize > totalCount {
		pageSize = totalCount
	}

	if currPage > (totalCount/pageSize)+1 {
		return nil, totalCount
	}

	beginnum := (currPage - 1) * pageSize
	endnum := pageSize

	if currPage == 1 {
		retList = allList[:endnum]
	} else {
		endnum = pageSize * currPage
		//数量不够，endnum=开始num+剩余数量
		if (totalCount - beginnum) < pageSize {
			endnum = beginnum + (totalCount - beginnum)
		}

		if beginnum >= totalCount {
			return nil, totalCount
		}
		retList = allList[beginnum:endnum]
	}
	return retList, totalCount
}

// GetTrainListByArea 根据地区获取培训课程
func (service *TrainService) GetTrainListByDateAndArea(pid string, date time.Time, area string) ([]*trainmodel.NetworkMeetingInfo, error) {
	trainlist, err := service.GetTrainListFromRedis(pid)

	var retList []*trainmodel.NetworkMeetingInfo
	for _, v := range trainlist {
		//根据地区和日期过滤培训课程
		class_date := v.Class_date.Format("2006-01-02")
		req_date := date.Format("2006-01-02")
		if (v.Ddlarea == area || v.Ddlarea == "全国") && class_date == req_date {
			v.TrainTagName = GetTrainTagName(v.TrainTag)
			retList = append(retList, v)
		}
	}

	return retList, err
}


// GetTrainListByNew 分页获取当日最新（用户培训最近一个月课程+实战培训最近一个月课程）
func (service *TrainService) GetTrainListByNew(pid string, currpage int, pageSize int) ([]*trainmodel.NetworkMeetingInfo, int, error) {
	var allList []*trainmodel.NetworkMeetingInfo
	var retList []*trainmodel.NetworkMeetingInfo
	trainList, err := service.GetTrainListFromRedis(pid)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetTrainListByNew 获取最近一个月培训课程列表失败 pid=%s", pid)
		return nil, 0, err
	}
	for _, v := range trainList {
		//取2018-11-15号之后的数据显示
		begintime, _ := time.Parse(baseTimeFormat, trainBeginTime)
		if v.Class_date.After(begintime) {
			v.TrainTagName = GetTrainTagName(v.TrainTag)
			allList = append(allList, v)
		}
	}

	//获取最近一个月实战培训课程数据
	newsservice := strategyservice.NewColumnInfoService()
	newsList, _, err := newsservice.GetStrategyNewsListByPage_MultiMedia1Month(1, 1000)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetTrainListByNew 获取最近一个月实战培训课程列表失败")
		return nil, 0, err
	}

	//栏目id和策略字典
	dicColumnAndClName, err := agent.GetColumnIdStrategyNameDict()

	for _, v := range newsList {
		trainInfo := new(trainmodel.NetworkMeetingInfo)
		trainInfo.ID = v.ID
		trainInfo.Mtg_name = v.Title
		trainInfo.Mtg_id = 0
		trainInfo.Ddlarea = "全国"
		trainInfo.CoverImg = v.CoverImg
		trainInfo.TrainTagName = "策略实战"

		classdate := v.Live_StartTime
		if classdate.Year() == 1 {
			classdate = v.CreateTime
		}
		trainInfo.Class_date = classdate

		enddate := v.Live_EndTime
		if enddate.Year() == 1 {
			enddate = v.CreateTime
		}
		trainInfo.EndDate = enddate

		//录播有值给录播地址，直播后的录播有值给直播后的录播地址
		if len(v.VideoPlayURL) > 0 {
			trainInfo.Video_url = v.VideoPlayURL
		}
		if len(v.LiveVideoURL) > 0 {
			trainInfo.Video_url = v.LiveVideoURL
		}

		trainInfo.Gensee_URL = v.LiveURL

		if len(v.ShowPosition) > 0 {
			var showPositionInfoList []*agent.ShowPositionInfo

			err := _json.Unmarshal(v.ShowPosition, &showPositionInfoList)

			if err != nil {
				shareTrainLogger.ErrorFormat(err, "GetTrainListByNew 实战培训栏目反序列化失败")
			} else {
				if len(showPositionInfoList) > 0 {
					trainInfo.Txtteachar = dicColumnAndClName[showPositionInfoList[0].ShowPosition]
				}
			}
		}
		allList = append(allList, trainInfo)
	}
	//排序
	sort.Sort(NetworkMeetingInfoData(allList))

	//分页
	retList, totalCount := service.GetPage(allList, pageSize, currpage)

	return retList, totalCount, err
}




// GetTrainListByDateAndPid 根据日期和版本获取培训课程 (最新迭代需求)
func (service *TrainService) GetTrainListByDateAndPid(pid string, date *time.Time) ([]*trainmodel.NetworkMeetingInfo, error) {
	dateFormat := date.Format("2006-01-02")
	redisKey := fmt.Sprintf(CacheKey_TrainListByPIDAndDate+":%s_%s", pid, dateFormat)
	var retList []*trainmodel.NetworkMeetingInfo
	err := service.RedisCache.GetJsonObj(redisKey, &retList)
	if err == redis.ErrNil {
		err = nil
	}
	if retList == nil || len(retList) == 0 {
		trainList, err := service.trainRepo.GetTrainListByDateAndPid(pid, date)
		if err != nil {
			shareTrainLogger.ErrorFormat(err, "GetTrainListByDateAndPid 根据日期和版本获取培训课程失败 pid=%s", pid)
			return nil, err
		}

		if len(trainList) > 0 {
			service.RedisCache.Set(redisKey, _json.GetJsonString(trainList), TrainRedisTTL)
			retList = trainList
		}
	}
	return retList, nil
}

// GetTrainListByNow 当日最新 （用户培训当天数据+实战培训数据 最新迭代需求）
func (service *TrainService) GetTrainListByNow(pid string, date *time.Time) ([]*trainmodel.NetworkMeetingInfo, error) {
	var allList []*trainmodel.NetworkMeetingInfo

	//根据日期和产品id获取培训课表
	trainList, err := service.GetTrainListByDateAndPid(pid, date)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetTrainListByNow 根据日期获取培训课程列表失败 pid=%s", pid)
		return nil, err
	}
	//查询培训标签信息
	ret, err := resapi.GetTrainTagInfoConfig()
	var tagList []*trainmodel.TrainTagInfo
	err = _json.Unmarshal(ret, &tagList)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetTrainListByDateAndPidAndTag 获取培训标签反序列化失败 ret=%s", ret)
	}
	//匹配标签名称后显示
	for _, v := range trainList {
		for _, tag := range tagList {
			if strconv.Itoa(v.TrainTag) == tag.TrainTagID {
				v.TrainTagName = tag.TrainTagName
			}
		}
		allList = append(allList, v)
	}

	// 合并实战培训数据
	newsService := strategyservice.NewColumnInfoService()
	newsList, err := newsService.GetStrategyNewsListByPage_MultiMedia_ByDate_All(date)

	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetStrategyNewsListByPage_MultiMedia_ByDate_All 根据日期获取实战培训课程列表失败")
		return nil, err
	}
	//栏目id和策略字典
	srv := strategyservice.NewClientStrategyInfoService()
	result, err := srv.GetClientStrategyInfoList(true)
	if err != nil {
		shareTrainLogger.ErrorFormat(err, "GetClientStrategyInfoList 用户培训获取栏目id和策略字典失败")
	}

	dicColumnAndClName := make(map[string]string)
	for _, item := range result {
		dicColumnAndClName[strconv.Itoa(item.ColumnInfoId)] = item.ClientStrategyName
	}

	for _, v := range newsList {
		trainInfo := new(trainmodel.NetworkMeetingInfo)
		trainInfo.ID = v.ID
		trainInfo.Mtg_name = v.Title
		trainInfo.Mtg_id = 0
		trainInfo.Ddlarea = "全国"
		trainInfo.CoverImg = v.CoverImg
		trainInfo.TrainTagName = "策略实战"

		classdate := v.Live_StartTime
		if classdate.Year() == 1 {
			classdate = v.CreateTime
		}
		trainInfo.Class_date = classdate

		enddate := v.Live_EndTime
		if enddate.Year() == 1 {
			enddate = v.CreateTime
		}
		trainInfo.EndDate = enddate

		//录播有值给录播地址，直播后的录播有值给直播后的录播地址
		if len(v.VideoPlayURL) > 0 {
			trainInfo.Video_url = v.VideoPlayURL
		}
		if len(v.LiveVideoURL) > 0 {
			trainInfo.Video_url = v.LiveVideoURL
		}

		trainInfo.Gensee_URL = v.LiveURL

		if len(v.ShowPosition) > 0 {
			var showPositionInfoList []*agent.ShowPositionInfo
			err := _json.Unmarshal(v.ShowPosition, &showPositionInfoList)
			if err != nil {
				shareTrainLogger.ErrorFormat(err, "GetTrainListByNew 实战培训栏目反序列化失败")
			} else {
				if len(showPositionInfoList) > 0 && len(dicColumnAndClName) > 0 {
					trainInfo.Txtteachar = dicColumnAndClName[showPositionInfoList[0].ShowPosition]
				}
			}
		}
		//过滤中继回踩培训数据
		if trainInfo.Txtteachar != _const.TrainFilterCL_ZJHC {
			allList = append(allList, trainInfo)
		}
	}
	//排序
	sort.Sort(NetworkMeetingInfoData(allList))
	return allList, err
}

// GetTrainListByDateAndPid 根据日期和版本获取培训课程 (最新迭代需求)
func (service *TrainService) GetTrainListByDateAndPidAndTag(pid string, date *time.Time,trainTag int,trainType string,pageSize int, currPage int) ([]*trainmodel.NetworkMeetingInfo,int, error) {
	var redisKey = ""
	var retList []*trainmodel.NetworkMeetingInfo
	var filterList []*trainmodel.NetworkMeetingInfo
	var trainList []*trainmodel.NetworkMeetingInfo
	if date != nil {
		redisKey = fmt.Sprintf(CacheKey_TrainListByPIDAndDateAndTag+":%s_%s_%d_%s", pid, date.String(), trainTag, trainType)
	}else{

		redisKey = fmt.Sprintf(CacheKey_TrainListByPIDAndDateAndTag+":%s_%d_%s", pid, trainTag, trainType)
	}
	err := service.RedisCache.GetJsonObj(redisKey, &retList)
	if err == redis.ErrNil {
		err = nil
	}
	if retList == nil || len(retList) == 0 {
		//不根据日期则查询最近一个月数据
		if date != nil {
			trainList, err = service.trainRepo.GetTrainListByDateAndPid(pid, date)
		} else {
			trainList, err = service.trainRepo.GetTrainListByPID(pid)
		}
		if err != nil {
			shareTrainLogger.ErrorFormat(err, "GetTrainListByDateAndPidAndTag 根据日期版本和标签获取培训课程失败 pid=%s", pid)
			return nil, 0, err
		}

		//根据类型查询培训标签
		ret, err := resapi.GetTrainTagInfoConfig()
		var tagList []*trainmodel.TrainTagInfo
		var retTagList []*trainmodel.TrainTagInfo
		err = _json.Unmarshal(ret, &tagList)
		if err != nil {
			shareTrainLogger.ErrorFormat(err, "GetTrainListByDateAndPidAndTag 获取培训标签反序列化失败 ret=%s", ret)
		}
		dicTagInfo := make(map[string]string)
		for _, item := range tagList {
			dicTagInfo[item.TrainTagID] = item.TrainTagName
		}
		for _, v := range tagList {
			if v.Type == trainType {
				retTagList = append(retTagList, v)
			}
		}

		//过滤标签后返回
		for _, v := range trainList {
			//标签=0 获取分类下的所有标签综合
			if trainTag == 0 {
				for _, tag := range retTagList {
					if strconv.Itoa(v.TrainTag) == tag.TrainTagID {
						v.TrainTagName = tag.TrainTagName
						filterList = append(filterList, v)
					}
				}
			} else {
				if v.TrainTag == trainTag {
					v.TrainTagName = dicTagInfo[strconv.Itoa(v.TrainTag)]
					filterList = append(filterList, v)
				}
			}
		}

		//存入redis
		if len(filterList) > 0 {
			service.RedisCache.Set(redisKey, _json.GetJsonString(filterList), TrainRedisTTL)
			retList = filterList
		}
	}
	//分页
	retObj, totalCount := service.GetPage(retList, pageSize, currPage)
	return retObj,totalCount, nil
}


// GetServiceAgentName 获取用户所在地
func (service *TrainService) GetServiceAgentName(agentid string, username string) (string, error) {
	var area = ""
	err := service.RedisCache.GetJsonObj(CacheKey_ServiceAgentName, area)
	if area != "" {
		return area, nil
	}

	apiUrl := config.CurrentConfig.ServiceAgentNameApi
	apiUrl = fmt.Sprintf(apiUrl, agentid, username)
	body, contentType, intervalTime, errReturn := _http.HttpGet(apiUrl)
	if errReturn != nil {
		return "", errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		return "", err
	}
	if apiGatewayResp.RetCode != 0 {
		return "", errors.New(apiGatewayResp.RetMsg)
	}
	retarea := apiGatewayResp.Message
	if retarea != "" {
		service.RedisCache.SetJsonObj(CacheKey_ServiceAgentName, retarea)
	}

	return area, err
}

func GetTrainTagName(tagid int) string {
	tagname := ""
	switch tagid {
	case 1:
		tagname = "解盘分析"
		break
	case 2:
		tagname = "益晨会"
		break
	case 3:
		tagname = "价值发现"
		break
	case 4:
		tagname = "基础理论"
		break
	case 5:
		tagname = "实战运用"
		break
	}
	return tagname
}

func init() {
	protected.RegisterServiceLoader(trainServiceName, trainServiceLoader)
}

func trainServiceLoader() {
	shareTrainRepo = train.NewTrainRepository(protected.DefaultConfig)
	shareTrainLogger = dotlog.GetLogger(trainServiceName)
}

// 培训-开课时间排序
type NetworkMeetingInfoData []*trainmodel.NetworkMeetingInfo

func (s NetworkMeetingInfoData) Len() int {
	return len(s)
}
func (s NetworkMeetingInfoData) Less(i, j int) bool {
	return s[i].Class_date.After(s[j].Class_date)
}
func (s NetworkMeetingInfoData) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
