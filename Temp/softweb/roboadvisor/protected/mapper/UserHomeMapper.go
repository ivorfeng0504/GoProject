package mapper

import (
	"git.emoney.cn/softweb/roboadvisor/const"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"strconv"
	"time"
)

func MapperActivity(from []*userhome_model.Activity, userActivityMap map[int64]bool) []*viewmodel.Activity {
	if from == nil {
		return nil
	}
	var target []*viewmodel.Activity
	for _, activity := range from {
		vmActivity := &viewmodel.Activity{
			ActivityId:  activity.ActivityId,
			Title:       activity.Title,
			Summary:     activity.Summary,
			BeginTime:   time.Time(activity.BeginTime).Format("2006-01-02"),
			EndTime:     time.Time(activity.EndTime).Format("2006-01-02"),
			ImageUrl:    activity.ImageUrl,
			ActivityUrl: activity.ActivityUrl,
			OpenMode:    activity.OpenMode,
			NeedBind:    activity.NeedBind,
			NeedSSO:     activity.NeedSSO,
		}
		if userActivityMap != nil {
			_, exist := userActivityMap[activity.ActivityId]
			vmActivity.IsJoin = exist
		}
		target = append(target, vmActivity)
	}
	return target
}

func MapperUserAward(from []*userhome_model.UserAward) []*viewmodel.UserAward {
	if from == nil {
		return nil
	}
	var target []*viewmodel.UserAward
	for _, userAward := range from {
		expireTime := time.Time(userAward.CreateTime).Add(time.Duration(userAward.AvailableDay*24) * time.Hour)
		vmUserAward := &viewmodel.UserAward{
			UserAwardId:    userAward.UserAwardId,
			AwardName:      userAward.AwardName,
			AwardId:        userAward.AwardId,
			AwardImg:       userAward.AwardImg,
			CreateTime:     time.Time(userAward.CreateTime).Format("2006-01-02"),
			IntroduceVideo: userAward.IntroduceVideo,
			ExpireTime:     expireTime.Format("2006-01-02"),
		}
		if expireTime.Before(time.Now()) {
			vmUserAward.AwardStateDesc = "已过期"
		} else {
			vmUserAward.AwardStateDesc = "未过期"
		}
		if userAward.AvailableDay <= 0 {
			vmUserAward.AwardStateDesc = "未过期"
			vmUserAward.ExpireTime = "长期有效"
			vmUserAward.AvailableDay = "长期有效"
		} else {
			vmUserAward.AvailableDay = strconv.Itoa(userAward.AvailableDay) + "天"
		}
		if len(vmUserAward.IntroduceVideo) == 0 {
			vmUserAward.IntroduceVideo = "javascript:void(0)"
		}
		switch userAward.AwardType {
		case _const.AwardType_Function:
			vmUserAward.AwardTypeDesc = "功能"
			break
		case _const.AwardType_Video:
			vmUserAward.AwardTypeDesc = "视频"
			vmUserAward.IsVideo = true
			break
		case _const.AwardType_Report:
			vmUserAward.AwardTypeDesc = "研报"
			break
		case _const.AwardType_QQ:
			vmUserAward.AwardTypeDesc = "其他"
			break
		}
		target = append(target, vmUserAward)
	}
	return target
}

func MapperUserReceivedStock(from *userhome_model.UserReceivedStock) *viewmodel.UserReceivedStockViewModel {
	if from == nil {
		return nil
	}
	issueDate, err := time.Parse("20060102", from.IssueNumber)
	vm := &viewmodel.UserReceivedStockViewModel{
		CreateTime: time.Time(from.CreateTime).Format("2006-01-02"),
		ReportUrl:  from.ReportUrl,
	}
	if err == nil {
		vm.IssueNumber = issueDate.Format("2006-01-02")
	}

	if len(from.StockList) > 0 {
		err := _json.Unmarshal(from.StockList, &vm.StockList)
		if err != nil {
			vm.StockList = nil
		}
	}
	return vm
}

func MapperUserReceivedStockList(from []*userhome_model.UserReceivedStock) []*viewmodel.UserReceivedStockViewModel {
	if from == nil {
		return nil
	}
	var target []*viewmodel.UserReceivedStockViewModel
	for _, record := range from {
		target = append(target, MapperUserReceivedStock(record))
	}
	return target
}
