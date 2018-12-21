package strategyservice

import (
	sysSql "database/sql"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	strategyservice_model "git.emoney.cn/softweb/roboadvisor/protected/model/strategyservice"
	"git.emoney.cn/softweb/roboadvisor/protected/repository"
	"strings"
)

type ClientStrategyInfoRepository struct {
	repository.BaseRepository
}

func NewClientStrategyInfoRepository(conf *protected.ServiceConfig) *ClientStrategyInfoRepository {
	repo := &ClientStrategyInfoRepository{}
	repo.Init(conf.EMoney_RoboAdvisorDBConn)
	return repo
}

// GetClientStrategyInfoList 获取所有可用的策略及栏目信息
// includeParent是否包含父级
func (repo *ClientStrategyInfoRepository) GetClientStrategyInfoList(includeParent bool) (result []*strategyservice_model.ClientStrategyInfo, err error) {
	sql := `SELECT TB2.[ClientStrategyId],TB2.[ColumnInfoId],TB2.[ClientStrategyName],TB2.[IsTop],TB1.[ParentId],TB1.[ParentName]  FROM [ClientStrategyInfoRelation] AS TB1
JOIN [ClientStrategyInfo]  AS TB2 ON [TB1].[ClientStrategyId]=TB2.[ClientStrategyId] WHERE TB1.[IsDeleted]=0 AND TB2.[IsDeleted]=0 `
	if includeParent == false {
		sql += " AND TB2.[IsTop]=0"
	}
	err = repo.FindList(&result, sql)
	return result, err
}

// GetClientStrategyInfo 根据策略Id获取策略及栏目信息
func (repo *ClientStrategyInfoRepository) GetClientStrategyInfo(clientStrategyId int) (strategyInfo *strategyservice_model.ClientStrategyInfo, err error) {
	sql := "SELECT TOP 1 [ClientStrategyInfoId],[ColumnInfoId],[ClientStrategyId],[ClientStrategyName],[ClientStrategyInfo].[IsDeleted],[ClientStrategyInfo].[CreateTime],[ColumnInfo].[ColumnName] FROM [ClientStrategyInfo] JOIN [ColumnInfo] ON [ColumnInfoId]=ID WHERE [ClientStrategyInfo].[IsDeleted]=0 AND [ClientStrategyId]=?"
	strategyInfo = new(strategyservice_model.ClientStrategyInfo)
	err = repo.FindOne(strategyInfo, sql, clientStrategyId)
	if sysSql.ErrNoRows == err {
		return nil, nil
	}
	return strategyInfo, err
}

// InsertClientStrategyInfo 插入一个新的策略信息
func (repo *ClientStrategyInfoRepository) InsertClientStrategyInfo(strategyInfo strategyservice_model.ClientStrategyInfo) (err error) {
	sql := "INSERT INTO [ClientStrategyInfo] ([ColumnInfoId],[ClientStrategyId],[ClientStrategyName],[IsTop],[IsDeleted],[CreateTime],[ModifyTime]) VALUES(?,?,?,?,0,GETDATE(),GETDATE())"
	_, err = repo.Insert(sql, strategyInfo.ColumnInfoId, strategyInfo.ClientStrategyId, strategyInfo.ClientStrategyName, strategyInfo.IsTop)
	return err
}

// DeleteClientStrategyInfoNotInList 删除不在列表中的策略关联关系
func (repo *ClientStrategyInfoRepository) DeleteClientStrategyInfoNotInList(clientStrategyIdList []int) (err error) {
	if clientStrategyIdList == nil || len(clientStrategyIdList) == 0 {
		err = errors.New("clientStrategyIdList 不能为空")
		return err
	}
	insql := strings.Repeat("?,", len(clientStrategyIdList))
	insql = insql[:len(insql)-1]
	sql := "UPDATE [ClientStrategyInfo] SET [IsDeleted]=1,[ModifyTime]=GETDATE() WHERE [ClientStrategyId] NOT IN (" + insql + ")"
	var params []interface{}
	for _, item := range clientStrategyIdList {
		params = append(params, item)
	}
	_, err = repo.Update(sql, params...)
	return err
}

// UpdateClientStrategyInfo 更新策略信息
func (repo *ClientStrategyInfoRepository) UpdateClientStrategyInfo(strategyInfo strategyservice_model.ClientStrategyInfo) (err error) {
	sql := "UPDATE [ClientStrategyInfo] SET [ClientStrategyName]=?,[IsTop]=?,[ModifyTime]=GETDATE() WHERE [ClientStrategyInfoId]=?"
	_, err = repo.Update(sql, strategyInfo.ClientStrategyName, strategyInfo.IsTop, strategyInfo.ClientStrategyInfoId)
	return err
}

// InsertClientStrategyInfoRelation 插入一条新的关联关系
func (repo *ClientStrategyInfoRepository) InsertClientStrategyInfoRelation(clientStrategyId int, parentId int, parentName string) (err error) {
	sql := "INSERT INTO [ClientStrategyInfoRelation] ([ClientStrategyId],[ParentId],[ParentName],[IsDeleted],[CreateTime],[ModifyTime]) VALUES(?,?,?,0,GETDATE(),GETDATE())"
	_, err = repo.Insert(sql, clientStrategyId, parentId, parentName)
	return err
}

// DeleteAllClientStrategyInfoRelation 删除所有原有的关联关系
func (repo *ClientStrategyInfoRepository) DeleteAllClientStrategyInfoRelation() (err error) {
	sql := "UPDATE [ClientStrategyInfoRelation] SET [IsDeleted]=1,[ModifyTime]=GETDATE()"
	_, err = repo.Update(sql)
	return err
}
