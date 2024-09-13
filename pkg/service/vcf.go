package service

import (
	"fmt"
	"github.com/spf13/cast"
	"gomokuAI/pkg/dal"
	"gomokuAI/pkg/service/common"
	"gomokuAI/pkg/util"
)

type VCFInfo struct {
	HasVCF bool
	X      int
	Y      int
}

func CheckVcWithCache(ctx *common.Context, chess [][]int, player, depth int) int {
	originKey := util.GenCacheKey(struct {
		Chess  [][]int `json:"chess"`
		Player int     `json:"player"`
	}{
		Chess:  chess,
		Player: player,
	}) + "_vc_"
	tmpDepth := depth
	key := originKey + cast.ToString(tmpDepth)
	if v := util.GetCacheSync(key); v != nil {
		return cast.ToInt(v)
	}
	for {
		tmpDepth += 2
		if tmpDepth >= dal.MaxVcDepth {
			break
		}
		newKey := originKey + cast.ToString(tmpDepth)
		if v := util.GetCacheSync(newKey); v != nil {
			if v == 10 && (depth%2 == 0) && (depth >= 2) {
				return cast.ToInt(v)
			}
			if v == -10 && (depth%2 == 1) {
				return cast.ToInt(v)
			}
		}
	}

	v := CheckVc(ctx, chess, player, depth)
	if v == 100 || v == -100 {
		fmt.Printf("100有误了！！！！")
	}
	util.SetCacheSync(key, v)
	return v
}

// CheckVc depth初始值为0 只有杀棋才return pos，否则为-1
func CheckVc(ctx *common.Context, chess [][]int, player, depth int) int {
	if depth == dal.MaxVcDepth {
		// 没分出胜负，返回-1
		return -1
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
	maxVal := -100 // 负无穷
	// 对方4连先手必挡
	posList = common.Hm4Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
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
	//if depth == 0 {
	//	fmt.Printf("posH3: %+v,player: %+v",posList,player)
	//}
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
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
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
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
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
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
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
	// 对方眠3，可冲眠4，需要考虑对方是否能冲杀，可下棋，不强制return
	//posList = common.M3Pos(ctx, chess, 3-player)
	//if len(posList) > 0 {
	//	for _, pos := range posList {
	//		chess[pos[0]][pos[1]] = player
	//		ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
	//		chess[pos[0]][pos[1]] = 0
	//		if ret > maxVal {
	//			maxVal = ret
	//		}
	//		if maxVal == 10 {
	//			if depth == 0 {
	//				return pos[0]*100 + pos[1]
	//			} else {
	//				return 10
	//			}
	//		}
	//	}
	//}
	// 我方和对方都没啥好棋下，均势
	// 什么m3我不冲，h2我也不冲，可能冲了就是输，我就认为留下来就是均势
	return -1
}

// 后手层
func minLevel(ctx *common.Context, chess [][]int, player, depth int) int {
	posList := common.Hm4Pos(ctx, chess, player)
	// 我方4连先手下了就赢，不用递归进去check5连了
	if len(posList) > 0 {
		return -10
	}
	minVal := 100 // 类似正无穷
	// 对方4连先手必挡
	posList = common.Hm4Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				return minVal
			}
		}
		return minVal
	}
	// 我方活3必下 (也是必赢，对方没有4连)
	posList = common.H3Pos(ctx, chess, player)
	if len(posList) > 0 {
		return -10
	}
	// 可下棋，没有杀棋不强制return
	posList = common.M3Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				return minVal
			}
		}
	}
	// 对方活3必挡
	posList = common.H3Pos(ctx, chess, 3-player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				return minVal
			}
		}
		return minVal
	}
	// 我方活2，可冲活3，可下棋，不强制return
	posList = common.H2Pos(ctx, chess, player)
	if len(posList) > 0 {
		for _, pos := range posList {
			chess[pos[0]][pos[1]] = player
			ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
			chess[pos[0]][pos[1]] = 0
			if ret < minVal {
				minVal = ret
			}
			if minVal == -10 {
				return minVal
			}
		}
	}
	// 对方眠3，可冲眠4，需要考虑对方是否能冲杀，可下棋，不强制return
	//posList = common.M3Pos(ctx, chess, 3-player)
	//if len(posList) > 0 {
	//	for _, pos := range posList {
	//		chess[pos[0]][pos[1]] = player
	//		ret := CheckVcWithCache(ctx, chess, 3-player, depth+1)
	//		chess[pos[0]][pos[1]] = 0
	//		if ret < minVal {
	//			minVal = ret
	//		}
	//		if minVal == -10 {
	//			return minVal
	//		}
	//	}
	//}
	// 我方和对方都没啥好棋下，均势
	return -1
}
