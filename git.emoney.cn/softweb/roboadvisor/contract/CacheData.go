package contract

import "time"

type CacheData struct {

	//创建时间
	CreateTime time.Time

	//数据缓存过期时间（如果数据未过期，直接访问缓存，不请求API接口服务器）
	CacheExpireTime time.Time

	// 数据存储过期时间（如果数据存储未过期，会将缓存数据的Hash值发送给API服务器
	// 如果服务器的响应数据的Hash与请求中的一致，则只不会响应数据的内容(即ApiResponse.Message为空)）
	StoreExpireTime time.Time

	//数据
	Data ApiResponse
	}
