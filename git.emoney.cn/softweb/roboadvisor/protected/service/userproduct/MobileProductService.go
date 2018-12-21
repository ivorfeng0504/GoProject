package userproduct

import (
	"errors"
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/protected"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"time"
)

const (
	mobileProductServiceName = "MobileProductService"
)

var (
	shareMobileProductLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(mobileProductServiceName, func() {
		shareMobileProductLogger = dotlog.GetLogger(mobileProductServiceName)
	})
}

// GetUserMobileProductList 获取用户移动端产品列表
func GetUserMobileProductList(uid int64) (userProductList []*UserProduct, err error) {
	if uid <= 0 {
		return nil, errors.New("UID不正确")
	}
	apiUrl := config.CurrentConfig.MobileProductApi
	if len(apiUrl) == 0 {
		return nil, errors.New("未获取到移动产品API地址")
	}
	apiUrl = fmt.Sprintf(apiUrl, uid)
	body, contentType, intervalTime, err := _http.HttpGet(apiUrl)
	if err != nil {
		shareMobileProductLogger.ErrorFormat(err, "GetUserMobileProductList HTTP访问异常  请求地址为:%s   响应结果为:%s", apiUrl, body)
		return nil, err
	}
	_ = contentType
	_ = intervalTime
	apiResponse := new(MobileProductResponse)
	err = jsonutil.Unmarshal(body, &apiResponse)
	if err != nil {
		shareMobileProductLogger.ErrorFormat(err, "GetUserMobileProductList 响应结果反序列化异常  请求地址为:%s   响应结果为:%s", apiUrl, body)
		return nil, err
	}
	if apiResponse.Result.Code != 0 {
		shareMobileProductLogger.ErrorFormat(err, "GetUserMobileProductList 请求异常  请求地址为:%s   响应结果为:%s", apiUrl, body)
		return nil, err
	}
	if apiResponse.Detail == nil {
		return nil, nil
	}

	for _, product := range apiResponse.Detail {
		userProduct := &UserProduct{
			UID:         uid,
			ProductName: product.AuthName,
		}
		activateTime := time.Unix(product.CreateTime/1000, 0)
		expireTime := time.Unix(product.ExpiryTime/1000, 0)
		userProduct.ActivateTime = activateTime.Format("2006-01-02")
		userProduct.ExpireTime = expireTime.Format("2006-01-02")
		if expireTime.After(time.Now()) {
			userProduct.StateDesc = "未过期"
		} else {
			userProduct.StateDesc = "已过期"
		}
		userProductList = append(userProductList, userProduct)
	}
	return userProductList, nil
}
