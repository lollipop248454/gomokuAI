package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

var cacheMap map[string]interface{}

func InitCache() {
	cacheMap = make(map[string]interface{})
}

func GenCacheKey(params interface{}) string {
	param, _ := json.Marshal(params)
	has := md5.Sum(param)

	return fmt.Sprintf("%x", has)
}

func SetCache(key string, value interface{}) {
	cacheMap[key] = value
}

func GetCache(key string) interface{} {
	if v, ok := cacheMap[key]; ok {
		//fmt.Println("命中了")
		return v
	}
	return -1
}
