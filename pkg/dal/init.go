package dal

var Dx []int
var Dy []int

var CntMap map[string]int

var NumMap map[string]int

var Param int

var ScoreMultiNum float64

var Score map[string]int64

//记录先手与后手的得分
var (
	SecondScore map[int]int
	FirstScore  map[int]int
)

func InitAI(id string) {
	CntMap[id] = 0
	NumMap[id] = 9
}

func init() {
	CntMap = make(map[string]int)
	NumMap = make(map[string]int)

	ScoreMultiNum = 1.5

	Param = 2

	Dx = []int{1, 1, 1, 0}
	Dy = []int{-1, 0, 1, 1}

	Score = make(map[string]int64)

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

	// Score

	Score["11111"] = 500000
	Score["011110"] = 4320
	Score["011100"] = 720
	Score["001110"] = 720
	Score["011010"] = 720
	Score["010110"] = 720
	Score["11110"] = 720
	Score["01111"] = 720
	Score["11011"] = 720
	Score["10111"] = 720
	Score["11101"] = 720
	Score["001100"] = 120
	Score["001010"] = 120
	Score["010100"] = 120
	//Score["001010"] = 80
	//Score["010100"] = 80
	Score["000100"] = 20
	Score["001000"] = 20
}

func InitChess() [][]int {
	chess := make([][]int, 15)
	for i := 0; i < 15; i++ {
		chess[i] = make([]int, 15)
	}
	return chess
}
