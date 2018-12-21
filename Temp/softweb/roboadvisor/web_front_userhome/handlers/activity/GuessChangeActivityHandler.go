package activity

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel/userhome"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/mapper"
	"time"
)

// GuessChange 猜涨跌活动首页
func GuessChange(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "guess_change.html")
}

// GetCurrentGuessChange 获取当期猜涨跌活动情况
func GetCurrentGuessChange(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	currentGuessChange, err := agent.GetCurrentGuessChange(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取当期猜涨跌活动信息异常"
		global.InnerLogger.ErrorFormat(err, "获取当期猜涨跌活动信息异常")
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessChange
	return ctx.WriteJson(response)
}

// GuessChangeSubmit 提交竞猜结果
func GuessChangeSubmit(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	result := 0
	resultStr := ctx.PostFormValue("result")
	if resultStr == "1" {
		result = 1
	} else if resultStr == "-1" {
		result = -1
	} else {
		response.RetCode = -1
		response.RetMsg = "无效的竞猜结果"
		return ctx.WriteJson(response)
	}
	user := contract_userhome.UserHomeUserInfo(ctx)
	currentGuessChange, err := agent.GuessChangeSubmit(*user, result)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "提交竞猜结果异常"
		global.InnerLogger.ErrorFormat(err, "提交竞猜结果异常")
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessChange
	return ctx.WriteJson(response)
}

// GetMyGuessChangeJoinInfo 获取我的竞猜和我的奖品信息
func GetMyGuessChangeJoinInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	joinInfo := &viewmodel.GuessChangeJoinInfo{}
	user := contract_userhome.UserHomeUserInfo(ctx)
	myGuessChangeList, err := agent.GetMyGuessChangeInfoCurrentWeek(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取我的竞猜异常"
		global.InnerLogger.ErrorFormat(err, "获取我的竞猜信息异常")
		return ctx.WriteJson(response)
	}

	var guessChangeResultListSource []*viewmodel.GuessChangeResult
	if myGuessChangeList != nil && len(myGuessChangeList) > 0 {
		err = mapper.MapperSlice(myGuessChangeList, &guessChangeResultListSource)
	}

	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取我的竞猜-更多信息异常 请稍后再试"
		global.InnerLogger.ErrorFormat(err, "我的竞猜 Mapper异常 映射失败")
		return ctx.WriteJson(response)
	}

	myGuessChangeAwardList, err := agent.GetMyGuessageAwardList(*user)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取我的奖品异常"
		global.InnerLogger.ErrorFormat(err, "获取我的奖品异常")
		return ctx.WriteJson(response)
	}

	if myGuessChangeAwardList != nil && len(myGuessChangeAwardList) > 0 {
		err = mapper.MapperSlice(myGuessChangeAwardList, &joinInfo.GuessChangeAwardList)
	}

	if err != nil {
		response.RetCode = -4
		response.RetMsg = "获取我的奖品异常 请稍后再试"
		global.InnerLogger.ErrorFormat(err, "我的奖品 Mapper异常 映射失败")
		return ctx.WriteJson(response)
	}

	//获取到本周一
	weekResultDict := make(map[string]int)
	weekday := int(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7
	}
	count := 5
	now := time.Now()
	monday := time.Now().AddDate(0, 0, -weekday+1)
	for i := 0; i < count; i++ {
		date := monday.AddDate(0, 0, i)
		issueNum := date.Format("20060102")
		result := &viewmodel.GuessChangeResult{
			IssueNumber: issueNum,
			ResultDesc:  "-",
			StateDesc:   "未参与",
		}
		if (now.Hour() < 9 && date.Day() > now.Day()) || (now.Hour() >= 9 && date.Day() > now.Day()+1) {
			result.StateDesc = "未开局"
		}
		joinInfo.GuessChangeResultList = append(joinInfo.GuessChangeResultList, result)
		weekResultDict[issueNum] = i
	}

	//计算参与次数与猜中次数
	if guessChangeResultListSource != nil && len(guessChangeResultListSource) > 0 {
		joinInfo.JoinCount = len(guessChangeResultListSource)
		for _, record := range guessChangeResultListSource {
			if record.IsGuessed {
				joinInfo.GuessedCount++
				record.StateDesc = "猜中"
			} else {
				if record.IsPublish == false {
					record.StateDesc = "未公布"
				} else {
					record.StateDesc = "未猜中"
				}
			}
			//如果已经参与了 则替换默认空数据
			index, exist := weekResultDict[record.IssueNumber]
			if exist {
				joinInfo.GuessChangeResultList[index] = record
			}
		}
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = joinInfo
	return ctx.WriteJson(response)
}

// GetMyGuessChangeInfoNewst 获取我的竞猜-更多信息
func GetMyGuessChangeInfoNewst(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	myGuessChangeList, err := agent.GetMyGuessChangeInfoNewst(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取我的竞猜-更多信息异常"
		global.InnerLogger.ErrorFormat(err, "获取我的竞猜-更多信息异常")
		return ctx.WriteJson(response)
	}

	var myGuessChangeListVM []*viewmodel.GuessChangeResult
	if myGuessChangeList != nil && len(myGuessChangeList) > 0 {
		err = mapper.MapperSlice(myGuessChangeList, &myGuessChangeListVM)
	}
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取我的竞猜-更多信息异常 请稍后再试"
		global.InnerLogger.ErrorFormat(err, "获取我的竞猜-更多信息异常 映射失败")
		return ctx.WriteJson(response)
	}

	if myGuessChangeListVM != nil && len(myGuessChangeListVM) > 0 {
		for _, record := range myGuessChangeListVM {
			if record.IsGuessed {
				record.StateDesc = "猜中"
			} else {
				if record.IsPublish == false {
					record.StateDesc = "未公布"
				} else {
					record.StateDesc = "未猜中"
				}
			}
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = myGuessChangeListVM
	return ctx.WriteJson(response)
}

// GetGuessChangeAwardInfo 获取猜涨跌奖品信息
func GetGuessChangeAwardInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	awardList, err := agent.GetGuessChangeAward()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取历史信息异常"
		global.InnerLogger.ErrorFormat(err, "获取猜涨跌奖品信息异常")
		return ctx.WriteJson(response)
	}

	var awardSimpleList []*viewmodel.ActivityAwardSimple
	if awardList != nil && len(awardList) > 0 {
		mapper.MapperSlice(awardList, &awardSimpleList)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = awardSimpleList
	return ctx.WriteJson(response)
}