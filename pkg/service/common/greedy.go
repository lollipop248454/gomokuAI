package common

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"gomokuAI/pkg/dal"
	"gomokuAI/pkg/util"
	"sort"
	"sync"
	"time"
)

func solve(ctx *Context, chess [][]int, x, y, player int) int {
	var wg sync.WaitGroup
	wg.Add(1)
	first := 0
	go func() {
		defer wg.Done()
		first = solveFirst(ctx, chess, x, y, player)
	}()
	second := solveSecond(ctx, chess, x, y, 3-player)
	wg.Wait()
	return first + second
}

// 这个位置对自己的作用 我下在这能有多少分
// 不用考虑是否被0分割，因为会从头到尾遍历，靠近x,y的最终分也高
func solveFirst(ctx *Context, chess [][]int, x, y, player int) int {
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	initScore := 0
	k := float64(1)
	if player == 1 {
		k = dal.ScoreMultiNum
	}
	for l := 0; l < 4; l++ {
		line := "0"
		for i := 0; i < 4; i++ {
			xx := -dal.Dx[l]*(i+1) + x
			yy := -dal.Dy[l]*(i+1) + y
			s := "0"
			if util.Out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", s, line)
		}
		for i := 0; i < 4; i++ {
			xx := dal.Dx[l]*(i+1) + x
			yy := dal.Dy[l]*(i+1) + y
			s := "0"
			if util.Out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", line, s)
		}
		//s := exstrings.SubString(line, 0, 5)
		count1 := util.CountChar(line, 5, '1')
		count2 := util.CountChar(line, 5, '2')
		for i := 0; i < 5; i++ {
			if count2 == 0 {
				initScore += int(float64(dal.FirstScore[count1]) / k)
			}

			if i == 4 {
				break
			}
			if line[i] == '1' {
				count1--
			}
			if line[i] == '2' {
				count2--
			}
			if line[i+5] == '1' {
				count1++
			}
			if line[i+5] == '2' {
				count2++
			}
		}
	}
	return initScore
}

// 堵上这个位置对对方的影响 对方下在这有多少分
func solveSecond(ctx *Context, chess [][]int, x, y, player int) int {
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	initScore := 0
	k := float64(1)
	if player == 1 {
		k = dal.ScoreMultiNum
	}
	for l := 0; l < 4; l++ {
		line := "0"
		for i := 0; i < 4; i++ {
			xx := -dal.Dx[l]*(i+1) + x
			yy := -dal.Dy[l]*(i+1) + y
			s := "0"
			if util.Out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", s, line)
		}
		for i := 0; i < 4; i++ {
			xx := dal.Dx[l]*(i+1) + x
			yy := dal.Dy[l]*(i+1) + y
			s := "0"
			if util.Out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", line, s)
		}
		//s := exstrings.SubString(line, 0, 5)
		//count1 := strings.Count(s, "1")
		count1 := util.CountChar(line, 5, '1')
		count2 := util.CountChar(line, 5, '2')
		for i := 0; i < 5; i++ {
			if count2 == 0 {
				initScore += int(float64(dal.SecondScore[count1]) / k)
			}

			if i == 4 {
				break
			}
			if line[i] == '1' {
				count1--
			}
			if line[i] == '2' {
				count2--
			}
			if line[i+5] == '1' {
				count1++
			}
			if line[i+5] == '2' {
				count2++
			}
		}
	}
	return initScore
}

func GetMaxGreedyInfo(ctx *Context, chess [][]int, player int64) [][]int {
	defer func(tm time.Time) {
		util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	}(time.Now())
	mx := make([][]int, 0)
	key := util.GenCacheKey(struct {
		Chess  [][]int `json:"chess"`
		Player int64   `json:"player"`
	}{
		Chess:  chess,
		Player: player,
	})
	if v, ok := util.CacheMap[key]; ok {
		//fmt.Println("命中了！")
		json.Unmarshal([]byte(cast.ToString(v)), &mx)
		return mx
	}
	mx = getMaxGreedyInfo(ctx, chess, player)
	body, _ := json.Marshal(mx)
	util.SetCache(key, body)
	return mx
}

func getMaxGreedyInfo(ctx *Context, chess [][]int, player int64) [][]int {
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	mx := make([][]int, 0)
	var wt sync.WaitGroup
	ans := make(chan []int, 225)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] > 0 || ctx.UnRelativeMap[i*100+j] { // NotRelative(ctx, i, j, 2, chess)
				continue
			}
			// 这个是按照贪心平分
			wt.Add(1)
			go func(i, j int, player int64, chess [][]int) {
				defer func() {
					wt.Done()
				}()
				v := solve(ctx, chess, i, j, int(player))
				ans <- []int{v, i, j}
			}(i, j, player, chess)
		}
	}
	wt.Wait()
	close(ans)

	for res := range ans {
		mx = append(mx, res)
	}
	sort.Slice(mx, func(i, j int) bool {
		return mx[i][0] < mx[j][0]
	})
	return mx
}

func NotRelative(ctx *Context, x, y, length int, chess [][]int) bool {
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	if dal.CntMap[ctx.ID] <= 2 {
		dal.Param = 1
	} else {
		dal.Param = length
	}
	for i := -dal.Param; i <= dal.Param; i++ {
		for j := -dal.Param; j <= dal.Param; j++ {
			if util.Out(x+i, y+j) {
				continue
			}
			if chess[x+i][y+j] > 0 {
				return false
			}
		}
	}
	return true
}

func GetUnRelativeMap(ctx *Context, chess [][]int) map[int]bool {
	res := make(map[int]bool)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			res[i*100+j] = false
			// 这个点为空并且边上没有
			if chess[i][j] == 0 && NotRelative(ctx, i, j, 2, chess) {
				res[i*100+j] = true
			}
		}
	}
	return res
}

func UpdateUnRelativeMap(unRelativeMap map[int]bool, x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			if util.Out(x+i, y+j) {
				continue
			}
			unRelativeMap[(x+i)*100+(y+j)] = false
		}
	}
}
