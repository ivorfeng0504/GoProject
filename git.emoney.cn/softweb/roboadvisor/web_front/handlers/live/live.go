package live

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/tokenapi"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/util/time"
	"github.com/devfeel/dotweb"
	"strconv"
	"strings"
	"time"
)

// Index 图文直播首页
func Index(ctx dotweb.Context) error {
	submitTongJi(ctx, "pv", "visit")
	submitTongJi(ctx, "uv", "visitor")
	return contract.RenderHtml(ctx, "index.html")
}

// GetLiveContent 根据直播间Id（RoomId）与时间(Date)获取直播内容
// RoomId 直播间号
// Date 直播时间 如果Date为空则默认取当前时间
// IsIncre 是否获取递增的内容 如果为1则为递增
// IsQueryTop 是否查询指定消息 如果为1则为查询置顶
func GetLiveContent(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	isIncre := ctx.PostFormValue("IsIncre")
	isQueryTop := ctx.PostFormValue("IsQueryTop")
	dateStr := ctx.PostFormValue("Date")
	date, err := _time.ParseTime(dateStr)
	if err != nil {
		date = time.Now()
	}

	user := contract.UserInfo(ctx)
	if user == nil || user.RoomList == nil || len(user.RoomList) == 0 {
		response.RetCode = -1
		response.RetMsg = "您没有直播间查看权限"
		return ctx.WriteJson(response)
	}

	requestData := agent.LiveContentRequest{
		RoomIdList:       user.RoomList,
		Date:             date.Format(time.RFC3339),
		DisPlayLastTopic: true,
	}
	if isIncre == "1" {
		requestData.IsIncre = 1
	}
	if isQueryTop == "1" {
		requestData.IsQueryTop = 1
	}
	contentList, retMsg, err := agent.GetLiveContent(requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = retMsg
	response.Message = contentList

	return ctx.WriteJson(response)
}

// GetLiveQuestionAnswerList 获取指定直播间当天的问答列表，如果IsSelf=1，则筛选当前用户的提问列表
func GetLiveQuestionAnswerList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract.UserInfo(ctx)
	isSelf := ctx.PostFormValue("IsSelf")
	dateStr := ctx.PostFormValue("Date")
	date, err := _time.ParseTime(dateStr)
	if err != nil {
		date = time.Now()
	}
	if user == nil || user.RoomList == nil || len(user.RoomList) == 0 {
		response.RetCode = -1
		response.RetMsg = "您没有直播间问答查看权限"
		return ctx.WriteJson(response)
	}

	requestData := agent.QuestionAnswerRequest{
		RoomIdList: user.RoomList,
		UserId:     user.UserId,
		Date:       date.Format(time.RFC3339),
	}
	if isSelf == "1" {
		requestData.IsSelf = 1
	}
	questionList, err := agent.GetLiveQuestionAnswerList(requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "当前问答内容不存在！"
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = questionList
	return ctx.WriteJson(response)
}

// AddQuestion 用户提问
func AddQuestion(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract.UserInfo(ctx)
	askContent := ctx.PostFormValue("AskContent")
	if len(askContent) == 0 {
		response.RetCode = -1
		response.RetMsg = "提问内容不能为空"
		return ctx.WriteJson(response)
	}

	roomId, err := strconv.Atoi(config.CurrentConfig.DefaultAskRoom)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "房间号不正确"
		return ctx.WriteJson(response)
	}

	requestData := agent.AddQuestionRequest{
		AskContent: askContent,
		RoomId:     roomId,
		UserId:     user.UserId,
		CheckPID:   true,
		PID:        user.PID,
	}
	err = agent.AddQuestion(requestData)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		submitTongJi(ctx,"pv","question")
	}
	return ctx.WriteJson(response)
}

// IsTradeTime 当前是否为交易时间
func IsTradeTime(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	isTradeTime, err := agent.IsTradeTime()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = isTradeTime
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
func HasNewMessage(ctx dotweb.Context) error {
	queryKey := ctx.QueryString("querykey")
	roomId := ctx.QueryInt("groupid")
	if strings.Contains(queryKey, "CST") == false {
		queryKey += " CST"
	}
	date, err := _time.ParseTime(queryKey)
	if err != nil {
		date = time.Now()
	}

	requestData := agent.HasNewMessageRequest{
		Date:   date.Format(time.RFC3339),
		RoomId: roomId,
	}
	liveDate, err := agent.HasNewMessage(requestData)
	if err != nil {
		return ctx.WriteJson("")
	}
	if liveDate.ContentNum == 0 && liveDate.RepliedNum == 0 {
		return ctx.WriteJson("")
	}
	return ctx.WriteJson(liveDate)
}

// GetRoomList 获取用户拥有权限的直播间列表
func GetRoomList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract.UserInfo(ctx)
	if user == nil || user.RoomList == nil || len(user.RoomList) == 0 {
		response.RetCode = -1
		response.RetMsg = "您没有直播间查看权限"
	} else {
		//过期10分钟
		tokenLifeSeconds := 60 * 10
		tokenResp, err := tokenapi.CreateToken(config.CurrentConfig.AppId, tokenLifeSeconds, user.RoomList)
		if err != nil {
			response.RetCode = -2
			response.RetMsg = "获取直播Token失败，请刷新页面后重试！"
			return ctx.WriteJson(response)
		}
		roomInfo := viewmodel.UserRoomInfo{
			RoomList: user.RoomList,
			Token:    tokenResp.Message.Token,
		}
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = roomInfo
	}
	return ctx.WriteJson(response)
}

type PostData struct {
	Appid    string `json:"appid"`
	Time     int64  `json:"time"`
	Category string `json:"category"`
	Key      string `json:"key"`
	Increase int    `json:"increase"`
	Globalid string `json:"globalid"`
}

//submitTongJi 统计仪表板
//参数：
//ctx：上下文
//tongjiType：pv/uv
//flag：visit(访问);visitor(访客);question(提问);
//返回：
//result：响应消息
//e：错误消息
func submitTongJi(ctx dotweb.Context, tongjiType string, flag string) (body string, errReturn error) {
	var url string
	timeStamp := _time.GetTimestamp(time.Now())
	globalid, _ := ctx.ReadCookieValue("tongji_globalid")
	if globalid == "" {
		globalid = ctx.RemoteIP()
	}

	postData := new(PostData)
	postData.Appid = "10050"
	postData.Time = timeStamp
	postData.Globalid = globalid
	postData.Key = "TianYanDingPanLive-" + flag

	if tongjiType == "pv" {
		url = "http://api2.tongji.emoney.cn/counter/pv"
		postData.Increase = 1
		postData.Category = "TianYanDingPanLivePV"
	} else {
		url = "http://api2.tongji.emoney.cn/counter/uv"
		postData.Category = "TianYanDingPanLiveUV"
	}

	postBody := _json.GetJsonString(postData)
	postBody = "ActionData=" + postBody

	result, _, _, e := _http.HttpPost(url, postBody, "")

	return result, e
}
