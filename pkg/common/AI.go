package common

import (
	"fmt"
	"sort"
)

var dx []int
var dy []int

var score map[string]int64
var oppoScore map[string]int64
var comb map[string]int

var AIFirst bool
var chess [][]int

func evalScore(player int) int64 {
	vis := make([][][]int, 15)
	for i := 0; i < 15; i++ {
		vis[i] = make([][]int, 15)
		for j := 0; j < 15; j++ {
			vis[i][j] = make([]int, 4)
		}
	}
	finalScore := int64(0)
	for i := -1; i < 15; i++ {
		for j := -1; j < 15; j++ {
			if i >= 0 && j >= 0 && chess[i][j] == player {
				continue
			}
			for l := 0; l < 4; l++ {
				if i >= 0 && j >= 0 && vis[i][j][l] > 0 {
					continue
				}
				if i >= 0 && j >= 0 {
					vis[i][j][l] = 1
				}
				line := ""
				step := 0
				if i < 0 || j < 0 || chess[i][j] > 0 {
					line += "2"
				} else {
					line += "0"
				}
				zeroCnt := 0
				for len(line) < 9 {
					step++
					x := dx[l]*step + i
					y := dy[l]*step + j
					if out(x, y) || chess[x][y] == 3-player {
						line += "2"
						break
					}
					if chess[x][y] == 0 {
						zeroCnt++
					}
					if zeroCnt >= 3 {
						break
					}
					if step == 1 && chess[x][y] == 0 {
						break
					}
					vis[x][y][l] = 1
					if chess[x][y] > 0 {
						line += "1"
					} else {
						line += "0"
					}
				}
				mxScore := int64(0)
				for len(line) >= 3 {
					lc := []rune(line)
					for i := 0; i < len(line)/2; i++ {
						lc[i], lc[len(line)-1-i] = lc[len(line)-1-i], lc[i]
					}
					if player == 1 {
						if v, ok := oppoScore[line]; ok {
							mxScore = max(v, mxScore)
						} else if v, ok := oppoScore[string(lc)]; ok {
							mxScore = max(v, mxScore)
						}
					}
					if v, ok := score[line]; ok {
						mxScore = max(v, mxScore)
					} else if v, ok := score[string(lc)]; ok {
						mxScore = max(v, mxScore)
					}
					line = line[:(len(line) - 1)]
				}
				finalScore += mxScore
			}
		}
	}
	return finalScore + findComb(player)
}

func eval() int64 {
	if AIFirst {
		return evalScore(2) - evalScore(1)*2
	}
	return evalScore(2) - evalScore(1)*3
}

var px, py int

func AI(newChess [][]int) (int, int) {
	chess = newChess
	ab(5, -1000000000000000000, 1000000000000000, 2, 1)
	fmt.Println(evalScore(1), evalScore(2))
	return px, py
}

var num int

func ab(depth, alpha, beta, player, firstLevel int64) int64 {
	if depth == 0 {
		return eval()
	}
	mx := make([][]int64, 0)
	if player == 2 {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if chess[i][j] > 0 || notRelative(i, j) {
					continue
				}
				chess[i][j] = 2
				mx = append(mx, []int64{eval(), int64(i), int64(j)})
				chess[i][j] = 0
			}
		}
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
			if Check(int(i), int(j), 2) {
				ret = eval()
			} else {
				ret = ab(depth-1, alpha, beta, 3-player, 0)
			}
			chess[i][j] = 0
			if ret > alpha {
				alpha = ret
				if firstLevel > 0 {
					px = int(i)
					py = int(j)
				}
			}
			if beta <= alpha {
				break
			}
		}
		return alpha
	} else {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if chess[i][j] > 0 || notRelative(i, j) {
					continue
				}
				chess[i][j] = 1
				mx = append(mx, []int64{eval(), int64(i), int64(j)})
				chess[i][j] = 0
			}
		}
		sort.Slice(mx, func(i, j int) bool {
			return mx[i][0] > mx[j][0]
		})
		down := len(mx) - num
		if down < 0 {
			down = 0
		}
		for idx := len(mx) - 1; idx >= down; idx-- {
			i := mx[idx][1]
			j := mx[idx][2]
			chess[i][j] = 1
			if Check(int(i), int(j), 1) {
				beta = min(beta, eval())
			} else {
				beta = min(beta, ab(depth-1, alpha, beta, 3-player, 0))
			}
			chess[i][j] = 0
			if beta <= alpha {
				break
			}
		}
		return beta
	}
}

