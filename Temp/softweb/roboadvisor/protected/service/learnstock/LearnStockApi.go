package learnstock

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/util/http"
)

// HasStrategyPowerByCid 投教课程学习进度是否满足50% from @qiaoli
func HasStrategyPowerByCid(cid int64) (finished bool, err error) {
	if cid <= 0 {
		return false, nil
	}
	apiUrl := "http://ds.emoney.cn/LearnStock/toujiao/HasStrategyPower_userhome?cid=%d"
	apiUrl = fmt.Sprintf(apiUrl, cid)
	result, _, intervalTime, err := _http.HttpGet(apiUrl)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "HasStrategyPowerByCid 投教课程学习进度是否满足百分之50  异常 url=%s result=%s intervalTime=%d", apiUrl, result, intervalTime)
		return finished, err
	}
	global.InnerLogger.DebugFormat("HasStrategyPowerByCid 投教课程学习进度是否满足百分之50 url=%s result=%s intervalTime=%d", apiUrl, result, intervalTime)
	if result == "1" {
		finished = true
	}
	return finished, err
}
