package dbsync

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"strings"
	"time"
)

var (
	shareDbSyncLogger dotlog.Logger
)

const (
	dbSyncServiceName = "DBSyncService"
	waitSeconds       = 3
)

func init() {
	protected.RegisterServiceLoader(dbSyncServiceName, func() {
		shareDbSyncLogger = dotlog.GetLogger(dbSyncServiceName)
	})
}

// SyncCore 数据库同步
func SyncCore(table string) error {
	tableKey := "DBSync_" + table
	tableId := config.GetAppConfig(tableKey)
	if len(tableId) == 0 {
		shareDbSyncLogger.WarnFormat("数据库同步异常  未获取到配置%s", tableKey)
		return errors.New("未获取到配置" + tableKey)
	}
	//如果没配置地址 则直接返回 不进行同步
	if len(config.CurrentConfig.DBSyncApi) == 0 {
		return nil
	}
	url := fmt.Sprintf(config.CurrentConfig.DBSyncApi, tableId)
	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	if errReturn != nil {
		shareDbSyncLogger.WarnFormat("数据库同步异常 %s", errReturn.Error())
		return errReturn
	}
	_ = contentType
	_ = intervalTime
	//解析网关响应结果
	apiGatewayResp := contract.ApiGatewayResponse{}
	err := _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		shareDbSyncLogger.WarnFormat("数据库同步API网关响应结果解析异常 %s", err.Error())
		return err
	}
	if apiGatewayResp.RetCode != 0 {
		shareDbSyncLogger.WarnFormat("数据库同步API网关响应异常 %s",apiGatewayResp.RetMsg)
		return errors.New(apiGatewayResp.RetMsg)
	}
	//解析业务结果
	var result []SyncResult
	err = _json.Unmarshal(apiGatewayResp.Message, &result)
	if err != nil {
		shareDbSyncLogger.ErrorFormat(err, "数据库同步异常 反序列化异常  表%s 地址=%s  同步结果：%s", table, url, body)
		return err
	}

	if result == nil || len(result) == 0 || strings.Contains(result[0].Message, "Success") == false {
		shareDbSyncLogger.ErrorFormat(err, "数据库同步异常  表%s 地址=%s  同步结果：%s", table, url, body)
		return errors.New("数据库同步异常")
	}
	shareDbSyncLogger.DebugFormat("数据库同步成功  表%s 地址=%s  同步结果：%s", table, url, body)
	return nil
}

// Sync 数据库同步 如果失败了 waitSeconds秒后重试一次
func Sync(table string) {
	go func() {
		err := SyncCore(table)
		if err != nil {
			time.Sleep(time.Second * waitSeconds)
			SyncCore(table)
		}
	}()
}
