package _cache

import (
"git.emoney.cn/softweb/roboadvisor/config"
"github.com/devfeel/cache"
"github.com/devfeel/dotweb/session"
)

// GetRedisCacheProvider 获取Redis缓存组件
func GetRedisCacheProvider(configInfo *config.RedisInfo) cache.RedisCache {
	if configInfo == nil {
		return nil
	}
	cacheProvider := cache.GetRedisCachePoolConf(configInfo.ServerUrl, configInfo.MaxIdle, configInfo.MaxActive)
	cacheProvider.SetBackupServer(configInfo.BackupServer, configInfo.MaxIdle, configInfo.MaxActive)
	cacheProvider.SetReadOnlyServer(configInfo.ReadOnlyServer, configInfo.MaxIdle, configInfo.MaxActive)
	return cacheProvider
}

// GetSessionRedisConfig 获取用于SessionRedis使用的config
func GetSessionRedisConfig(configInfo *config.RedisInfo, sessionExpire int64, storeKeyPre string, cookieName string) *session.StoreConfig {
	if configInfo == nil {
		return nil
	}
	storeConfig := session.NewStoreConfig(session.SessionMode_Redis, sessionExpire, configInfo.ServerUrl, storeKeyPre)
	storeConfig.BackupServerUrl = configInfo.BackupServer
	if len(cookieName)>0{
		storeConfig.CookieName = cookieName
	}
	return storeConfig
}
