package sso

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"emoney.cn/fundchannel/config"
	"emoney.cn/fundchannel/global"
)

func main() {
	fmt.Printf("Hello Word!\n")
	DecryptSso("rand=21755&amp;Version=2016080801&amp;uid=1000479057&amp;pid=153000000&amp;sid=1239112&amp;tid=154&amp;agentid=200000000&amp;clienttype=12&amp;OutOfDate=0&amp;token=TpzfwmdhdxUxA0P7X0nmyh49GoTW2AefY919MP5zNO%2ffCm7eMpZhi2oLtYdTPVOQBBMiNF0pDHCAoUNOdBAA6BuWk7098hTF8agjOUcrkihPIsA5fHc3LxlTSxe82xbPVmEc91NBor9seDu9p1Xka02FUesvpTq4Pr5FcAtTIkQ%3d&amp;bata=0003", "127.0.0.1")

}

func DecryptSso(token, ip string) (Usr *UserInfoResult, err error) {
	//&需要替换为"&amp; 进行HtmlEncode
	token = template.HTMLEscapeString(token)
	url := config.CurrentConfig.SSOUrl

	reqBody :=
		`<?xml version="1.0" encoding="utf-16"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			<soap:Body>
				<GetUserInfo xmlns="http://tempuri.org/">
				<QueryString>` + token + `</QueryString>
				<UserIP>` + ip + `</UserIP>
				<rool>1</rool>
				</GetUserInfo>
			</soap:Body>
		</soap:Envelope>`

	res, err := http.Post(url, "text/xml; charset=UTF-8", strings.NewReader(reqBody))

	if nil != err {
		global.InnerLogger.Error(err, "SSO 解析异常 http post err:")
		return nil, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if nil != err {
		global.InnerLogger.ErrorFormat(err, "SSO 解析异常 ioutil ReadAll err url=%s data=%s \n", url, string(data))
		return nil, err
	}

	if http.StatusOK != res.StatusCode {
		global.InnerLogger.ErrorFormat(err, "WebService soap1.1 request fail, status: %d  url=%s data=%s \n", res.StatusCode, url, string(data))
		return nil, errors.New("SSO解析出错 状态码为" + strconv.Itoa(res.StatusCode))
	}

	SsoUsr := SsoRet{}

	err = xml.Unmarshal(data, &SsoUsr)

	if err != nil {
		global.InnerLogger.ErrorFormat(err, "SSO 解析异常 反序列化失败 data=%s", string(data))
		return nil, err
	}
	Usr = &SsoUsr.UserInfoResponses[0].UserInfoResults[0]
	accountInfo := strings.Split(Usr.SourceMsg, "_")
	if len(accountInfo) >= 4 {
		Usr.Account = accountInfo[1]
		Usr.Pwd = accountInfo[2]
	}
	if len(accountInfo) >= 5 {
		cidStr := accountInfo[4]
		Usr.Cid, err = strconv.ParseInt(cidStr, 10, 64)
		if err != nil {
			Usr.Cid = 0
		}
	}
	return Usr, nil
}

//SSO解析结果
type SsoRet struct {
	XMLName           xml.Name           `xml:"Envelope"`
	UserInfoResponses []UserInfoResponse `xml:"Body>GetUserInfoResponse"`
}
type UserInfoResponse struct {
	XMLName         xml.Name         `xml:"GetUserInfoResponse"`
	UserInfoResults []UserInfoResult `xml:"GetUserInfoResult"`
}
type UserInfoResult struct {
	Checked   bool   `xml:"Checked"`
	Uid       int64  `xml:"Uid"`
	Pid       int    `xml:"Pid`
	Msg       string `xml:"Msg"`
	SourceMsg string `xml:"SourceMsg"`
	Account   string
	Cid       int64
	Pwd		  string `xml:"Pwd"`
}
