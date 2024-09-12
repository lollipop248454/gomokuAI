package service

import (
	"bytes"
	"fmt"
	"gomokuAI/pkg/dal"
	"gomokuAI/pkg/service/common"
	"gomokuAI/pkg/util"
	"sync"
	"time"
)

func evalScorePara(ctx *common.Context, player int, chess [][]int) int64 {
	finalScore := int64(0)
	var wt sync.WaitGroup
	ans := make(chan int64, 900)
	for l := 0; l < 4; l++ {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if chess[i][j] == 3-player || ctx.UnRelativeMap[i*100+j] {
					continue
				}
				wt.Add(1)
				go func(i, j, l int) {
					defer wt.Done()
					count := 1
					var buffer bytes.Buffer
					step := 0
					if chess[i][j] > 0 {
						buffer.WriteString("1")
					} else {
						buffer.WriteString("0")
					}
					for count < 6 {
						step++
						x := dal.Dx[l]*step + i
						y := dal.Dy[l]*step + j
						if util.Out(x, y) || chess[x][y] == 3-player {
							break
						}
						if chess[x][y] > 0 {
							buffer.WriteString("1")
						} else {
							buffer.WriteString("0")
						}
						count++
					}
					line := buffer.String()
					// 用sumscore or maxscore
					maxScore := int64(0)
					for len(line) >= 5 {
						if v, ok := dal.Score[line]; ok {
							maxScore = util.Max(maxScore, v)
						}
						line = line[:(len(line) - 1)]
					}
					ans <- maxScore
				}(i, j, l)

			}
		}
	}
	wt.Wait()
	close(ans)
	for score := range ans {
		finalScore += score
	}
	return finalScore //+ findComb(player, chess)
}

// 不并行平均10s
func evalScore(ctx *common.Context, player int, chess [][]int) int64 {
	finalScore := int64(0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if chess[i][j] == 3-player || ctx.UnRelativeMap[i*100+j] { //common.NotRelative(ctx, i, j, 2, chess)
				continue
			}
			for l := 0; l < 4; l++ {
				count := 1
				var buffer bytes.Buffer
				step := 0
				if chess[i][j] > 0 {
					buffer.WriteString("1")
				} else {
					buffer.WriteString("0")
				}
				for count < 6 {
					step++
					x := dal.Dx[l]*step + i
					y := dal.Dy[l]*step + j
					if util.Out(x, y) || chess[x][y] == 3-player {
						break
					}
					if chess[x][y] > 0 {
						buffer.WriteString("1")
					} else {
						buffer.WriteString("0")
					}
					count++
				}
				line := buffer.String()
				// 用sumscore or maxscore
				maxScore := int64(0)
				for len(line) >= 5 {
					if v, ok := dal.Score[line]; ok {
						maxScore = util.Max(maxScore, v)
					}
					line = line[:(len(line) - 1)]
				}
				finalScore += maxScore
			}
		}
	}
	return finalScore //+ findComb(player, chess)
}

func eval(ctx *common.Context, chess [][]int) int64 {
	defer func(tm time.Time) {
		util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	}(time.Now())
	aiScore := int64(0)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		aiScore = evalScore(ctx, 2, chess)
	}()
	playerScore := evalScore(ctx, 1, chess)
	wg.Wait()
	playerScore = int64(float64(playerScore) / dal.ScoreMultiNum)
	//if ctx.FirstPlayer == "AI" {
	//	playerScore = int64(float64(playerScore) / dal.ScoreMultiNum)
	//} else {
	//	aiScore = int64(float64(aiScore) / dal.ScoreMultiNum)
	//}
	return aiScore - playerScore
}

