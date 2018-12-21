package mapper

import (
	"fmt"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel"
	"github.com/devfeel/mapper"
	"strings"
	"time"
)

func MapperLiveContent(from []*livemodel.LiveContent) []*viewmodel.LiveContent {
	var target []*viewmodel.LiveContent
	err := mapper.MapperSlice(from, &target)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	//处理股票列表与时间格式
	for _, content := range target {
		if content.StockList != "" {
			stockInfoList := strings.Split(content.StockList, ",")
			if len(stockInfoList) > 0 {
				for _, stockInfoStr := range stockInfoList {
					if stockInfoStr == "" {
						continue
					}
					stockInfo := strings.Split(stockInfoStr, "|")
					if len(stockInfo) == 4 {
						stock := &viewmodel.StockInfo{
							StockName: stockInfo[0],
							StockCode: stockInfo[1],
							Change:    stockInfo[2],
							StockType: stockInfo[3],
						}
						content.StockInfoList = append(content.StockInfoList, stock)
					}
				}
			}
		}
		content.CreateTimeStr = time.Time(content.CreateTime).Format("15:04")
	}
	return target
}

func MapperLiveQuestionAnswer(from []*livemodel.LiveQuestionAnswer, room *livemodel.LiveRoom) []*viewmodel.LiveQuestionAnswer {
	var target []*viewmodel.LiveQuestionAnswer
	if room == nil {
		return target
	}
	err := mapper.MapperSlice(from, &target)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	for _, question := range target {
		question.AskTimeStr = time.Time(question.AskTime).Format("15:04")

		question.AdviserNo = room.AdviserNo
		question.AdviserName = room.AdviserName

		//是否有投顾资质信息
		if question.AdviserName != "" && question.AdviserNo != "" {
			question.HasAdviserInfo = true
		}
		//是否已经回答
		if len(question.AnswerContent) > 0 {
			question.IsAnswered = true
		}
	}
	return target
}
