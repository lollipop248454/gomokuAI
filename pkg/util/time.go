package util

import (
	"fmt"
	"time"
)

type timeInfo struct {
	CostTime time.Duration
	Count    int64
}

func add(t, data timeInfo) timeInfo {
	t.CostTime += data.CostTime
	t.Count += data.Count
	return t
}

var timeMap map[string]timeInfo

func init() {
	timeMap = make(map[string]timeInfo)
}

func AddTime(id string, tm time.Duration) {
	if _, ok := timeMap[id]; !ok {
		timeMap[id] = timeInfo{
			CostTime: 0,
			Count:    0,
		}
	}
	timeMap[id] = add(timeMap[id], timeInfo{
		CostTime: tm,
		Count:    1,
	})
}

func ShowAllTime() {
	for id, tm := range timeMap {
		fmt.Printf("%s total cost: %+vs, cnt: %+v, avgCost: %+v\n", id, tm.CostTime.Seconds(), tm.Count, tm.CostTime.Seconds()/float64(tm.Count))
	}
}

func ShowTime(id string) {
	tm := timeMap[id]
	fmt.Printf("%s total cost: %+vs, cnt: %+v, avgCost: %+v\n", id, tm.CostTime.Seconds(), tm.Count, tm.CostTime.Seconds()/float64(tm.Count))
}

func ClearTimeMap() {
	timeMap = make(map[string]timeInfo)
}