func getInitDeep(ctx *common.Context) int64 {
	cnt := dal.CntMap[ctx.ID]
	deep := int64(4)
	if cnt > 1 {
		deep = 6
	}
	if cnt > 3 {
		deep = 8
	}
	// 先手前期直接追杀
	if ctx.FirstPlayer == "AI" {
		deep++
		deep = util.Min(deep, 9)
	}
	if cnt > 5 {
		deep = 11
	}
	//if cnt > 8 {
	//	deep = 13
	//switch cnt % 2 {
	//case 0:
	//	deep = 11
	//case 1:
	//	deep = 9
	//case 2:
	//	deep = 8
	//}
	//}
	if cnt > 25 {
		deep = 9
	}
	//if cnt > 34 {
	//	deep = 7
	//}
	//if Cnt > 17 {
	//	switch Cnt % 2 {
	//	case 0:
	//		deep = 14
	//	case 1:
	//		deep = 12
	//		//case 2:
	//		//	deep = 8
	//	}
	//}
	return deep
}

func AI(ctx *common.Context, chess [][]int) (int, int) {
	dal.CntMap[ctx.ID]++
	util.ClearTimeMap()
	ctx.UnRelativeMap = common.GetUnRelativeMap(ctx, chess)
	t := time.Now()
	var px, py int
	deep := getInitDeep(ctx)
	dal.NumMap[ctx.ID] = 10
	fmt.Println("deep: ", deep, " cnt: ", dal.CntMap[ctx.ID])
	v := int(ab(ctx, deep, -100000000000000, 1000000000000, 2, 1, chess))
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
	util.ShowAllTime()
	fmt.Println(evalScore(ctx, 1, chess), evalScore(ctx, 2, chess))
	elapsed := time.Since(t)
	fmt.Println("cost time:", elapsed)
	return px, py
}

func abWithCache(ctx *common.Context, depth, alpha, beta, player, firstLevel int64, chess [][]int) int64 {
	// 并非总时间，会翻倍，主要看缓存命中率
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	// 第n层的这一步，获取上一步n+2的缓存结果
	// 读取缓存读2个，设置缓存也设置2个 有点麻烦，而且11 -> 9也实际作用只为了提升平均性能，对max没影响
	//key := util.GenCacheKey(struct {
	//	Chess  [][]int `json:"chess"`
	//	Depth  int64   `json:"depth"`
	//	Player int64   `json:"player"`
	//	//Alpha  int64   `json:"alpha"`
	//	//Beta   int64   `json:"beta"`
	//	Cnt int `json:"cnt"`
	//}{
	//	Chess:  chess,
	//	Depth:  depth,
	//	Player: player,
	//	//Alpha:  alpha,
	//	//Beta:   beta,
	//	Cnt: dal.CntMap[ctx.ID],
	//})
	//// 避免低层级获取缓存时读到评分
	//if firstLevel != 1 {
	//	cacheRes := cast.ToInt64(util.GetCache(key))
	//	if cacheRes != -1 {
	//		return cacheRes
	//	}
	//}
	res := ab(ctx, depth, alpha, beta, player, firstLevel, chess)
	//util.SetCache(key, res)
	return res
}

