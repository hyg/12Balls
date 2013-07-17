// balls project main.go
package main

import (
	"fmt"
	"time"
)

const m1 = 0x55555555 //binary: 0101...
const m2 = 0x33333333 //binary: 00110011..
const m4 = 0x0f0f0f0f //binary:  4 zeros,  4 ones ...
const m8 = 0x00ff00ff //binary:  8 zeros,  8 ones ...

func bitCount(x int) (n int) {
	x = (x & m1) + ((x >> 1) & m1) //put count of each  2 bits into those  2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each  4 bits into those  4 bits
	x = (x & m4) + ((x >> 4) & m4) //put count of each  8 bits into those  8 bits
	x = (x & m8) + ((x >> 8) & m8) //put count of each 16 bits into those 16 bits
	return x
}

type result struct {
	left    int
	right   int
	free    int
	ballcnt int
	mask0h  int //平衡时的偏重掩码
	mask0l  int //平衡时的偏轻掩码
	mask1h  int //左倾时的偏重掩码
	mask1l  int //左倾时的偏轻掩码
	mask2h  int //右倾时的偏重掩码
	mask2l  int //右倾时的偏轻掩码
}

type solution struct {
	left    int
	right   int
	left0   int
	right0  int
	left00  int
	right00 int
	left01  int
	right01 int
	left02  int
	right02 int
	left1   int
	right1  int
	left10  int
	right10 int
	left11  int
	right11 int
	left12  int
	right12 int
	left2   int
	right2  int
	left20  int
	right20 int
	left21  int
	right21 int
	left22  int
	right22 int
}

var resultmap = make(map[int]map[int]result)
var bitCnt [4096]int

func try(step int, maskh int, maskl int, s solution) (ret bool) {

	if bitCnt[maskh]+bitCnt[maskl] <= 1 {
		fmt.Printf("\nmaskh=%d,%012b\tmaskl=%d,%012b", maskh, maskh, maskl, maskl)
		return true
	}

	if step == 4 {
		return false
	}

	for left, rightmap := range resultmap {
		for right, themap := range rightmap {
			mask0h := maskh & themap.mask0h
			mask1h := maskh & themap.mask1h
			mask2h := maskh & themap.mask2h
			mask0l := maskl & themap.mask0l
			mask1l := maskl & themap.mask1l
			mask2l := maskl & themap.mask2l

			if (step == 1) && ((bitCnt[mask0h]+bitCnt[mask0l] > 9) || (bitCnt[mask1h]+bitCnt[mask1l] > 9) || (bitCnt[mask2h]+bitCnt[mask2l] > 9)) {
				continue
			}

			if (step == 2) && ((bitCnt[mask0h]+bitCnt[mask0l] > 3) || (bitCnt[mask1h]+bitCnt[mask1l] > 3) || (bitCnt[mask2h]+bitCnt[mask2l] > 3)) {
				continue
			}

			if (step == 3) && ((bitCnt[mask0h]+bitCnt[mask0l] > 1) || (bitCnt[mask1h]+bitCnt[mask1l] > 1) || (bitCnt[mask2h]+bitCnt[mask2l] > 1)) {
				continue
			}

			if try(step+1, mask0h, mask0l, s) {
				fmt.Printf("\nstep=%d\t平衡\tleft=%d,%012b\tright=%d,%012b", step, left, left, right, right)
			}
			if try(step+1, mask1h, mask1l, s) {
				fmt.Printf("\nstep=%d\t左倾\tleft=%d,%012b\tright=%d,%012b", step, left, left, right, right)
			}
			if try(step+1, mask2h, mask2l, s) {
				fmt.Printf("\nstep=%d\t右倾\tleft=%d,%012b\tright=%d,%012b", step, left, left, right, right)
			}
		}
	}
	return false
}

func main() {
	begin := time.Now()
	fmt.Println("\nbegin:", begin.String())

	//init the bit count
	for i := 0; i < 4096; i++ {
		bitCnt[i] = bitCount(i)
	}

	//init the result map
	for i := 0; i <= 4096; i++ {
		resultmap[i] = make(map[int]result)
	}

	//init result map
	CurLeftBit := 0
	CurRightBit := 0
	resultCnt := 0

	for left := 1; left < 4095; left++ {
		CurLeftBit = bitCnt[left]
		if CurLeftBit <= 6 {
			for right := left + 1; right < 4096; right++ {
				CurRightBit = bitCnt[right]
				if (CurLeftBit == CurRightBit) && (left&right == 0) {
					//free := (left & right) &^ 0xfff
					free := 0xfff - left - right
					resultmap[left][right] =
						result{
							left,
							right,
							free,
							CurLeftBit,
							free,
							free,
							left,
							right,
							right,
							left}

					resultCnt++
					//fmt.Printf("\nleft=%012b\tright=%012b\tfree1=%b\tfree2=%b", left, right, free, 0xfff-left-right)
					//fmt.Printf("No.%d\tresultmap[%012b][%012b]=%v\n", resultCnt, left, right, resultmap[left][right])

				}
			}
		}
	}

	var s solution
	try(1, 0xfff, 0xfff, s)

	fmt.Println("\nbegin:", begin.String(), "\nnow:", time.Now().String(), "\nused:", time.Since(begin))
}
