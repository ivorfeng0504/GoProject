package expertnews

import (
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expertnews"
	"encoding/json"
)

type TopicRepository struct {
	repository.BaseRepository
}

func NewTopicRepository(conf *protected.ServiceConfig) *TopicRepository {
	repo := &TopicRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}


// GetTopicList 获取主题列表
func (repo *TopicRepository) GetTopicList() (topicList []*expertnews.ExpertNews_Topic, err error) {
	sql := "SELECT * FROM ExpertNews_Topic with(nolock)  WHERE IsDeleted=0 "
	err = repo.FindList(&topicList, sql)

	if len(topicList)>0 {

		for _,v := range topicList  {
			stocks := v.RelatedStockInfo
			json.Unmarshal([]byte(stocks),&v.RelatedStockList)
		}

	}
	return topicList, err
}





