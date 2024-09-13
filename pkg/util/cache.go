package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
)

var CacheMap map[string]interface{}

var CacheMapSync sync.Map

func init() {
	CacheMap = make(map[string]interface{})
	CacheMapSync = sync.Map{}
}

func GenCacheKey(params interface{}) string {
	param, _ := json.Marshal(params)
	has := md5.Sum(param)

	return fmt.Sprintf("%x", has)
}

func SetCache(key string, value interface{}) {
	CacheMap[key] = value
}

func SetCacheSync(key string, value interface{}) {
	CacheMapSync.Store(key, value)
}

func GetCache(key string) interface{} {
	if v, ok := CacheMap[key]; ok {
		//fmt.Println("命中了")
		return v
	}
	return nil
}

func GetCacheSync(key string) interface{} {
	if v, ok := CacheMapSync.Load(key); ok {
		return v
	}
	return nil
}
