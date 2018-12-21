package mystocksyn

import (
	"encoding/xml"
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/global"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	fmt.Printf("Hello Word!\n")

}

func CopyAndWrite(UidFrom string,UidTo string, UidOld string) (retFlag bool, err error) {
	//&需要替换为"&amp; 进行HtmlEncode
	url := config.CurrentConfig.MyStockSynURL

	reqBody :=
		`<?xml version="1.0" encoding="utf-16"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			<soap:Body>
				<CopyAndWrite xmlns="http://tempuri.org/">
				<UidFrom>` + UidFrom + `</UidFrom>
				<UidTo>` + UidTo + `</UidTo>
				<UidOld>` + UidOld + `</UidOld>
				</CopyAndWrite>
			</soap:Body>
		</soap:Envelope>`

	res, err := http.Post(url, "text/xml; charset=UTF-8", strings.NewReader(reqBody))

	if nil != err {
		global.InnerLogger.Error(err, "自选股云同步解析异常 http post err:"+err.Error())
		return false, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if nil != err {
		global.InnerLogger.ErrorFormat(err, "自选股云同步解析异常 ioutil ReadAll err url=%s data=%s \n", url, string(data))
		return false, err
	}

	if http.StatusOK != res.StatusCode {
		global.InnerLogger.ErrorFormat(err, "WebService soap1.1 request fail, status: %d  url=%s data=%s \n", res.StatusCode, url, string(data))
		return false, errors.New("自选股云同步解析出错 状态码为" + strconv.Itoa(res.StatusCode))
	}

	MyStockSynRet := MyStockSynRet{}
	err = xml.Unmarshal(data, &MyStockSynRet)

	if err != nil {
		global.InnerLogger.ErrorFormat(err, "自选股云同步解析异常 反序列化失败 data=%s", string(data))
		return false, err
	}

	global.InnerLogger.InfoFormat("调用云自选股同步成功 UidFrom=%s UidTo=%s UidOld=%s ", UidFrom, UidTo, UidOld)

	retFlag = MyStockSynRet.CopyAndWriteResponse[0].CopyAndWriteResult

	return retFlag, nil
}

//MyStockSyn解析结果
type MyStockSynRet struct {
	XMLName           xml.Name           `xml:"Envelope"`
	CopyAndWriteResponse   []CopyAndWriteResponse `xml:"Body>CopyAndWriteResponse"`
}
type CopyAndWriteResponse struct {
	XMLName         xml.Name         `xml:"CopyAndWriteResponse"`
	CopyAndWriteResult   bool
}