// https://oi-wiki.org/search/alpha-beta/
// 为什么能剪枝还得看图说话，即当前节点已经搜到其他节点比其他节点差的了，没必要了
// 可以试下第一层的5个节点之下的子节点广度降低，因为次数alpha,beta值都比较严格了，没必要遍历那么广了
// alpha,beta不同，其他相同的情况下，返回的估分也是不同的，先搜11层，在搜9层就必须把a,b给当做key
// 如果是自身的遍历时，访问到相同的情况，此时可忽略a,b的不同，从源头剪枝的策略了
func ab(ctx *common.Context, depth, alpha, beta, player, firstLevel int64, chess [][]int) int64 {
	// 并非总时间，会翻倍，主要看缓存命中率
	//defer func(tm time.Time) {
	//	util.AddTime(util.GetCurrentFuncName()+ctx.ID, time.Since(tm))
	//}(time.Now())
	tm := time.Now()
	_ = tm
	var px, py int
	if depth == 0 {
		res := eval(ctx, chess)
		return res
	}
	mx := common.GetMaxGreedyInfo(ctx, chess, player)
	if len(mx) == 0 {
		fmt.Println("mx有误！ ", chess, player)
		return -500000
	}
	down := len(mx) - dal.NumMap[ctx.ID]
	if down < 0 {
		down = 0
	}
	hasDowngrade := false
	// AI
	if player == 2 {
		// 直接5个
		if mx[len(mx)-1][0] >= 500000 {
			alpha = 500000
			// 相当于跳过循环
			down = len(mx)
		}
		if firstLevel == 1 && mx[len(mx)-1][0] >= 100000 {
			px, py = mx[len(mx)-1][1], mx[len(mx)-1][2]
			// 相当于跳过循环
			down = len(mx)
		}
		for idx := len(mx) - 1; idx >= down; idx-- {
			i := mx[idx][1]
			j := mx[idx][2]

			if time.Since(tm).Seconds() > 12 && firstLevel == 1 && hasDowngrade == false {
				dal.NumMap[ctx.ID] = 6
				depth = 9
				hasDowngrade = true
				//break
			}
			// 最多也就遍历9个
			//if time.Since(tm).Seconds() < 2 && down > 0 && len(mx)-down < 9 {
			//	down--
			//}
			chess[i][j] = 2
			originUnRelativeMap := ctx.UnRelativeMap
			common.UpdateUnRelativeMap(ctx.UnRelativeMap, i, j)
			ret := int64(0)
			if util.Check(i, j, 2, chess) {
				//ret = eval(chess)
				ret = 500000
			} else {
				ret = ab(ctx, depth-1, alpha, beta, 3-player, firstLevel-1, chess)
			}
			ctx.UnRelativeMap = originUnRelativeMap
			if firstLevel > 0 {
				fmt.Printf("AI层 位置,深度，广度，得分,ret：%d %d %d %d %d %d\n", i, j, depth, dal.NumMap[ctx.ID], mx[idx][0], ret)
				if len(mx)-1-idx > 5 && dal.NumMap[ctx.ID] > 7 {
					dal.NumMap[ctx.ID] = 7
				}
			}
			chess[i][j] = 0
			// 递归第一层判断估分最高的位置进行赋值，更深的递归return的都是估分，
			if ret > alpha {
				alpha = ret
				if firstLevel > 0 {
					px, py = i, j
				}
			}
			if beta <= alpha {
				break
			}
			// 必胜了
			if ret > 400000 {
				break
			}
			// 必须挡
			if mx[idx][0] >= 100000 {
				break
			}
		}
		if firstLevel > 0 {
			if alpha <= 1000 {
				dal.ScoreMultiNum = 1
			}
			if alpha > 1000 {
				dal.ScoreMultiNum = 1.5
			}
			if alpha > 2000 {
				dal.ScoreMultiNum = 2
			}
			if alpha > 4000 {
				dal.ScoreMultiNum = 2.5
			}
			fmt.Println("k: ", dal.ScoreMultiNum)
			res := int64(px*100 + py)
			return res
		}
		return alpha
	} else {
		if mx[len(mx)-1][0] >= 500000 {
			return -500000
		}
		for idx := len(mx) - 1; idx >= down; idx-- {
			if time.Since(tm).Seconds() > 2 && firstLevel == 0 && hasDowngrade == false {
				dal.NumMap[ctx.ID] = 6
				depth = 8
				hasDowngrade = true
				//break
			}
			i := mx[idx][1]
			j := mx[idx][2]
			chess[i][j] = 1
			originUnRelativeMap := ctx.UnRelativeMap
			common.UpdateUnRelativeMap(ctx.UnRelativeMap, i, j)
			ret := int64(0)
			if util.Check(i, j, 1, chess) {
				//ret = eval(chess)
				ret = -500000
			} else {
				ret = ab(ctx, depth-1, alpha, beta, 3-player, firstLevel-1, chess)
			}
			ctx.UnRelativeMap = originUnRelativeMap
			chess[i][j] = 0
			if firstLevel > 0 {
				fmt.Printf("玩家层 位置,得分,alpha：%d %d %d %d\n", i, j, mx[idx][0], ret)
			}
			// 在敌方最有的一系列可能得步骤中，选取使得我方AI得分最小的局面
			beta = util.Min(beta, ret)
			if beta <= alpha {
				break
			}
			// 必败了
			if ret < -400000 {
				break
			}
			// 必须挡
			if mx[idx][0] >= 100000 {
				break
			}
		}
		return beta
	}
}
