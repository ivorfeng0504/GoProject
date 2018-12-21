package userproduct

import (
	"git.emoney.cn/softweb/roboadvisor/protected"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/dotweb/framework/json"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/util/cache"
)

const (
	userProductServiceName                             = "UserProductService"
	UserProductService_GetUserProductList_CacheKey     = "UserProductService:GetUserProductList:"
	UserProductService_GetUserProductList_CacheSeconds = 15 * 60
)

var (
	shareUserProductLogger dotlog.Logger
)

func init() {
	protected.RegisterServiceLoader(userProductServiceName, func() {
		shareUserProductLogger = dotlog.GetLogger(userProductServiceName)
	})
}

// GetUserProductList获取用户产品列表
func GetUserProductList(uid int64) (userProductList []*UserProduct, err error) {
	cacheKey := UserProductService_GetUserProductList_CacheKey + strconv.FormatInt(uid, 10)
	cacheProvider := _cache.GetRedisCacheProvider(protected.UserHomeRedisConfig)
	err = cacheProvider.GetJsonObj(cacheKey, &userProductList)
	if err == nil {
		return userProductList, nil
	}
	userProductList, err = getUserProductListCore(uid)
	if err == nil && userProductList != nil {
		err = cacheProvider.Set(cacheKey, jsonutil.GetJsonString(userProductList), UserProductService_GetUserProductList_CacheSeconds)
	}
	return userProductList, nil
}

// getUserProductListCore 获取用户产品列表-直接读取接口
func getUserProductListCore(uid int64) (userProductList []*UserProduct, err error) {
	pcProductList, err := GetUserPCProductList(uid)
	if err != nil {
		shareUserProductLogger.ErrorFormat(err, "获取用户PC端产品失败  uid=%s", uid)
		pcProductList = nil
	}
	mobileProductList, err := GetUserMobileProductList(uid)
	if err != nil {
		//shareUserProductLogger.ErrorFormat(err, "获取用户移动端产品失败  uid=%s", uid)
		mobileProductList = nil
	}
	if pcProductList != nil {
		userProductList = append(userProductList, pcProductList...)
	}
	if mobileProductList != nil {
		userProductList = append(userProductList, mobileProductList...)
	}
	return userProductList, nil
}
