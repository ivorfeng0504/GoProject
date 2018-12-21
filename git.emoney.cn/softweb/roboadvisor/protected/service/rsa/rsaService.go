package rsa

import (
	"git.emoney.cn/softweb/roboadvisor/config"
	"net/http"
	"strings"
	"git.emoney.cn/softweb/roboadvisor/global"
	"io/ioutil"
	"strconv"
	"errors"
	"encoding/xml"
	"html/template"
)

func DecryptRSA(content string) (retMsg string, err error) {
	url := config.CurrentConfig.RsaDecryptUrl

	key := `-----BEGIN RSA PRIVATE KEY-----
BwIAAACkAABSU0EyAAQAAAEAAQBLhEbvl7OXxek1aewFQH12UVqoNU12Pu0I1wdmH96NHNN/iQqBEiWXjX7od6bHywD+ckz7urftHcRz4Qf4judh0jf5+xjpmxhsTnRd3zJQrw2Mw8CTc41DBXTq9Bg93Zxk7YO/caNJdr2SyTK7HtXTS7+PX83Du7j/Yq5TsjxIlqG4k8JND9s3dskrw2LUNOjtUoK6ikPvi7TWx9VHTzfde5IY1HKuwr/eOFaaDWBlc7mAEeIvXSnerfxx1Hv+Qc5ruZy5yju9V3UAEWmGvuUaFl24Cr6iJHRPJmWGvLT25WmUG8bhiI1tbIRSFsawbXnH2tXthuZQD2FAQCRVZ4a6AXKkwNxKxw9wOFW+5aGvYCeXP1QatiJxjAWcbUv9qkOHYTSerEvd3C3RUP0iv84Is+b3h8+umA18czobj0emOMEc2gskpYs5c9m9xRlY6SPsgWSfQ26tYbHDJl80Yutx1zhTFW3ohUGrKJMaGcJp55uq7HamAqAC1Qym18Ub34DQzFbGcAoucQC+3EIIU5AQoe9IYKF2zXGXllSeekZOginYHVeDmiOST2PbhAH/U2NA8q+GEtimg9vBKE4REfpxQR+2hOTvJxmpkejHcBHzncDvKlMTclXbBk8hwiYkeYcF2LQ2wdytFUGAMM7Nklw2s4S83oFsu3Hmt+h46GqOpS6H75/uopo4pzQebPKDg8fe0vI5DikGF+bumk6P08dj8pd7e+Yuf5BTqHDCPXA1YZhQGsPjlGugF88lszYMI2Q=
-----END RSA PRIVATE KEY-----`
	key = template.HTMLEscapeString(key)

	content = template.HTMLEscapeString(content)

	reqBody :=
		`<?xml version="1.0" encoding="utf-16"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			  <soap:Body>
				<Decrypt xmlns="http://tempuri.org/">
				  <key>` + key + `</key>
				  <content>` + content + `</content>
				</Decrypt>
			  </soap:Body>
		</soap:Envelope>`

	res, err := http.Post(url, "text/xml; charset=utf-8", strings.NewReader(reqBody))

	if nil != err {
		global.InnerLogger.Error(err, "RSA 解密异常 http post err:")
		return "", err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if nil != err {
		global.InnerLogger.ErrorFormat(err, "RSA 解密异常 ioutil ReadAll err  data=%s \n", string(data))
		return "", err
	}

	if http.StatusOK != res.StatusCode {
		global.InnerLogger.ErrorFormat(err, "RSAWebService soap1.1 request fail, status: %d  data=%s \n", res.StatusCode, string(data))
		return "", errors.New("RSA解析出错 状态码为" + strconv.Itoa(res.StatusCode))
	}

	RsaRet := RsaRet{}
	err = xml.Unmarshal(data, &RsaRet)

	if err != nil {
		global.InnerLogger.ErrorFormat(err, "RSA 解析异常 反序列化失败 data=%s", string(data))
		return "", err
	}
	retMsg = RsaRet.DecryptResponse[0].DecryptResult

	return retMsg, nil
}

//RSA解析结果
type RsaRet struct {
	XMLName           xml.Name           `xml:"Envelope"`
	DecryptResponse   []DecryptResponse `xml:"Body>DecryptResponse"`
}
type DecryptResponse struct {
	XMLName         xml.Name         `xml:"DecryptResponse"`
	DecryptResult   string
}