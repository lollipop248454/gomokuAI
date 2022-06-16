package common

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

var dx []int
var dy []int

var Cnt int

var score map[string]int64

//记录先手与后手的得分
var (
	SecondScore map[int]int
	FirstScore  map[int]int
)

func calScore(block, num int) int {
	if num >= 5 {
		return 100000
	}
	if block == 2 {
		return 0
	}
	c := 1
	if block == 0 {
		c *= 10
	}
	if num == 1 {
		return 1 * c
	}
	if num == 2 {
		return 10 * c
	}
	if num == 3 {
		return 100 * c
	}
	if num == 4 {
		return 1000 * c
	}
	return -1
}

// 无法拦截11101这种情况
func evalScore1(player int, chess [][]int) int64 {
	vis := make([][][]int, 15)
	for i := 0; i < 15; i++ {
		vis[i] = make([][]int, 15)
		for j := 0; j < 15; j++ {
			vis[i][j] = make([]int, 4)
		}
	}
	score := int64(0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != player {
				continue
			}
			for l := 0; l < 4; l++ {
				if vis[i][j][l] > 0 {
					continue
				}
				step := 0
				block := 0
				px := i - dx[l]
				py := j - dy[l]
				if out(px, py) || chess[px][py] == 3-player {
					block++
				}
				for true {
					step++
					xx := i + dx[l]*step
					yy := j + dy[l]*step
					if out(xx, yy) || chess[xx][yy] != player {
						if out(xx, yy) || chess[xx][yy] == 3-player {
							block++
						}
						break
					}
					vis[xx][yy][l] = 1
				}
				score += int64(calScore(block, step))
			}
		}
	}
	return score
}

func evalScorePara(player int, chess [][]int) int64 {
	finalScore := int64(0)
	var wt sync.WaitGroup
	ans := make(chan int64, 4)
	for l := 0; l < 4; l++ {
		wt.Add(1)
		go func(chess [][]int, l, player int) {
			defer func() {
				wt.Done()
			}()
			midScore := int64(0)
			for i := 0; i < 15; i++ {
				for j := 0; j < 15; j++ {
					if chess[i][j] == 3-player {
						continue
					}
					line := ""
					step := 0
					if chess[i][j] > 0 {
						line += "1"
					} else {
						line += "0"
					}
					for len(line) < 6 {
						step++
						x := dx[l]*step + i
						y := dy[l]*step + j
						if out(x, y) || chess[x][y] == 3-player {
							break
						}
						if chess[x][y] > 0 {
							line += "1"
						} else {
							line += "0"
						}
					}
					// 用sumscore or maxscore
					maxScore := int64(0)
					for len(line) >= 5 {
						if v, ok := score[line]; ok {
							maxScore = max(maxScore, v)
						}
						line = line[:(len(line) - 1)]
					}
					finalScore += maxScore
				}
			}
			ans <- midScore
		}(chess, l, player)
	}
	wt.Wait()
	for l := 0; l < 4; l++ {
		finalScore += <-ans
	}
	close(ans)
	return finalScore //+ findComb(player, chess)
}

// 不并行平均10s
func evalScore(player int, chess [][]int) int64 {
	finalScore := int64(0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] == 3-player {
				continue
			}
			for l := 0; l < 4; l++ {
				line := ""
				step := 0
				if chess[i][j] > 0 {
					line += "1"
				} else {
					line += "0"
				}
				for len(line) < 6 {
					step++
					x := dx[l]*step + i
					y := dy[l]*step + j
					if out(x, y) || chess[x][y] == 3-player {
						break
					}
					if chess[x][y] > 0 {
						line += "1"
					} else {
						line += "0"
					}
				}
				// 用sumscore or maxscore
				maxScore := int64(0)
				for len(line) >= 5 {
					if v, ok := score[line]; ok {
						maxScore = max(maxScore, v)
					}
					line = line[:(len(line) - 1)]
				}
				finalScore += maxScore
			}
		}
	}
	return finalScore //+ findComb(player, chess)
}

