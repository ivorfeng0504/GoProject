package myoptional

import (
	sysSql "database/sql"
	"git.emoney.cn/softweb/roboadvisor/protected"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
)

type EvaluationResultRepository struct {
	repository.BaseRepository
}

func NewEvaluationResultRepository(conf *protected.ServiceConfig) *EvaluationResultRepository {
	repo := &EvaluationResultRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// InsertEvaluationResult 新增一个评测记录
func (repo *EvaluationResultRepository) InsertEvaluationResult(result myoptional_model.EvaluationResult) (err error) {
	sql := "INSERT INTO [EvaluationResult] ([UID],[InvestTarget],[InvestTargetDesc],[InvestTargetTip],[ChooseStockReason],[ChooseStockReasonDesc],[ChooseStockReasonTip],[HoldStockTime],[HoldStockTimeDesc],[HoldStockTimeTip],[BuyStyle],[BuyStyleDesc],[BuyStyleTip],[Result],[ResultDesc],[CreateTime],[ModifyTime]) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,GETDATE(),GETDATE())"
	_, err = repo.Insert(sql, result.UID, result.InvestTarget, result.InvestTargetDesc, result.InvestTargetTip, result.ChooseStockReason, result.ChooseStockReasonDesc, result.ChooseStockReasonTip, result.HoldStockTime, result.HoldStockTimeDesc, result.HoldStockTimeTip, result.BuyStyle, result.BuyStyleDesc, result.BuyStyleTip, result.Result, result.ResultDesc)
	return err
}

// UpdateEvaluationResult 更新评测记录
func (repo *EvaluationResultRepository) UpdateEvaluationResult(result myoptional_model.EvaluationResult) (err error) {
	sql := "UPDATE [EvaluationResult] SET [InvestTarget]=?,[InvestTargetDesc]=?,[InvestTargetTip]=?,[ChooseStockReason]=?,[ChooseStockReasonDesc]=?,[ChooseStockReasonTip]=?,[HoldStockTime]=?,[HoldStockTimeDesc]=?,[HoldStockTimeTip]=?,[BuyStyle]=?,[BuyStyleDesc]=?,[BuyStyleTip]=?,[Result]=?,[ResultDesc]=?,[ModifyTime]=GETDATE() WHERE [UID]=? "
	_, err = repo.Update(sql, result.InvestTarget, result.InvestTargetDesc, result.InvestTargetTip, result.ChooseStockReason, result.ChooseStockReasonDesc, result.ChooseStockReasonTip, result.HoldStockTime, result.HoldStockTimeDesc, result.HoldStockTimeTip, result.BuyStyle, result.BuyStyleDesc, result.BuyStyleTip, result.Result, result.ResultDesc, result.UID)
	return err
}

// GetEvaluationResult 获取评测结果
func (repo *EvaluationResultRepository) GetEvaluationResult(uid string) (result *myoptional_model.EvaluationResult, err error) {
	sql := "SELECT TOP 1 * FROM [EvaluationResult] WHERE [UID]=?"
	result = new(myoptional_model.EvaluationResult)
	err = repo.FindOne(result, sql, uid)
	if err == sysSql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
