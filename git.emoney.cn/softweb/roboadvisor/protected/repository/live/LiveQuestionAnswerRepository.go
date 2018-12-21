package live

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type LiveQuestionAnswerRepository struct {
	repository.BaseRepository
}

func NewLiveQuestionAnswerRepository(conf *protected.ServiceConfig) *LiveQuestionAnswerRepository {
	repo := &LiveQuestionAnswerRepository{}

	repo.Init(conf.ContentLivePlatDBConn)
	return repo
}

// GetLiveQuestionAnswerList 获取指定直播间某天的问答列表，如果UserId>0则筛选指定用户的提问列表
func (repo *LiveQuestionAnswerRepository) GetLiveQuestionAnswerList(roomId int, date time.Time, userId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	sql := "SELECT * FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [LId]=? AND DATEDIFF(dd,[AskTime],?)=0 AND (?<=0 OR [AskUserId]=?) ORDER BY [AskTime] ASC"
	err = repo.FindList(&answerList, sql, roomId, date.Format("2006-01-02"), userId, userId)
	return answerList, err
}

// AddLiveQuestionAnswer 用户提问
func (repo *LiveQuestionAnswerRepository) AddLiveQuestionAnswer(askUserId int, askUserName string, askContent string, roomId int, source string) (id int, err error) {
	sql := "INSERT INTO [LiveQuestionAnswer] ([LId],[Status],[AskContent],[AskUserId],[AskTime],[AskUserName],[Source]) VALUES(?,?,?,?,?,?,?)"
	id64, err := repo.Insert(sql, roomId, 1, askContent, askUserId, time.Now(), askUserName, source)
	id = int(id64)
	return id, err
}

// GetTodayRepliedQuestionNum 获取今日已回复问题总数
func (repo *LiveQuestionAnswerRepository) GetTodayRepliedQuestionNum(roomId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [LId]=? AND [Status]=2 AND DATEDIFF(dd,[AskTime],getdate())=0"
	return repo.Count(sql, roomId)
}

// GetTodayAllQuestionNum 获取今日所有问题总数
func (repo *LiveQuestionAnswerRepository) GetTodayAllQuestionNum(roomId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [LId]=? AND DATEDIFF(dd,[AskTime],getdate())=0"
	return repo.Count(sql, roomId)
}

// GetTodayNoRepliedQuestionNum 获取今日未回复问题总数
func (repo *LiveQuestionAnswerRepository) GetTodayNoRepliedQuestionNum(roomId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [LId]=? AND [Status]=1 AND DATEDIFF(dd,[AskTime],getdate())=0"
	return repo.Count(sql, roomId)
}

// GetAllLiveQuestionAnswerListByAskUserId 获取指定直播间某天的问答列表，如果UserId>0则筛选指定用户的提问列表
func (repo *LiveQuestionAnswerRepository) GetAllLiveQuestionAnswerListByAskUserId(askUserId int) (answerList []*livemodel.LiveQuestionAnswer, err error) {
	sql := "SELECT * FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [AskUserId]=?"
	err = repo.FindList(&answerList, sql, askUserId)
	return answerList, err
}

// GetMyAskCount 获取某个用户的提问总数
func (repo *LiveQuestionAnswerRepository) GetMyAskCount(askUserId int) (count int64, err error) {
	sql := "SELECT COUNT(1) FROM [LiveQuestionAnswer] WITH(NOLOCK) WHERE [AskUserId]=?"
	return repo.Count(sql, askUserId)
}