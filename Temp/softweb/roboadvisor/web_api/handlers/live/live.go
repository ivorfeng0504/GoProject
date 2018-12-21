package live

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/global"
	vmMapper "git.emoney.cn/softweb/roboadvisor/protected/mapper"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service/dataapi"
	liveservice "git.emoney.cn/softweb/roboadvisor/protected/service/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/protected/service/tokenapi"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel"
	"git.emoney.cn/softweb/roboadvisor/util/array"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/mapper"
	"html/template"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// IsTradeTime 当前是否为交易时间
// Deprecated:迁移到web_api_live项目
func IsTradeTime(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}

	isTradeTime, err := dataapi.IsTradeTimeNow()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = isTradeTime
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// AddQuestion 用户提问
// Deprecated:迁移到web_api_live项目
func AddQuestion(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.AddQuestionRequest{}
	err = mapper.MapperMap(request.RequestData.(map[string]interface{}), &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if len(requestData.AskContent) == 0 {
		response.RetCode = -3
		response.RetMsg = "提问内容不能为空"
		return ctx.WriteJson(response)
	}
	if requestData.RoomId <= 0 {
		response.RetCode = -4
		response.RetMsg = "房间号不正确"
		return ctx.WriteJson(response)
	}
	if len(requestData.Source) < 5 {
		response.RetCode = -5
		response.RetMsg = "用户来源不正确"
		return ctx.WriteJson(response)
	}

	//对提问内容进行编码
	requestData.AskContent = template.HTMLEscapeString(requestData.AskContent)

	userService := liveservice.NewUserService()
	questionService := liveservice.NewLiveQuestionAnswerService()
	var userInfo *livemodel.User
	//1、优先使用UserId进行用户查询
	//2、如果没有传递UserId，则通过手机号查询UserId
	//3、如果手机号查询不存在则使用UID进行查询
	if requestData.UserId <= 0 {
		//如果没有传递昵称 使用掩码手机号作为昵称
		if len(requestData.NickName) == 0 {
			requestData.NickName = requestData.MaskMobile
		}
		//如果没有传递UserId则使用Mobile进行用户检测，此时必须传递NickName
		if len(requestData.NickName) == 0 {
			response.RetCode = -6
			response.RetMsg = "昵称不能为空"
			return ctx.WriteJson(response)
		}
		if len(requestData.Mobile) == 0 || requestData.UID <= 0 {
			response.RetCode = -7
			response.RetMsg = "手机号或UID不正确"
			return ctx.WriteJson(response)
		}
		//首先使用手机号校验用户
		//十六进制编码成base64
		requestData.Mobile = mobile.Hex2Base64(requestData.Mobile)
		userInfo, err = userService.GetUserByAccount(requestData.Mobile)
		if err != nil {
			response.RetCode = -8
			response.RetMsg = "查询用户信息失败"
			return ctx.WriteJson(response)
		}

		//用户不存在则使用UID进行二次校验
		if userInfo == nil {
			userInfo, err = userService.GetUserByUID(requestData.UID)
			if err != nil {
				response.RetCode = -9
				response.RetMsg = "查询用户信息失败"
				return ctx.WriteJson(response)
			}
		}

		//如果用户仍然不存在 则自动注册用户记录
		if userInfo == nil {
			requestData.UserId, err = userService.AddUser(requestData.Mobile, requestData.NickName, requestData.UID, requestData.Source)
			if err != nil {
				response.RetCode = -10
				response.RetMsg = "自动注册用户失败," + err.Error()
				return ctx.WriteJson(response)
			}
		} else {
			requestData.UserId = userInfo.UserId
		}
	}

	userInfo, err = userService.GetUserById(requestData.UserId)

	if err != nil {
		response.RetCode = -11
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if userInfo == nil {
		response.RetCode = -12
		response.RetMsg = "该用户不存在"
		return ctx.WriteJson(response)
	}
	if requestData.CheckPID {
		canIAsk, tip, err := questionService.CheckPID(requestData.PID, requestData.UserId, requestData.RoomId)
		if err != nil {
			response.RetCode = -13
			response.RetMsg = "用户权限校验异常"
			return ctx.WriteJson(response)
		}
		if canIAsk == false {
			response.RetCode = -14
			response.RetMsg = tip
			return ctx.WriteJson(response)
		}
	}
	_, err = questionService.AddLiveQuestionAnswer(userInfo.UserId, userInfo.NickName, requestData.AskContent, requestData.RoomId, requestData.Source)
	if err != nil {
		response.RetCode = -15
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	} else {
		//递增用户提问次数
		questionService.IncreAskPermission(requestData.PID, requestData.UserId, requestData.RoomId)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetLiveContent 根据直播间Id（RoomId）与时间(Date)获取直播内容
// RoomId 直播间号
// Date 直播时间 如果Date为空则默认取当前时间
// IsIncre 是否获取递增的内容 如果为1则为递增
// IsQueryTop 是否查询指定消息 如果为1则为查询置顶
// Deprecated:迁移到web_api_live项目
func GetLiveContent(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.LiveContentRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if requestData.RoomIdList == nil || len(requestData.RoomIdList) == 0 {
		response.RetCode = -3
		response.RetMsg = "房间号不正确"
		return ctx.WriteJson(response)
	}

	date, err := _time.ParseTime(requestData.Date)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = "时间格式不正确"
		return ctx.WriteJson(response)
	}
	topicService := liveservice.NewLiveTopicService()
	contentService := liveservice.NewLiveContentService()
	var liveContentListResult viewmodel.LiveContentList

	//显示上一个交易日的主题 并且当前是交易日时需要提供倒计时
	needCountdown := false
	//如果指定日期中有一个直播间有开播记录 则忽略DisPlayLastTopic相关逻辑
	ignoreLastTopic := topicService.HaveTopicOpen(date, requestData.RoomIdList...)
	now := time.Now()
	preStartTime, err := time.Parse("2006-1-2 15:04 MST", now.Format("2006-1-2")+" 08:20 CST")

	for _, roomId := range requestData.RoomIdList {
		topic, err := topicService.GetTopic(roomId, date)
		if err != nil {
			response.RetCode = -5
			response.RetMsg = err.Error()
			continue
		}
		if topic == nil {
			//是否为交易日
			isTradeTime, err := dataapi.IsTradeDayToDay()
			if err != nil {
				global.InnerLogger.ErrorFormat(err, "查询交易日信息异常")
			}
			//如果查看当天直播 并且在8:20之前则显示前一天的主题内容(requestData.DisPlayLastTopic为true的前提下)
			if ignoreLastTopic == false && requestData.DisPlayLastTopic && date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day() && (now.Before(preStartTime) || isTradeTime == false) {
				//如果是交易日才需要为前端提供倒计时
				if isTradeTime {
					needCountdown = true
				}
				topic, err = topicService.GetNewestTopic(roomId)
				if err != nil {
					response.RetCode = -6
					response.RetMsg = err.Error()
					continue
				}
				if topic == nil {
					response.RetCode = -7
					response.RetMsg = "当前直播主题不存在！"
					continue
				}
			} else {
				response.RetCode = -8
				response.RetMsg = "当前直播主题不存在！"
				continue
			}
		}
		var contentList []*livemodel.LiveContent
		var topContent *livemodel.LiveContent

		//是否查询递增的内容
		if requestData.IsIncre == 1 {
			contentList, err = contentService.GetContentIncre(topic.Id, date)
		} else if requestData.IsQueryTop == 1 {
			topContent, err = contentService.GetTopContent(topic.Id)
			if topContent != nil {
				contentList = append(contentList, topContent)
			}
		} else {
			contentList, err = contentService.GetContentList(topic.Id)
		}

		if err != nil {
			response.RetCode = -9
			response.RetMsg = "当前直播内容不存在！"
			continue
		}

		liveContentListResult = append(liveContentListResult, vmMapper.MapperLiveContent(contentList)...)
	}

	sort.Sort(liveContentListResult)
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = liveContentListResult
	if needCountdown {
		response.RetMsg = strconv.FormatFloat(preStartTime.Sub(now).Seconds()*1000, 'f', -1, 64)
	}
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetLiveQuestionAnswerList 获取指定直播间当天的问答列表，如果IsSelf=1，则筛选当前用户的提问列表
// Deprecated:迁移到web_api_live项目
func GetLiveQuestionAnswerList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.QuestionAnswerRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if requestData.RoomIdList == nil || len(requestData.RoomIdList) == 0 {
		response.RetCode = -3
		response.RetMsg = "房间号不正确"
		return ctx.WriteJson(response)
	}

	date, err := _time.ParseTime(requestData.Date)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = "时间格式不正确"
		return ctx.WriteJson(response)
	}

	//如果查询自己的问答并且没有传递UserId，则通过手机号查询UserId
	if requestData.IsSelf == 1 && requestData.UserId <= 0 {
		userService := liveservice.NewUserService()
		user, err := userService.GetUserByAccount(mobile.Hex2Base64(requestData.Mobile))
		if err != nil {
			response.RetCode = -5
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
		if user == nil {
			response.RetCode = -6
			response.RetMsg = "未获取到用户信息"
			return ctx.WriteJson(response)
		}
		requestData.UserId = user.UserId
	}

	roomService := liveservice.NewLiveRoomService()
	questionService := liveservice.NewLiveQuestionAnswerService()
	var questionListResult viewmodel.LiveQuestionAnswerList

	for _, roomId := range requestData.RoomIdList {
		roomInfo, err := roomService.GetLiveRoom(roomId)
		if err != nil {
			response.RetCode = -7
			response.RetMsg = err.Error()
			continue
		}
		if roomInfo == nil {
			response.RetCode = -8
			response.RetMsg = "直播间不存在"
			continue
		}
		var questionList []*livemodel.LiveQuestionAnswer
		if requestData.IsSelf == 1 {
			questionList, err = questionService.GetLiveQuestionAnswerList(roomId, date, requestData.UserId)
		} else {
			questionList, err = questionService.GetLiveQuestionAnswerListByAnswered(roomId, date)
		}

		if err != nil {
			response.RetCode = -9
			response.RetMsg = "当前问答内容不存在！"
			continue
		}
		questionListResult = append(questionListResult, vmMapper.MapperLiveQuestionAnswer(questionList, roomInfo)...)
	}

	sort.Sort(questionListResult)
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = questionListResult
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

//HasNewMessage Long-Polling时间窗口数据拉取
// Deprecated:迁移到web_api_live项目
func HasNewMessage(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.HasNewMessageRequest{}
	err = mapper.MapperMap(request.RequestData.(map[string]interface{}), &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	date, err := _time.ParseTime(requestData.Date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "时间格式不正确"
		return ctx.WriteJson(response)
	}
	liveDate := &viewmodel.LiveDataWithTime{}
	topicService := liveservice.NewLiveTopicService()
	topic, err := topicService.GetTopic(requestData.RoomId, date)
	if err != nil || topic == nil {
		liveDate.ContentNum = 0
	} else {
		contentService := liveservice.NewLiveContentService()
		liveDate.ContentNum = contentService.GetContentListByLidIncreNum(topic.Id, date)
	}
	liveDate.RepliedNum = 0
	liveDate.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = liveDate
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// AddLiveUserInRoom 添加用户直播间权限
// Deprecated:迁移到web_api_live项目
func AddLiveUserInRoom(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	jsonStr := _json.GetJsonString(request.RequestData)
	requestData := agent.AddLiveUserInRoomRequest{}
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.Rooms == nil || len(requestData.Rooms) == 0 {
		response.RetCode = -3
		response.RetMsg = "请传递需要添加的直播间信息"
		return ctx.WriteJson(response)
	}
	//十六进制编码成base64
	requestData.Mobile = mobile.Hex2Base64(requestData.Mobile)

	liveUserInRoomService := liveservice.NewLiveUserInRoomService()
	liveRoomPidMapService := liveservice.NewLiveRoomPidMapService()
	hasError := false
	errInfos := ""
	var successRoomInfo []agent.UserRoomInfo
	for _, room := range requestData.Rooms {
		room.Source = requestData.Source
		//如果roomid小于等于0 则查询用pid查找映射关系
		if room.RoomId <= 0 {
			room.RoomId = liveRoomPidMapService.GetRoomId(room.PID)
			if room.RoomId <= 0 {
				hasError = true
				errInfo := "注册直播间权限异常,未找到PID对应的RoomId，mobile=" + requestData.Mobile + " roomId=" + strconv.Itoa(room.RoomId) + " pid=" + strconv.Itoa(room.PID) + " expireDay=" + strconv.Itoa(room.ExpireDay) + " 订单号=" + room.OrderId
				errInfos += errInfo + "|"
				global.InnerLogger.Error(err, errInfo)
				continue
			}
		}
		id, err := liveUserInRoomService.AddLiveUserInRoom(requestData.Mobile, room.RoomId, room.ExpireDay, room.OrderId, room.Source)
		if err != nil {
			hasError = true
			errInfo := "注册直播间权限异常,mobile=" + requestData.Mobile + " roomId=" + strconv.Itoa(room.RoomId) + " pid=" + strconv.Itoa(room.PID) + " expireDay=" + strconv.Itoa(room.ExpireDay) + " 订单号=" + room.OrderId
			errInfos += errInfo + "|"
			global.InnerLogger.Error(err, errInfo)
			continue
		}
		//操作成功
		room.LiveOrderId = agent.GetLiveOrderId(int(id))
		successRoomInfo = append(successRoomInfo, room)
	}
	if hasError {
		response.RetCode = -999
		response.RetMsg = errInfos
		response.Message = successRoomInfo
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = successRoomInfo
	return ctx.WriteJson(response)
}

// RemoveLiveUserInRoom 删除用户直播间权限
// Deprecated:迁移到web_api_live项目
func RemoveLiveUserInRoom(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	liveOrderId, success := request.RequestData.(string)
	if success == false {
		response.RetCode = -2
		response.RetMsg = "无效的订单号"
		return ctx.WriteJson(response)
	}

	if len(liveOrderId) < 10 || liveOrderId[:4] != "live" {
		response.RetCode = -3
		response.RetMsg = "无效的订单号"
		return ctx.WriteJson(response)
	}
	deleteOrderIdStr := liveOrderId[4:]
	deleteOrderId, err := strconv.Atoi(deleteOrderIdStr)
	if err != nil || deleteOrderId <= 0 {
		response.RetCode = -4
		response.RetMsg = "无效的订单号"
		return ctx.WriteJson(response)
	}

	liveUserInRoomService := liveservice.NewLiveUserInRoomService()
	err = liveUserInRoomService.DeleteById(deleteOrderId)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetUserRoomList 获取当前用户拥有的直播间权限
// Deprecated:迁移到web_api_live项目
func GetUserRoomList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	mobile := request.RequestData.(string)
	liveUserInRoomService := liveservice.NewLiveUserInRoomService()
	roomList, err := liveUserInRoomService.GetRoomList(mobile)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	defaultRoomList, err := liveUserInRoomService.GetDefaultRoomList()
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	roomList = append(roomList, defaultRoomList...)
	//去重
	roomList = _array.Distinct(roomList)
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = roomList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// CheckToken websocket授权回调接口
// Deprecated:迁移到web_api_live项目
func CheckToken(ctx dotweb.Context) error {
	roomId := ctx.QueryString("groupid")
	roomIds := ctx.QueryString("groupids")
	roomIds, _ = url.QueryUnescape(roomIds)

	global.InnerLogger.Debug(fmt.Sprintf("call CheckToken url=【%s】 groupid=【%s】 groupids=【%s】", ctx.Request().Url(), roomId, roomIds))

	appId := ctx.QueryString("appid")
	userId := ctx.QueryString("userid")
	token := ctx.QueryString("token")
	response := contract.AuthResponse{}
	handlerResp, err := tokenapi.QueryToken(appId, token)

	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else if handlerResp.RetCode != 0 {
		response.RetCode = handlerResp.RetCode
		response.RetMsg = handlerResp.RetMsg
	} else if appId != handlerResp.Message.AppID {
		response.RetCode = -2
		response.RetMsg = "AppId不一致"
	} else {
		//检查传递的直播间是否真的有权限
		//是否有没有权限的直播间
		hasErrRoom := false
		var rooms []int
		err = _json.Unmarshal(handlerResp.Message.TokenBody, &rooms)
		if err != nil || len(rooms) == 0 {
			response.RetCode = -3
			response.RetMsg = "没有直播间权限"
			return ctx.WriteJson(response)
		}
		if roomId != "" {
			if isInRooms(roomId, rooms) == false {
				hasErrRoom = true
			}
		}
		socketListenGroups := contract.SocketListenGroups{}
		err = _json.Unmarshal(roomIds, &socketListenGroups)
		if err != nil {
			response.RetCode = -4
			response.RetMsg = "没有直播间权限"
			return ctx.WriteJson(response)
		}

		if roomId == "" && len(socketListenGroups.IDs) == 0 {
			response.RetCode = -5
			response.RetMsg = "没有直播间权限"
			return ctx.WriteJson(response)
		}
		for _, roomStr := range socketListenGroups.IDs {
			if isInRooms(roomStr, rooms) == false {
				hasErrRoom = true
			}
		}

		if hasErrRoom {
			response.RetCode = -6
			response.RetMsg = "直播间权限有误，包含未授权的直播间"
			return ctx.WriteJson(response)
		}

		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.AppID = appId
		response.UserID = userId
	}
	return ctx.WriteJson(response)
}

//HasNewMessage Long-Polling时间窗口数据拉取
//参数：
//appid：应用编号，统一申请
//groupid：用户组编号，应用自定义
//userid：用户编号，应用需保证userid在同Appid下的唯一性
//querykey：透传key，会透传给应用messageapi，一般用于决定是否有需要马上返回的数据
//返回：
//为空，则表示无需马上返回的数据，请求会继续等待；
//不为空，请求马上返回，向客户端推送业务消息
// Deprecated:迁移到web_api_live项目
func GetHasNewMessage(ctx dotweb.Context) error {
	//global.InnerLogger.Debug("call GetHasNewMessage " + ctx.Request().Url())
	queryKey := ctx.QueryString("querykey")
	roomId := ctx.QueryInt("groupid")
	if strings.Contains(queryKey, "CST") == false {
		queryKey += " CST"
	}
	date, err := _time.ParseTime(queryKey)
	if err != nil {
		date = time.Now()
	}

	liveDate := &viewmodel.LiveDataWithTime{}
	topicService := liveservice.NewLiveTopicService()
	topic, err := topicService.GetTopic(roomId, date)
	if err != nil || topic == nil {
		liveDate.ContentNum = 0
	} else {
		contentService := liveservice.NewLiveContentService()
		liveDate.ContentNum = contentService.GetContentListByLidIncreNum(topic.Id, date)
	}
	liveDate.RepliedNum = 0
	liveDate.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")

	if liveDate.ContentNum == 0 && liveDate.RepliedNum == 0 {
		return ctx.WriteString("")
	}
	return ctx.WriteJson(liveDate)
}

// Deprecated:迁移到web_api_live项目
func isInRooms(roomStr string, rooms []int) bool {
	room, err := strconv.Atoi(roomStr)
	if err != nil {
		return false
	}
	for _, item := range rooms {
		if item == room {
			return true
		}
	}
	return false
}
