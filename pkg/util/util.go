package util

import (
	"runtime"
	"sync"
)

// GetCurrentFuncName 获取正在运行的函数名
func GetCurrentFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// Check 走x,y能赢
func Check(x, y, k int, chess [][]int) bool {
	//l := 0
	//r := 0
	//t := 0
	//b := 0
	ans := make(chan bool, 4)
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		l := 0
		r := 0
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
			ans <- true
		}
	}()

	go func() {
		defer wg.Done()
		t := 0
		b := 0
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
			ans <- true
		}
	}()

	go func() {
		defer wg.Done()
		t := 0
		b := 0
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
			ans <- true
		}
	}()

	go func() {
		defer wg.Done()
		t := 0
		b := 0
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
			ans <- true
		}
	}()
	wg.Wait()
	close(ans)
	for j := range ans {
		if j {
			return true
		}
	}
	return false
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Out(x, y int) bool {
	return x < 0 || x >= 15 || y < 0 || y >= 15
}

func CountChar(s string, n int, c byte) int {
	ans := 0
	for i := 0; i < n; i++ {
		if s[i] == c {
			ans++
		}
	}
	return ans
}

func DeepCopyChess(chess [][]int) [][]int {
	data := make([][]int, 0)
	for i := 0; i < 15; i++ {
		tmp := make([]int, 0)
		for j := 0; j < 15; j++ {
			tmp = append(tmp, chess[i][j])
		}
		data = append(data, tmp)
	}
	return data
}
