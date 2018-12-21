package dataapi

import (
	"encoding/json"
	"errors"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotlog"
	"io/ioutil"
	"net/http"
)

const (
	dataApiUtilServiceName = "DataApiUtilService"
)

var (
	shareDataApiUtilLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(dataApiUtilServiceName, func() {
		shareDataApiUtilLogger = dotlog.GetLogger(dataApiUtilServiceName)
	})
}

func GetDataApi(apiUrl string) (result [][]string, err error) {
	if len(apiUrl) == 0 {
		shareDataApiUtilLogger.Error(err, "dataapi请求地址不能为空")
		return result, errors.New("dataapi请求地址不能为空")
	}
	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		shareDataApiUtilLogger.ErrorFormat(err, "DataAPI http请求异常 URL=%s", apiUrl)
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("请求异常 StatusCode Not OK")
		shareDataApiUtilLogger.ErrorFormat(err, "DataAPI http请求异常 状态码为=%d URL=%s", resp.StatusCode, apiUrl)
		return result, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		shareDataApiUtilLogger.ErrorFormat(err, "DataAPI 读取body异常 URL=%s", apiUrl)
		return result, err
	}
	//这个接口是utf8带BOM头格式，需要手动去除BOM
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	respData := new(DataApiResponse)
	err = json.Unmarshal(data, respData)
	if err != nil {
		shareDataApiUtilLogger.ErrorFormat(err, "DataAPI Unmarshal异常 data=%s URL=%s", string(data), apiUrl)
		return result, err
	}
	if respData.Success && len(respData.Data) > 0 && respData.Data[0].Success && len(respData.Data[0].Result) > 2 {
		return respData.Data[0].Result, nil
	} else if respData.Success && len(respData.Data) > 0 && respData.Data[0].Success && len(respData.Data[0].Result) <= 2 {
		shareDataApiUtilLogger.DebugFormat("DataAPI查询结果为空 respData=%s", _json.GetJsonString(respData))
		return nil, err
	} else {
		err = errors.New("DataAPI查询异常")
		shareDataApiUtilLogger.ErrorFormat(err, "DataAPI查询异常 respData=%s", _json.GetJsonString(respData))
		return nil, err
	}
}
