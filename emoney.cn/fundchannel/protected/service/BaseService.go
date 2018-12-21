package service

import "github.com/devfeel/cache"

type BaseService struct {
	RedisCache cache.RedisCache
}
const (
	//开头标识需要与.NET项目中的KEY保持一致
	FundChannelKey      = "EMoney.FundChannel:"
)