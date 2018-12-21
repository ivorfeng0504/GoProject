package repository

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type RefundHistoryRepository struct {
	repository.BaseRepository
}

func NewRefundHistoryRepository(conf *protected.ServiceConfig) *RefundHistoryRepository {
	repo := &RefundHistoryRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertRefundHistory 新增一条退款记录
func (repo *RefundHistoryRepository) InsertRefundHistory(model *userhome_model.RefundHistory) (id int64, err error) {
	if model == nil {
		return 0, errors.New("RefundHistory 不能为空")
	}
	sql := "INSERT INTO [UserHome_RefundHistory] ([Account],[Mobile],[OrderId],[Reason],[RefundMode],[SubmitData],[CreateTime]) VALUES(?,?,?,?,?,?,GETDATE())"
	id, err = repo.Insert(sql, model.Account, model.Mobile, model.OrderId, model.Reason, model.RefundMode, model.SubmitData)
	return id, err
}
