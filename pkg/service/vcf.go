package service

import (
	"gomokuAI/pkg/dal"
	"gomokuAI/pkg/service/common"
)

type VCFInfo struct {
	HasVCF bool
	X      int
	Y      int
}

// CheckVc depth初始值为0 只有杀棋才return pos，否则为-1
func CheckVc(ctx *common.Context, chess [][]int, player, depth int) int {
	if depth == dal.MaxVcDepth {
		// 没分出胜负，返回0
		return 0
	}
	if depth%2 == 0 {
		return maxLevel(ctx, chess, player, depth)
	} else {
		return minLevel(ctx, chess, player, depth)
	}
}

// 先手层
func maxLevel(ctx *common.Context, chess [][]int, player, depth int) int {
	posList := common.Hm4Pos(ctx, chess, player)
	// 我方4连先手下了就赢，不用递归进去check5连了
	if len(posList) > 0 {
		pos := posList[0]
		if depth == 0 {
			return pos[0]*100 + pos[1]
		}
		return 10
	}
	maxVal := -10
	// 对方4连先手必挡
	posList = common.Hm4Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret > maxVal {
				maxVal = ret
			}
			if maxVal == 10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return 10
				}
			}
		}
		return maxVal
	}
	// 我方活3必下 (也是必赢，对方没有4连)
	posList = common.H3Pos(ctx, chess, player)
	if len(posList) > 0 {
		pos := posList[0]
		if depth == 0 {
			return pos[0]*100 + pos[1]
		}
		return 10
	}
	// 可下棋，没有杀棋不强制return
	posList = common.M3Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret > maxVal {
				maxVal = ret
			}
			if maxVal == 10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return 10
				}
			}
		}
	}
	// 对方活3必挡
	posList = common.H3Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret > maxVal {
				maxVal = ret
			}
			if maxVal == 10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return 10
				}
			}
		}
		return maxVal
	}
	// 我方活2，可冲活3，可下棋，不强制return
	posList = common.H2Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret > maxVal {
				maxVal = ret
			}
			if maxVal == 10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return 10
				}
			}
		}
	}
	return maxVal
}

// 后手层
func minLevel(ctx *common.Context, chess [][]int, player, depth int) int {
	posList := common.Hm4Pos(ctx, chess, player)
	// 我方4连先手下了就赢，不用递归进去check5连了
	if len(posList) > 0 {
		pos := posList[0]
		if depth == 0 {
			return pos[0]*100 + pos[1]
		}
		return -10
	}
	minVal := 10
	// 对方4连先手必挡
	posList = common.Hm4Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return minVal
				}
			}
		}
		return minVal
	}
	// 我方活3必下 (也是必赢，对方没有4连)
	posList = common.H3Pos(ctx, chess, player)
	if len(posList) > 0 {
		pos := posList[0]
		if depth == 0 {
			return pos[0]*100 + pos[1]
		}
		return -10
	}
	// 可下棋，没有杀棋不强制return
	posList = common.M3Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return minVal
				}
			}
		}
	}
	// 对方活3必挡
	posList = common.H3Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return minVal
				}
			}
		}
		return minVal
	}
	// 我方活2，可冲活3，可下棋，不强制return
	posList = common.H2Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = 1
			ret := CheckVc(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 1
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				if depth == 0 {
					return pos[0]*100 + pos[1]
				} else {
					return minVal
				}
			}
		}
	}
	return minVal
}
