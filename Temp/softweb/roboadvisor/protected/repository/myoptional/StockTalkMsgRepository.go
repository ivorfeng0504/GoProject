package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"time"
)

type StockTalkMsgRepository struct {
	repository.BaseRepository
}

func NewStockTalkMsgRepository(conf *protected.ServiceConfig) *StockTalkMsgRepository {
	repo := &StockTalkMsgRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertStockTalkMsg 插入一条数据
func (repo *StockTalkMsgRepository) InsertStockTalkMsg(model *myoptional_model.StockTalkMsg) (id int64, err error) {
	sql := "INSERT INTO [StockTalkMsg] ([LiveRoomId],[LiveRoomName],[StockCode],[StockName],[ImageList],[Content],[SendTime],[IsSendStockTalk],[CreateTime],[IsDeleted],[ModifyTime],[ModifyUser]) VALUES(?,?,?,?,?,?,?,0,GETDATE(),0,GETDATE(),?)"
	id, err = repo.Insert(sql, model.LiveRoomId, model.LiveRoomName, model.StockCode, model.StockName, model.ImageList, model.Content, time.Time(model.SendTime), model.ModifyUser)
	return id, err
}
