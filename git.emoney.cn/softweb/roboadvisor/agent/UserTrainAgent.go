package agent

import (
	"time"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/config"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"git.emoney.cn/softweb/roboadvisor/protected/model/train"
)

// GetTrainListByDate 根据日期获取培训课程列表
func GetTrainListByDate(pid string, date *time.Time) (result []* train.NetworkMeetingInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = UserTrainRequest{
		PID:        pid,
		FilterDate: date,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/train/gettrainlistbydate", req)

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

// GetTrainListByDateAndTag 根据日期和标签获取培训课程列表
func GetTrainListByDateAndTag(pid string, date *time.Time,trainTag int,trainType string,pageSize int, currPage int) (result []* train.NetworkMeetingInfo,totalCount int, err error) {
	req := contract.NewApiRequest()
	req.RequestData = UserTrainRequest{
		PID:         pid,
		TrainTag:    trainTag,
		TrainType:   trainType,
		FilterDate:  date,
		PageSize:    pageSize,
		CurrentPage: currPage,
	}
	response, err := Post(config.CurrentConfig.WebApiHost+"/api/train/gettrainlistbydateandtag", req)
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

type UserTrainRequest struct {
	//PID
	PID string
	//培训标签
	TrainTag int
	//培训类型
	TrainType string
	//过滤时间 2018-11-20
	FilterDate *time.Time
	//当前页
	CurrentPage int
	//页码
	PageSize int
}
