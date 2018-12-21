package mobile

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"strconv"
	"strings"
)

const (
	EncryptMobileCacheKey = "mobile.EncryptMobileBase64."
)

// EncryptMobileBase64 将手机号加密后返回Base64格式
func EncryptMobileBase64(mobile string) (base64Str string, err error) {
	mobile = strings.Trim(mobile, " ")
	if mobile == "" {
		return mobile, nil
	}
	hex, err := EncryptMobileHex(mobile)
	if err != nil {
		return hex, err
	}
	base64Str = Hex2Base64(hex)
	return base64Str, nil
}

// EncryptMobileHex 将手机号加密后返回0x十六进制格式
func EncryptMobileHex(mobile string) (mobileHex string, err error) {
	mobile = strings.Trim(mobile, " ")
	if mobile == "" {
		return mobileHex, nil
	}
	cacheProvider := _cache.GetRedisCacheProvider(protected.LiveRedisConfig)
	cacheKey := EncryptMobileCacheKey + mobile
	err = cacheProvider.GetJsonObj(cacheKey, &mobileHex)
	if err == nil && len(mobileHex) > 0 {
		return mobileHex, err
	}
	url := config.CurrentConfig.EncryptMobileApi
	url = fmt.Sprintf(url, mobile)
	url = strings.Trim(url, " ")
	body, contentType, intervalTime, errReturn := _http.HttpGet(url)
	if errReturn != nil {
		return mobileHex, errReturn
	}
	_ = contentType
	_ = intervalTime
	apiGatewayResp := contract.ApiGatewayResponse{}
	err = _json.Unmarshal(body, &apiGatewayResp)
	if err != nil {
		return mobileHex, err
	}
	if apiGatewayResp.RetCode != 0 {
		return mobileHex, errors.New(apiGatewayResp.RetMsg)
	}
	data := MobileResponse{}
	err = _json.Unmarshal(apiGatewayResp.Message, &data)
	if data.RetCode == "1" {
		mobileHex = data.Message
		cacheProvider.SetJsonObj(cacheKey, mobileHex)
		return mobileHex, nil
	} else {
		err = errors.New(apiGatewayResp.Message)
		global.InnerLogger.ErrorFormat(err, "EncryptMobileHex 账号加密失败")
		return mobileHex, err
	}
}

//将十六进制格式的加密传转为Base64格式
func Hex2Base64(mobileHex string) (mobileBase64 string) {
	if mobileHex == "" {
		return mobileHex
	}
	head := mobileHex[0:2]
	if head != "0x" && head != "0X" {
		return mobileHex
	}
	mobileHex = mobileHex[2:]
	var bytes []byte
	for index, _ := range mobileHex {
		if index%2 == 0 {
			continue
		}
		hexStr := mobileHex[index-1 : index+1]
		base, _ := strconv.ParseInt(hexStr, 16, 10)
		bytes = append(bytes, byte(base))
	}
	mobileBase64 = base64.StdEncoding.EncodeToString(bytes)
	return mobileBase64
}

// Base642Hex Base64转十六进制
func Base642Hex(mobileBase64 string) (mobileHex string, err error) {
	if len(mobileBase64) == 0 {
		return mobileHex, nil
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(mobileBase64)
	if err != nil {
		return "", err
	}
	mobileHex = hex.EncodeToString(decodeBytes)
	mobileHex = "0x" + strings.ToUpper(mobileHex)
	return mobileHex, nil
}
