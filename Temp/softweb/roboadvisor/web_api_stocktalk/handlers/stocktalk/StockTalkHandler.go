package stocktalk

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	myoptional_srv "git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	myoptional_vmmodel "git.emoney.cn/softweb/roboadvisor/protected/viewmodel/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// InsertStockTalk 插入一个微股吧评论
func InsertStockTalk(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := myoptional_model.StockTalk{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if len(requestData.NickName) == 0 {
		response.RetCode = -3
		response.RetMsg = "用户昵称不正确"
		return ctx.WriteJson(response)
	}

	//对提问内容进行编码
	requestData.Content = template.HTMLEscapeString(requestData.Content)
	//处理表情
	imojiBaseUrl := "http://static.emoney.cn/live/Scripts/xheditor/xheditor_emot/newface/default/"
	for i := 1; i <= 36; i++ {
		old := "[em_" + strconv.Itoa(i) + "]"
		newStr := `<img src="` + imojiBaseUrl + strconv.Itoa(i) + ".gif" + `" border="0" />`
		requestData.Content = strings.Replace(requestData.Content, old, newStr, -1)
	}
	userSrv := service.NewUserInfoService()
	uid, err := strconv.ParseInt(requestData.UID, 10, 64)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	userInfo, err := userSrv.GetUserInfoByUID(uid)

	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if userInfo != nil {
		if len(userInfo.Headportrait) > 0 {
			requestData.Avatar = fmt.Sprintf(config.CurrentConfig.UserHome_Headportrait_UrlFormat, userInfo.Headportrait)
		}
		if len(userInfo.NickName) > 0 {
			requestData.NickName = userInfo.NickName
		}
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()
	err = stockTalkSrv.InsertStockTalk(&requestData)
	if err != nil {
		response.RetCode = -6
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetStockTalkListByDate 获取指定日期的微股吧评论
func GetStockTalkListByDate(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.StockTalkRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()
	var result []*myoptional_model.StockTalk
	stockTalkList, err := stockTalkSrv.GetStockTalkListByDate(requestData.Date)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	//是否是递增查询
	if requestData.IsIncre {
		if stockTalkList != nil && len(stockTalkList) > 0 {
			for _, item := range stockTalkList {
				if time.Time(item.CreateTime).After(requestData.Date) {
					result = append(result, item)
				}
			}
		}
	} else {
		if requestData.CheckDate == false {
			result, err = stockTalkSrv.GetStockTalkNewst()
			if err != nil {
				response.RetCode = -4
				response.RetMsg = err.Error()
				return ctx.WriteJson(response)
			}
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStockTalkListByPage 分页获取微股吧评论
func GetStockTalkListByPage(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.StockTalkRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()
	stockTalkList, err := stockTalkSrv.GetStockTalkListByPage(requestData.StockCodeList, requestData.PageIndex, requestData.PageSize)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = stockTalkList

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// SendStockTalkMsg 发送微股吧消息到后台备用
func SendStockTalkMsg(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	requestData := myoptional_vmmodel.StockTalkMsgVM{}
	err := agent.Bind(ctx, request, &requestData)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if len(requestData.LiveRoomName) == 0 {
		response.RetCode = -2
		response.RetMsg = "直播间名称不正确"
		return ctx.WriteJson(response)
	}

	if len(requestData.Content) == 0 {
		response.RetCode = -3
		response.RetMsg = "直播内容不能为空"
		return ctx.WriteJson(response)
	}
	if len(requestData.StockInfoList) == 0 {
		response.RetCode = -4
		response.RetMsg = "股票代码不正确"
		return ctx.WriteJson(response)
	}

	//对提问内容进行编码
	requestData.Content = template.HTMLEscapeString(requestData.Content)
	stockTalkMsgSrv := myoptional_srv.NewStockTalkMsgService()
	err = stockTalkMsgSrv.InsertStockTalkMsgQueue(&requestData)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetStockTalkByStockCode 根据股票代码获取微股吧数据
func GetStockTalkByStockCode(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	if ctx.Request().Method == http.MethodGet {
		request.RequestData = ctx.QueryString("StockCode")
		request.AllowClientCache = true
		request.MessageHash = ctx.QueryString("MessageHash")
	} else {
		err := ctx.Bind(request)
		if err != nil {
			response.RetCode = -1
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
	}
	stockCode := request.RequestData.(string)
	if len(stockCode) == 0 {
		response.RetCode = -2
		response.RetMsg = "股票代码不能为空"
		return ctx.WriteJson(response)
	}
	if len(stockCode) > 6 {
		stockCode = stockCode[len(stockCode)-6:]
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()
	result, err := stockTalkSrv.GetStockTalkByStockCode(stockCode, 30)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStockCodeList 获取股票列表
func GetStockCodeList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	if ctx.Request().Method == http.MethodGet {
		request.AllowClientCache = true
		request.MessageHash = ctx.QueryString("MessageHash")
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()
	result, err := stockTalkSrv.GetStockCodeListForPCClient()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// SyncStockTalkForStockBasicInfoByStockList 根据指定的股票列表同步最新的基本面信息到微股吧
func SyncStockTalkForStockBasicInfoByStockList(ctx dotweb.Context) error {
	response := contract.NewApiResponse()
	stockListStr := ctx.QueryString("StockList")
	// IsSync 是否同步实行 1同步 其他异步
	isSync := ctx.QueryString("IsSync")
	stockCodeList := strings.Split(stockListStr, ",")
	if stockCodeList == nil || len(stockCodeList) == 0 {
		response.RetCode = -1
		response.RetMsg = "股票列表为空"
		return ctx.WriteJson(response)
	}
	var stockList []*model.StockInfo
	for _, code := range stockCodeList {
		if len(code) == 0 {
			continue
		}
		stockList = append(stockList, &model.StockInfo{
			StockCode:  code,
			CreateTime: time.Now(),
		})
	}
	if stockList == nil || len(stockList) == 0 {
		response.RetCode = -2
		response.RetMsg = "股票列表为空"
		return ctx.WriteJson(response)
	}
	stockTalkSrv := myoptional_srv.NewStockTalkService()

	//异步执行 不返回真实操作结果 执行结果请查看相关日志
	if isSync != "1" {
		go stockTalkSrv.SyncStockTalkForStockBasicInfoByStockList(stockList)
		response.RetCode = 0
		response.RetMsg = "SUCCESS 执行结果请查看相关日志"
		response.Message = stockList
		return ctx.WriteJson(response)
	}
	err := stockTalkSrv.SyncStockTalkForStockBasicInfoByStockList(stockList)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS By" + time.Now().Format("yyyy-MM-dd HH:mm:ss")
	response.Message = stockList
	return ctx.WriteJson(response)
}
