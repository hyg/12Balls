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
var bitCnt [8192]int

func initdata() {
	//init the bit count
	for i := 0; i < 8192; i++ {
		bitCnt[i] = bitCount(i)
	}

	//init the result map
	for i := 0; i <= 8192; i++ {
		resultmap[i] = make(map[int]result)
	}

	//init result map
	CurLeftBit := 0
	CurRightBit := 0
	resultCnt := 0

	for left := 1; left < 8191; left++ {
		CurLeftBit = bitCnt[left]
		if CurLeftBit <= 6 {
			for right := left + 1; right < 8192; right++ {
				CurRightBit = bitCnt[right]
				if (CurLeftBit == CurRightBit) && (left&right == 0) {
					//free := (left & right) &^ 0xfff
					free := 0x1fff - left - right
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
				}
			}
		}
	}

	fmt.Printf("resultCnt:%d\r\n", resultCnt)
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

var ind = [5]string{"", "        ", "      ", "   ", ""}
var setbitmax = []int{27, 9, 3, 1}
var bitNo = map[int]int{
	1: 1, 2: 2, 4: 3, 8: 4, 16: 5, 32: 6, 64: 7, 128: 8, 256: 9, 512: 10, 1024: 11, 2048: 12, 4096: 13,
}

func findstep(level int, seth int, setl int) (bool, PossibleSet) {
	var ps PossibleSet

	ps.level = level
	ps.seth = seth
	ps.setl = setl
	ps.child = make([]step, 0)

	//child := 0

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

			set0 := bitCnt[set0h|set0l]
			set1 := bitCnt[set1h|set1l]
			set2 := bitCnt[set2h|set2l]

			if (set0 > setbitmax[level]) || (set1 > setbitmax[level]) || (set2 > setbitmax[level]) {
				continue
			}

			s.left = left
			s.right = right

			if level == 3 {
				//叶子节点
				ps0 := PossibleSet{level + 1, set0h, set0l, nil}
				ps1 := PossibleSet{level + 1, set1h, set1l, nil}
				ps2 := PossibleSet{level + 1, set2h, set2l, nil}

				s.outset0 = ps0
				s.outset1 = ps1
				s.outset2 = ps2

				ps.child = append(ps.child, s)

				return true, ps

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

					if level == 2 {
						return true, ps
					} /*else if level == 1 {
						child++
						if child == 3 {
							return true, ps
						}
					}*/
				}
			}

		}
	}

	if len(ps.child) > 1 {
		return true, ps
	}
	return false, ps
}

func print(ps PossibleSet, prefix string) {
	if ps.level == 1 {
		for psID, step := range ps.child {
			fmt.Printf("\n\n方案:%d\r\n", psID+1)
			fmt.Printf("%s:(%013b)-(%013b)\n", ind[1], step.left, step.right)
			print(step.outset0, "平")
			print(step.outset1, "左")
			print(step.outset2, "右")
		}
	} else if ps.level == 4 {
		if (ps.seth > 0) && (ps.seth == ps.setl) {
			fmt.Printf("%s:%d异常\n", prefix, bitNo[ps.seth])
		} else if ps.seth > 0 {
			fmt.Printf("%s:%d重\n", prefix, bitNo[ps.seth])
		} else if ps.setl > 0 {
			fmt.Printf("%s:%d轻\n", prefix, bitNo[ps.setl])
		} else {
			fmt.Printf("%s:不可能\n", prefix)
		}
	} else {
		fmt.Printf("%s%s:(%013b)-(%013b)\n", prefix, ind[ps.level], ps.child[0].left, ps.child[0].right)
		print(ps.child[0].outset0, prefix+">平")
		print(ps.child[0].outset1, prefix+">左")
		print(ps.child[0].outset2, prefix+">右")
	}
}

func main() {
	begin := time.Now()
	fmt.Println("\nbegin:", begin.String())

	initdata()

	ret, ps := findstep(1, 0x1fff, 0x1fff)

	if ret {
		//fmt.Println("\n%v", ps)
		fmt.Println("\n\n一级方案总数：", len(ps.child))

		print(ps, "")
	}

	fmt.Println("\nbegin:", begin.String(), "\nnow:", time.Now().String(), "\nused:", time.Since(begin))
}
