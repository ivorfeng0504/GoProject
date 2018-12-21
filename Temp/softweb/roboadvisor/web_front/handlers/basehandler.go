package handlers

import "git.emoney.cn/softweb/roboadvisor/global"

type WebRuntimeCache struct{
	cachekeypre string
	cacheexpired int64
}

func NewWebRuntimeCache(runtimecachekeypre string,runtimecacheexpired int64) *WebRuntimeCache{
	var instance = new (WebRuntimeCache)
	instance.cachekeypre =  runtimecachekeypre
	instance.cacheexpired = runtimecacheexpired
	return instance
}

////RunTimeCache
func (inst *WebRuntimeCache)GetCache(key string)(value interface{}, exists bool) {
	key = inst.cachekeypre + key
	value, err := global.DotApp.Cache().Get(key)
	if err != nil || value == nil {
		return nil, false
	} else {
		return value, true
	}
}

func (inst *WebRuntimeCache)RemoveCache(key string) {
	key = inst.cachekeypre + key
	global.DotApp.Cache().Delete(key)
}

func (inst *WebRuntimeCache)SetCache(key string, value interface{}) error {
	key = inst.cachekeypre + key
	return global.DotApp.Cache().Set(key, value, inst.cacheexpired)
}
