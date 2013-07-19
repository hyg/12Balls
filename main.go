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

var resultmap = make(map[int]map[int]result)
var bitCnt [4096]int

func init() {

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
	return
}

type PossibleSet struct {
	level int
	seth  int
	setl  int
	child []step
}

type step struct {
	level   int
	left    int
	right   int
	outset0 PossibleSet
	outset1 PossibleSet
	outset2 PossibleSet
}

var ind = [5]string{"", "  ", "    ", "      ", "        "}
var setbitmax = []int{27, 9, 3, 1}

func findstep(level int, seth int, setl int) (bool, PossibleSet) {
	var ps PossibleSet

	ps.level = level
	ps.seth = seth
	ps.setl = setl
	ps.child = make([]step, 1)

	child := 1
	//fmt.Printf("\n%s:%d(%012b)--%d(%012b)", ind[level-1], seth, seth, setl, setl)

	var s step
	s.level = level

	for left, rightmap := range resultmap {
		for right, themap := range rightmap {
			//the next 3 branches set
			set0h := seth & themap.mask0h
			set1h := seth & themap.mask1h
			set2h := seth & themap.mask2h
			set0l := setl & themap.mask0l
			set1l := setl & themap.mask1l
			set2l := setl & themap.mask2l

			setbit0 := bitCnt[set0h] + bitCnt[set0l]
			setbit1 := bitCnt[set1h] + bitCnt[set1l]
			setbit2 := bitCnt[set2h] + bitCnt[set2l]

			if (setbit0 > setbitmax[level]) || (setbit1 > setbitmax[level]) || (setbit2 > setbitmax[level]) {
				continue
			}
			s.left = left
			s.right = right

			fmt.Printf("\n%s:%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)", ind[level-1], seth, seth, setl, setl, left, left, right, right, set0h, set0h, set0l, set0l)
			fmt.Printf("\n%s:%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)", ind[level-1], seth, seth, setl, setl, left, left, right, right, set1h, set1h, set1l, set1l)
			fmt.Printf("\n%s:%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)\t%d(%012b)--%d(%012b)", ind[level-1], seth, seth, setl, setl, left, left, right, right, set2h, set2h, set2l, set2l)

			if level == 3 {
				//叶子节点
				ps0 := PossibleSet{level + 1, set0h, set0l, nil}
				ps1 := PossibleSet{level + 1, set1h, set1l, nil}
				ps2 := PossibleSet{level + 1, set2h, set2l, nil}

				s.outset0 = ps0
				s.outset1 = ps1
				s.outset2 = ps2

				ps.child = append(ps.child, s)

				child++
				if child > 3 {
					return true, ps
				}

			} else {
				//递归
				r0, ps0 := findstep(level+1, set0h, set0l)
				r1, ps1 := findstep(level+1, set1h, set1l)
				r2, ps2 := findstep(level+1, set2h, set2l)

				if r0 && r1 && r2 {
					s.outset0 = ps0
					s.outset1 = ps1
					s.outset2 = ps2

					ps.child = append(ps.child, s)

					child++
					if child > 3 {
						return true, ps
					}
				}

			}

		}
	}

	if len(ps.child) > 1 {
		fmt.Println("\nnow:", time.Now().String())
		//fmt.Printf("\n%s:%d(%012b)--%d(%012b)", ind[level-1], ps.child[1].left, ps.child[1].left, ps.child[1].right, ps.child[1].right)
		return true, ps
	} else {
		return false, ps
	}
}

func main() {
	begin := time.Now()
	fmt.Println("\nbegin:", begin.String())

	//init()
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

	ret, ps := findstep(1, 0xfff, 0xfff)

	if ret {
		fmt.Println("\n%v", ps)
	}

	fmt.Println("\nbegin:", begin.String(), "\nnow:", time.Now().String(), "\nused:", time.Since(begin))
}
