package repository

import (
	"github.com/devfeel/database/mysql"
	"github.com/devfeel/dotlog"
	"git.emoney.cn/softweb/roboadvisor/const"
)

type MySqlBaseRepository struct {
	mysql.MySqlDBContext
	databaseLogger dotlog.Logger
}

func (base *MySqlBaseRepository) InitLogger() {
	base.databaseLogger = dotlog.GetLogger(_const.LoggerName_Repository)
	base.DBCommand.OnTrace = base.Trace
	base.DBCommand.OnDebug = base.Debug
	base.DBCommand.OnInfo = base.Info
	base.DBCommand.OnWarn = base.Warn
	base.DBCommand.OnError = base.Error
}

func (base *MySqlBaseRepository) Trace(content interface{}) {
	if base.databaseLogger != nil {
		base.databaseLogger.Trace(content)
	}
}

func (base *MySqlBaseRepository) Debug(content interface{}) {
	if base.databaseLogger != nil {
		base.databaseLogger.Debug(content)
	}
}

func (base *MySqlBaseRepository) Info(content interface{}) {
	if base.databaseLogger != nil {
		base.databaseLogger.Info(content)
	}
}

func (base *MySqlBaseRepository) Warn(content interface{}) {
	if base.databaseLogger != nil {
		base.databaseLogger.Warn(content)
	}
}

func (base *MySqlBaseRepository) Error(err error, content interface{}) {
	if base.databaseLogger != nil {
		base.databaseLogger.Error(err, content)
	}
}