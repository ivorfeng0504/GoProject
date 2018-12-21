package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"strconv"
	"strings"
)

func IsTradeTime() (isTradeTime bool, err error) {
	req := contract.NewApiRequest()
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/istradetime", req)
	if err != nil {
		return false, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return false, errors.New(response.RetMsg)
	}
	return response.Message.(bool), nil
}

func AddQuestion(model AddQuestionRequest) (err error) {
	req := contract.NewApiRequest()
	model.Source = _const.UserSource_PC
	req.RequestData = model
	response, err := PostNoCache(config.CurrentConfig.LiveWebApiHost+"/api/live/addquestion", req)
	if err != nil {
		return errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return errors.New(response.RetMsg)
	}
	return nil
}

func GetLiveQuestionAnswerList(model QuestionAnswerRequest) (answerList []viewmodel.LiveQuestionAnswer, err error) {
	req := contract.NewApiRequest()
	req.RequestData = model
	response, err := PostNoCache(config.CurrentConfig.LiveWebApiHost+"/api/live/livequestion", req)
	if err != nil {
		return nil, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}

	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &answerList)
	return answerList, err
}

func HasNewMessage(model HasNewMessageRequest) (data viewmodel.LiveDataWithTime, err error) {
	req := contract.NewApiRequest()
	req.RequestData = model
	response, err := PostNoCache(config.CurrentConfig.LiveWebApiHost+"/api/live/hasnewmessage", req)
	if err != nil {
		return data, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return data, errors.New(response.RetMsg)
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &data)
	return data, err
}

func GetLiveContent(model LiveContentRequest) (contentList []viewmodel.LiveContent, retMsg string, err error) {
	req := contract.NewApiRequest()
	req.RequestData = model
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/livecontent", req)
	retMsg = "FAILED"
	if response != nil {
		retMsg = response.RetMsg
	}
	if err != nil {
		return nil, retMsg, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return nil, retMsg, errors.New(retMsg)
	}
	if response.Message == nil {
		return nil, retMsg, nil
	}

	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &contentList)
	return contentList, retMsg, err
}

// GetUserRoomList 获取当前用户拥有的直播间权限
func GetUserRoomList(mobile string) (rooms []int, err error) {
	req := contract.NewApiRequest()
	req.RequestData = mobile
	response, err := Post(config.CurrentConfig.LiveWebApiHost+"/api/live/getuserroomlist", req)
	if err != nil {
		return rooms, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return rooms, errors.New(response.RetMsg)
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &rooms)
	return rooms, err
}

//用户提问请求Model
type AddQuestionRequest struct {
	AskContent string
	RoomId     int
	UserId     int
	//加密手机号
	Mobile string
	//来源
	Source string
	//昵称
	NickName string
	//手机号掩码（当前等同于NickName）
	MaskMobile string
	//UID
	UID int64
	PID int
	//是否对PID进行权限校验
	CheckPID bool
}

//直播内容请求Model
type LiveContentRequest struct {
	RoomIdList []int
	IsIncre    int
	IsQueryTop int
	Date       string
	//如果当前直播未开播 是否显示上一个交易日的内容
	DisPlayLastTopic bool
}

type QuestionAnswerRequest struct {
	RoomIdList []int
	IsSelf     int
	UserId     int
	Mobile     string
	Date       string
}

type HasNewMessageRequest struct {
	Date   string
	RoomId int
}

type AddLiveUserInRoomRequest struct {
	//手机号
	Mobile string
	//来源标识
	Source string
	//注册的房间号
	Rooms []UserRoomInfo
}

type UserRoomInfo struct {
	//房间号
	RoomId int
	//PID（优先判断RoomId，如果RoomId未传递，则通过PID去查询映射的RoomId）
	PID int
	//过期天数
	ExpireDay int
	//订单Id
	OrderId string
	//来源
	Source string
	//图文直播权限分配唯一标识
	LiveOrderId string
}

//构造直播权限分配虚拟订单号
func GetLiveOrderId(id int) (orderId string) {
	maxLen := 12
	orderId = "live"
	idStr := strconv.Itoa(id)
	orderId = orderId + strings.Repeat("0", maxLen-len(idStr)-4) + idStr
	return orderId
}