func findComb(player int) int64 {
	score := 0
	vis := make([][]int, 15)
	for i := 0; i < 15; i++ {
		vis[i] = make([]int, 15)
	}
	for i := -1; i < 15; i++ {
		for j := -1; j < 15; j++ {
			if i >= 0 && j >= 0 && chess[i][j] == player {
				continue
			}
			for l := 0; l < 4; l++ {
				line := ""
				step := 0
				if i < 0 || j < 0 || chess[i][j] > 0 {
					line += "2"
				} else {
					line += "0"
				}
				zeroCnt := 0
				for len(line) < 8 {
					step++
					x := dx[l]*step + i
					y := dy[l]*step + j
					if out(x, y) || chess[x][y] == 3-player {
						line += "2"
						break
					}
					if chess[x][y] == 0 {
						zeroCnt++
					}
					if zeroCnt >= 3 {
						break
					}
					if step == 1 && chess[x][y] == 0 {
						break
					}
					if chess[x][y] > 0 {
						line += "1"
					} else {
						line += "0"
					}
				}
				for len(line) >= 5 {
					lc := []rune(line)
					for i := 0; i < len(line)/2; i++ {
						lc[i], lc[len(line)-1-i] = lc[len(line)-1-i], lc[i]
					}
					f := 0
					if _, ok := comb[line]; ok {
						for idx := 1; idx < len(line)-1; idx++ {
							x := i + idx*dx[l]
							y := j + idx*dy[l]
							if chess[x][y] == player {
								vis[x][y]++
							}
						}
						f = 1
					} else if _, ok := comb[string(lc)]; ok {
						for idx := 1; idx < len(string(lc))-1; idx++ {
							x := i + idx*dx[l]
							y := j + idx*dy[l]
							if chess[x][y] == player {
								vis[x][y]++
							}
						}
						f = 1
					}
					if f == 1 {
						break
					}
					line = line[:(len(line) - 1)]
				}
			}
		}
	}
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if vis[i][j] >= 2 {
				score += 1000 * (1 << vis[i][j])
			}
		}
	}
	return int64(score)
}

func findComb1(player int) int64 {
	score := 0
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != player {
				continue
			}
			c3 := 0
			c4 := 0
			for l := 0; l < 4; l++ {
				sep1 := 0
				sep2 := 0
				num1 := 0
				num2 := 0
				step := 0
				for true {
					step++
					x := i + dx[l]*step
					y := j + dy[l]*step
					if out(x, y) || chess[x][y] == 3-player {
						sep1 = 1
						sep2 = 1
						break
					}
					if step > 1 && chess[x][y] == 0 {
						break
					}
					if step == 1 && chess[x][y] == 0 {
						sep1 = 1
						num1--
					}
					num1++
				}
				if sep1*sep2 == 1 {
					continue
				}
				step = 0
				for true {
					step++
					x := i - dx[l]*step
					y := j - dy[l]*step
					if out(x, y) || chess[x][y] == 3-player {
						sep1 = 1
						sep2 = 1
						break
					}
					// 此处考虑不周，需要更改
					if step > 1 && chess[x][y] == 0 {
						break
					}
					if step == 1 && chess[x][y] == 0 {
						sep2 = 1
						num2--
					}
					num2++
				}
				if sep1*sep2 == 1 {
					continue
				}
				if num1+num2 == 2 {
					c3++
				}
				if num1+num2 == 3 {
					c4++
					if sep1+sep2 == 1 {
						if sep1 == 1 && num1 == 2 {
							c4--
						}
						if sep2 == 1 && num2 == 2 {
							c4--
						}
					}
				}
			}
			if c3+c4 >= 2 {
				v := 1000 * (1 << c3) * (1 << c4)
				score += v + 500*c4
			}
		}
	}
	return int64(score)
}

var param int

func notRelative(x, y int) bool {
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
	param = 2

	num = 5

	dx = []int{1, 1, 1, 0}
	dy = []int{-1, 0, 1, 1}

	score = make(map[string]int64)

	// 这个棋型在对面分要高一点，我就需要进行堵
	oppoScore = make(map[string]int64)

	comb = make(map[string]int)
	// 1
	score["010"] = 10

	// 2
	score["0110"] = 100
	score["2110"] = 20
	score["01010"] = 50
	score["21010"] = 10

	//3

	score["01110"] = 1000
	score["21110"] = 200
	score["011010"] = 800
	score["211010"] = 200
	score["210110"] = 200
	score["0110010"] = 300
	score["2110010"] = 30
	score["2100110"] = 200
	score["0101010"] = 500
	score["2101010"] = 400

	//4
	score["011110"] = 100000
	score["211110"] = 6000
	score["0111010"] = 6000
	score["2111010"] = 3000
	score["2101110"] = 6000
	score["2101112"] = 3000
	score["0110110"] = 6000
	score["2110110"] = 2000

	// 5
	score["0111110"] = 999999999
	score["2111110"] = 999999999
	score["2111112"] = 999999999
	score["01011112"] = 8000
	score["21011112"] = 8000
	score["01101110"] = 9000
	score["21101110"] = 9000
	score["01101112"] = 8500
	score["21101112"] = 8000

	//6
	score["011101110"] = 8000
	score["211101110"] = 7000
	score["211101112"] = 6000
	score["011011112"] = 9000
	score["211011112"] = 8000

	// oppoScore
	oppoScore["01110"] = 5000
	oppoScore["010110"] = 5000
	oppoScore["011010"] = 5000
	oppoScore["0111010"] = 60000
	oppoScore["2111010"] = 50000
	oppoScore["2101110"] = 60000
	oppoScore["2101112"] = 50000
	oppoScore["0110110"] = 40000
	oppoScore["2110110"] = 35000
	oppoScore["21101102"] = 35000
	oppoScore["011011112"] = 60000
	oppoScore["211011112"] = 60000
	oppoScore["011101110"] = 80000
	oppoScore["211101110"] = 70000
	oppoScore["211101112"] = 60000

	//comb
	// 3
	comb["01110"] = 0
	comb["011010"] = 0
	// 4
	comb["0101110"] = 1
	comb["0101112"] = 1
	comb["2101110"] = 1
	comb["2101112"] = 1
	comb["0110110"] = 1
	comb["2110110"] = 1
	comb["2110112"] = 1
	// 5
	comb["01101110"] = 1
	comb["21101110"] = 1
	comb["01101112"] = 1
	comb["21101112"] = 1
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

func Check(x, y, k int) bool {
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