func eval(chess [][]int) int64 {
	return evalScore(2, chess) - evalScore(1, chess)
}

func AI(chess [][]int) (int, int) {
	t := time.Now()
	var px, py int
	v := int(ab(6, -100000000000000, 1000000000000, 2, 1, chess))
	py = v % 100
	px = v / 100
	//for d := int64(2); d <= 6; d += 2 {
	//	v := int(ab(d, -100000000000000, 1000000000000, 2, 1, chess))
	//	newScore := v / 10000
	//	v -= newScore * 10000
	//	py = v % 100
	//	px = v / 100
	//	// 有必胜方案直接退出
	//	if newScore >= 10000 {
	//		break
	//	}
	//}

	fmt.Println(evalScore(1, chess), evalScore(2, chess))
	elapsed := time.Since(t)
	fmt.Println("cost time:", elapsed)
	return px, py
}

var num int

func ab(depth, alpha, beta, player, firstLevel int64, chess [][]int) int64 {
	var px, py int
	if depth == 0 {
		return eval(chess)
	}
	mx := make([][]int, 0)
	var wt sync.WaitGroup
	ans := make(chan []int, 225)
	length := 0
	if player == 2 {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if chess[i][j] > 0 || notRelative(i, j, chess) {
					continue
				}
				length++
				// 这个是按照贪心平分
				wt.Add(1)
				go func(i, j int, player int64, chess [][]int) {
					defer func() {
						wt.Done()
					}()
					v := solveFirst(chess, i, j, int(player)) + solveSecond(chess, i, j, 3-int(player))
					ans <- []int{v, i, j}
				}(i, j, player, chess)
				//v := solveFirst(chess, i, j, int(player)) + solveSecond(chess, i, j, 3-int(player))
				//mx = append(mx, []int{v, i, j})
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
		down := len(mx) - num
		if down < 0 {
			down = 0
		}
		for idx := len(mx) - 1; idx >= down; idx-- {
			i := mx[idx][1]
			j := mx[idx][2]
			chess[i][j] = 2
			ret := int64(0)
			if Check(i, j, 2, chess) {
				ret = eval(chess)
			} else {
				newChess := deepCopy(chess)
				ret = ab(depth-1, alpha, beta, 3-player, firstLevel-1, newChess)
			}
			if firstLevel > 0 {
				fmt.Printf("AI层 位置,得分,alpha：%d %d %d %d\n", i, j, mx[idx][0], ret)
			}
			chess[i][j] = 0
			if ret > alpha {
				alpha = ret
				if firstLevel > 0 {
					px, py = i, j
				}
			}
			if beta <= alpha {
				break
			}
		}
		if firstLevel > 0 {
			return int64(px*100 + py)
		}
		return alpha
	} else {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if chess[i][j] > 0 || notRelative(i, j, chess) {
					continue
				}
				length++
				// 这个是按照贪心平分
				wt.Add(1)
				go func(i, j int, player int64, chess [][]int) {
					defer func() {
						wt.Done()
					}()
					v := solveFirst(chess, i, j, int(player)) + solveSecond(chess, i, j, 3-int(player))
					ans <- []int{v, i, j}
				}(i, j, player, chess)
				//v := solveFirst(chess, i, j, int(player)) + solveSecond(chess, i, j, 3-int(player))
				//mx = append(mx, []int{v, i, j})
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
		down := len(mx) - num
		if down < 0 {
			down = 0
		}
		for idx := len(mx) - 1; idx >= down; idx-- {
			i := mx[idx][1]
			j := mx[idx][2]
			chess[i][j] = 1
			ret := int64(0)
			if Check(i, j, 1, chess) {
				ret = eval(chess)
			} else {
				newChess := deepCopy(chess)
				ret = ab(depth-1, alpha, beta, 3-player, 0, newChess)
			}
			chess[i][j] = 0
			if firstLevel > 0 {
				fmt.Printf("玩家层 位置,得分,alpha：%d %d %d %d\n", i, j, mx[idx][0], ret)
			}
			beta = min(beta, ret)
			if beta <= alpha {
				break
			}
		}
		return beta
	}
}

