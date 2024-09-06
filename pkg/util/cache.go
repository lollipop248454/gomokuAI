package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

var CacheMap map[string]interface{}

func init() {
	CacheMap = make(map[string]interface{})
}

func GenCacheKey(params interface{}) string {
	param, _ := json.Marshal(params)
	has := md5.Sum(param)

	return fmt.Sprintf("%x", has)
}

func SetCache(key string, value interface{}) {
	CacheMap[key] = value
}

func GetCache(key string) interface{} {
	if v, ok := CacheMap[key]; ok {
		//fmt.Println("命中了")
		return v
	}
	return -1
}
