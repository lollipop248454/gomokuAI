package common

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/thinkeridea/go-extend/exstrings"
	"sort"
	"strings"
	"sync"
)

// 这个位置对自己的作用 我下在这能有多少分
// 不用考虑是否被0分割，因为会从头到尾遍历，靠近x,y的最终分也高
func solveFirst(ctx *Context, chess [][]int, x, y, player int) int {
	initScore := 0
	k := 1
	// 后手方下 主要考虑对对方的影响，进行围堵
	if player == 2 && ctx.FirstPlayer != "AI" {
		k = 3
	}
	if player == 1 && ctx.FirstPlayer == "AI" {
		k = 3
	}
	for l := 0; l < 4; l++ {
		line := "0"
		for i := 0; i < 4; i++ {
			xx := -dx[l]*(i+1) + x
			yy := -dy[l]*(i+1) + y
			s := "0"
			if out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", s, line)
		}
		for i := 0; i < 4; i++ {
			xx := dx[l]*(i+1) + x
			yy := dy[l]*(i+1) + y
			s := "0"
			if out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", line, s)
		}
		s := exstrings.SubString(line, 0, 5)
		count1 := strings.Count(s, "1")
		count2 := strings.Count(s, "2")
		for i := 0; i < 5; i++ {
			if count2 == 0 {
				initScore += FirstScore[count1] / k
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
	initScore := 0
	k := 1
	// 先手方下 不太会考虑对对方的影响 主要看自己的收益 前面用了3-，所以逻辑不变
	if player == 2 && ctx.FirstPlayer != "AI" {
		k = 3
	}
	if player == 1 && ctx.FirstPlayer == "AI" {
		k = 3
	}
	for l := 0; l < 4; l++ {
		line := "0"
		for i := 0; i < 4; i++ {
			xx := -dx[l]*(i+1) + x
			yy := -dy[l]*(i+1) + y
			s := "0"
			if out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", s, line)
		}
		for i := 0; i < 4; i++ {
			xx := dx[l]*(i+1) + x
			yy := dy[l]*(i+1) + y
			s := "0"
			if out(xx, yy) || chess[xx][yy] == 3-player {
				s = "2"
			} else if chess[xx][yy] == player {
				s = "1"
			}
			line = fmt.Sprintf("%s%s", line, s)
		}
		s := exstrings.SubString(line, 0, 5)
		count1 := strings.Count(s, "1")
		count2 := strings.Count(s, "2")
		for i := 0; i < 5; i++ {
			if count2 == 0 {
				initScore += SecondScore[count1] / k
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
	mx := make([][]int, 0)
	key := GenCacheKey(struct {
		Chess  [][]int `json:"chess"`
		Player int64   `json:"player"`
	}{
		Chess:  chess,
		Player: player,
	})
	if v, ok := cacheMap[key]; ok {
		//fmt.Println("命中了！")
		json.Unmarshal([]byte(cast.ToString(v)), &mx)
		return mx
	}
	mx = getMaxGreedyInfo(ctx, chess, player)
	body, _ := json.Marshal(mx)
	SetCache(key, body)
	return mx
}

func getMaxGreedyInfo(ctx *Context, chess [][]int, player int64) [][]int {
	mx := make([][]int, 0)
	var wt sync.WaitGroup
	ans := make(chan []int, 225)
	length := 0
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] > 0 || notRelative(i, j, 2, chess) {
				continue
			}
			length++
			// 这个是按照贪心平分
			wt.Add(1)
			go func(i, j int, player int64, chess *[][]int) {
				defer func() {
					wt.Done()
				}()
				v := solveFirst(ctx, *chess, i, j, int(player)) + solveSecond(ctx, *chess, i, j, 3-int(player))
				ans <- []int{v, i, j}
			}(i, j, player, &chess)
		}
	}
	wt.Wait()

	for i := 0; i < length; i++ {
		mx = append(mx, <-ans)
	}
	close(ans)
	sort.Slice(mx, func(i, j int) bool {
		return mx[i][0] < mx[j][0]
	})
	return mx
}