var param int

func notRelative(x, y int, chess [][]int) bool {
	if Cnt <= 0 {
		param = 1
	} else {
		param = 2
	}
	for i := -param; i <= param; i++ {
		for j := -param; j <= param; j++ {
			if out(x+i, y+j) {
				continue
			}
			if chess[x+i][y+j] > 0 {
				return false
			}
		}
	}
	return true
}

func InitAI() {
	// 后续需要存储以保证并行
	Cnt = 0

	param = 2

	num = 10

	dx = []int{1, 1, 1, 0}
	dy = []int{-1, 0, 1, 1}

	score = make(map[string]int64)

	FirstScore = make(map[int]int)
	SecondScore = make(map[int]int)

	// first
	FirstScore[0] = 7
	FirstScore[1] = 35
	FirstScore[2] = 800
	FirstScore[3] = 15000
	FirstScore[4] = 800000

	//second
	SecondScore[0] = 0
	SecondScore[1] = 15
	SecondScore[2] = 400
	SecondScore[3] = 1800
	SecondScore[4] = 100000

	// score

	score["11111"] = 500000
	score["011110"] = 4320
	score["011100"] = 720
	score["001110"] = 720
	score["011010"] = 720
	score["010110"] = 720
	score["11110"] = 720
	score["01111"] = 720
	score["11011"] = 720
	score["10111"] = 720
	score["11101"] = 720
	score["001100"] = 120
	score["001010"] = 120
	score["010100"] = 120
	score["000100"] = 20
	score["001000"] = 20
}

func deepCopy(chess [][]int) [][]int {
	newChess := make([][]int, 15)
	for i := 0; i < 15; i++ {
		newChess[i] = make([]int, 15)
		copy(newChess[i], chess[i])
	}
	return newChess
}

func out(x, y int) bool {
	return x < 0 || x >= 15 || y < 0 || y >= 15
}

func InitChess() [][]int {
	chess := make([][]int, 15)
	for i := 0; i < 15; i++ {
		chess[i] = make([]int, 15)
	}
	return chess
}

func Check(x, y, k int, chess [][]int) bool {
	l := 0
	r := 0
	t := 0
	b := 0
	for i := x - 1; ; i-- {
		l = x - 1 - i
		if i < 0 || chess[i][y] != k {
			break
		}
	}
	for i := x + 1; ; i++ {
		r = i - x - 1
		if i >= 15 || chess[i][y] != k {
			break
		}
	}
	if l+r+1 >= 5 {
		return true
	}

	for j := y - 1; ; j-- {
		t = y - 1 - j
		if j < 0 || chess[x][j] != k {
			break
		}
	}
	for j := y + 1; ; j++ {
		b = j - y - 1
		if j >= 15 || chess[x][j] != k {
			break
		}
	}
	if t+b+1 >= 5 {
		return true
	}

	for l := -1; ; l-- {
		i := x + l
		j := y + l
		if i < 0 || j < 0 {
			b = -l - 1
			break
		}
		if chess[i][j] != k {
			b = -l - 1
			break
		}
	}
	for l := 1; ; l++ {
		i := x + l
		j := y + l
		if i >= 15 || j >= 15 {
			t = l - 1
			break
		}
		if chess[i][j] != k {
			t = l - 1
			break
		}
	}
	if t+b+1 >= 5 {
		return true
	}

	for l := -1; ; l-- {
		i := x + l
		j := y - l
		if i < 0 || j >= 15 {
			b = -l - 1
			break
		}
		if chess[i][j] != k {
			b = -l - 1
			break
		}
	}
	for l := 1; ; l++ {
		i := x + l
		j := y - l
		if i >= 15 || j < 0 {
			t = l - 1
			break
		}
		if chess[i][j] != k {
			t = l - 1
			break
		}
	}
	if t+b+1 >= 5 {
		return true
	}
	return false
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
