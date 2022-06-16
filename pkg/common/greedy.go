package common

import (
	"fmt"
	"github.com/thinkeridea/go-extend/exstrings"
	"strings"
)

func solveFirst(chess [][]int, x, y, player int) int {
	initScore := 0
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
		for i := 0; i < 5; i++ {
			s := exstrings.SubString(line, i, 5)
			if strings.Count(s, "2") > 0 {
				continue
			}
			initScore += FirstScore[strings.Count(s, "1")]
		}
	}
	return initScore
}

func solveSecond(chess [][]int, x, y, player int) int {
	initScore := 0
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
		for i := 0; i < 5; i++ {
			s := exstrings.SubString(line, i, 5)
			if strings.Count(s, "2") > 0 {
				continue
			}
			initScore += SecondScore[strings.Count(s, "1")]
		}
	}
	return initScore
}

func GreedyScore(chess [][]int, x, y, player int) int {
	return solveFirst(chess, x, y, 3-player) + solveSecond(chess, x, y, player)
}
