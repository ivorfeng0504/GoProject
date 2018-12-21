package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/protected/model/yqq"
)

type IndexDataInfo struct {
	//今日头条
	JrttList []*agent.ExpertNewsInfo
	//盘后预测
	PhycList []*agent.ExpertNewsInfo
	//策略看盘
	ClkpList []*ExpertNews_StrategyNewsInfo
	//主题
	TopicList []*ExpertNews_Topic
	//专家直播
	ZjzbList []yqq.YqqRoom
}