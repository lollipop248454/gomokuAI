package common

import (
	"fmt"
	"gomokuAI/pkg/dal"
	"gomokuAI/pkg/util"
)

// h -> 活
// m -> 眠

// H2Pos 活2 （能下出活3，下其他的没意义
func H2Pos(ctx *Context, chess [][]int, player int) [][]int {
	ans := make([][]int, 0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != 0 || NotRelative(ctx, i, j, 2, chess) {
				continue
			}
			if checkH3(chess, player, i, j) {
				ans = append(ans, []int{i, j})
			}
		}
	}
	return ans
}

// M3Pos 眠3能下的有效地方 (能下出眠4，下其他的没意义
func M3Pos(ctx *Context, chess [][]int, player int) [][]int {
	ans := make([][]int, 0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != 0 || NotRelative(ctx, i, j, 2, chess) { // ctx.UnRelativeMap[i*100+j] 后续能这么优化，担心有bug
				continue
			}
			if checkM4(chess, player, i, j) {
				ans = append(ans, []int{i, j})
			}
		}
	}
	return ans
}

// H3Pos 活3能下的有效地方 （判断标准：能下一些pos形成活4。 也只走能形成活4的这些pos, 比如有些pos能形成眠4，但是没意义
func H3Pos(ctx *Context, chess [][]int, player int) [][]int {
	ans := make([][]int, 0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != 0 || NotRelative(ctx, i, j, 1, chess) {
				continue
			}
			if checkH4(chess, player, i, j) {
				ans = append(ans, []int{i, j})
			}
		}
	}
	return ans
}

// Hm4Pos 活眠4能下的有效地方 下完能活5
func Hm4Pos(ctx *Context, chess [][]int, player int) [][]int {
	ans := make([][]int, 0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] != 0 || NotRelative(ctx, i, j, 1, chess) {
				continue
			}
			if checkHm5(chess, player, i, j) {
				ans = append(ans, []int{i, j})
			}
		}
	}
	return ans
}

// checkH3 下x,y能形成活3 判断010110 011100 边上2个必须是0，内部4个有3个1就行了 checked
func checkH3(chess [][]int, player, x, y int) bool {
	// 默认"1"为player自己
	for l := 0; l < 4; l++ {
		line := "1"
		// 因为1和最远的边界距离也就是4
		for i := 1; i <= 4; i++ {
			xx := x + i*dal.Dx[l]
			yy := y + i*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				line = fmt.Sprintf("%s%s", line, "1")
			} else if chess[xx][yy] == 0 {
				line = fmt.Sprintf("%s%s", line, "0")
			} else {
				break
			}
		}
		for i := 1; i <= 4; i++ {
			xx := x - i*dal.Dx[l]
			yy := y - i*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				line = fmt.Sprintf("%s%s", "1", line)
			} else if chess[xx][yy] == 0 {
				line = fmt.Sprintf("%s%s", "0", line)
			} else {
				break
			}
		}
		checkLen := 6
		if len(line) < checkLen {
			return false
		}
		count1 := util.CountChar(line, checkLen, '1')
		for i := 0; i <= len(line)-checkLen; i++ {
			// 边上都为0，内部有3个1
			if count1 == 3 && line[i] == '0' && line[i+checkLen-1] == '0' {
				return true
			}

			if i == len(line)-checkLen {
				break
			}
			if line[i] == '1' {
				count1--
			}
			if line[i+checkLen] == '1' {
				count1++
			}
		}
	}
	return false
}

// 下x,y能形成眠4 判断10111 11011 11110 checked
func checkM4(chess [][]int, player, x, y int) bool {
	// 默认1为player自己
	for l := 0; l < 4; l++ {
		line := "1"
		for i := 1; i <= 4; i++ {
			xx := x + i*dal.Dx[l]
			yy := y + i*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				line = fmt.Sprintf("%s%s", line, "1")
			} else if chess[xx][yy] == 0 {
				line = fmt.Sprintf("%s%s", line, "0")
			} else {
				break
			}
		}
		for i := 1; i <= 4; i++ {
			xx := x - i*dal.Dx[l]
			yy := y - i*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				line = fmt.Sprintf("%s%s", "1", line)
			} else if chess[xx][yy] == 0 {
				line = fmt.Sprintf("%s%s", "0", line)
			} else {
				break
			}
		}
		checkLen := 5
		if len(line) < checkLen {
			return false
		}
		count1 := util.CountChar(line, checkLen, '1')
		for i := 0; i <= len(line)-checkLen; i++ {
			// 边上都为0，内部有3个1
			if count1 == 3 && line[i] == '0' && line[i+checkLen-1] == '0' {
				return true
			}

			if i == len(line)-checkLen {
				break
			}
			if line[i] == '1' {
				count1--
			}
			if line[i+checkLen] == '1' {
				count1++
			}
		}
	}
	return false
}

// 下x,y能活4 判断011110 checked
func checkH4(chess [][]int, player, x, y int) bool {
	for l := 0; l < 4; l++ {
		count := 1
		idx := 0
		for {
			idx++
			xx := x + idx*dal.Dx[l]
			yy := y + idx*dal.Dy[l]
			if util.Out(xx, yy) {
				// 碰到2就不是活4了
				return false
			}
			if chess[xx][yy] == player {
				count++
			} else if chess[xx][yy] == 0 {
				// 边界需要为0
				count++
				break
			} else {
				return false
			}
		}
		idx = 0
		for {
			idx--
			xx := x + idx*dal.Dx[l]
			yy := y + idx*dal.Dy[l]
			if util.Out(xx, yy) {
				return false
			}
			if chess[xx][yy] == player {
				count++
			} else if chess[xx][yy] == 0 {
				// 边界需要为0
				count++
				break
			} else {
				return false
			}
		}
		if count >= 6 {
			return true
		}
	}
	return false
}

// 下x,y能5连 判断11111 checked
func checkHm5(chess [][]int, player, x, y int) bool {
	for l := 0; l < 4; l++ {
		count := 1
		idx := 0
		for {
			idx++
			xx := x + idx*dal.Dx[l]
			yy := y + idx*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				count++
			} else {
				break
			}
		}
		idx = 0
		for {
			idx--
			xx := x + idx*dal.Dx[l]
			yy := y + idx*dal.Dy[l]
			if util.Out(xx, yy) {
				break
			}
			if chess[xx][yy] == player {
				count++
			} else {
				break
			}
		}
		if count >= 5 {
			return true
		}
	}
	return false
}
