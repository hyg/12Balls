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

var freemap = make(map[int]map[int]int)
var bitCnt [0x1000]int

func initdata() {
	//init the bit count
	for i := 0; i < 0x1000; i++ {
		bitCnt[i] = bitCount(i)
	}

	//init result map
	CurLeftBit := 0
	CurRightBit := 0
	freemapCnt := 0

	for left := 1; left < 0xfff; left++ {
		CurLeftBit = bitCnt[left]
		if CurLeftBit == 4 {
			freemap[left] = make(map[int]int)

			for right := left + 1; right < 0x1000; right++ {
				CurRightBit = bitCnt[right]
				if (CurLeftBit == CurRightBit) && (left&right == 0) {
					freemap[left][right] = 0xfff - left - right

					freemapCnt++
				}
			}
		}
	}

	fmt.Printf("freemapCnt:%d\r\n", freemapCnt)
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
	1: 1, 2: 2, 4: 3, 8: 4, 16: 5, 32: 6, 64: 7, 128: 8, 256: 9, 512: 10, 1024: 11, 2048: 12,
	0x1000: 13, 0x2000: 14, 0x4000: 15, 0x8000: 16, 0x10000: 17, 0x20000: 18, 0x40000: 19, 0x80000: 20,
	0x100000: 21, 0x200000: 22, 0x400000: 23, 0x800000: 24, 0x1000000: 25, 0x2000000: 26, 0x4000000: 27,
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

	max := setbitmax[level]

	for left, rightmap := range freemap {
		if left&(seth|setl) == 0 {
			//fmt.Printf("left:(%013b)\nseth:(%013b)\nsetl:(%013b)\n", left, seth, setl)
			continue
		}
		for right, free := range rightmap {
			/*if right&(seth|setl) == 0 {
				fmt.Printf("right:(%013b)\n seth:(%013b)\n setl:(%013b)\n", right, seth, setl)
				continue
			}*/

			//the next 3 branches set
			set0h := seth & free
			set0l := setl & free
			//set0 := bitCnt[set0h|set0l]
			set0 := bitCnt[set0h] + bitCnt[set0l]
			if set0 > max {
				continue
			}

			set1h := seth & left
			set1l := setl & right
			//set1 := bitCnt[set1h|set1l]
			set1 := bitCnt[set1h] + bitCnt[set1l]
			if set1 > max {
				continue
			}

			set2h := seth & right
			set2l := setl & left
			//set2 := bitCnt[set2h|set2l]
			set2 := bitCnt[set2h] + bitCnt[set2l]
			if set2 > max {
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
				if r0 {
					r1, ps1 := findstep(level+1, set1h, set1l)
					if r1 {
						r2, ps2 := findstep(level+1, set2h, set2l)
						if r2 {
							s.outset0 = ps0
							s.outset1 = ps1
							s.outset2 = ps2

							ps.child = append(ps.child, s)

							if level == 2 {
								return true, ps
							} /*else if level == 1 {
								child++
								if child == 1 {
									return true, ps
								}
							}*/
						}
					}

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
			fmt.Printf("%s:(%012b)-(%012b)\n", ind[1], step.left, step.right)
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
		fmt.Printf("%s%s:(%012b)-(%012b)\n", prefix, ind[ps.level], ps.child[0].left, ps.child[0].right)
		print(ps.child[0].outset0, prefix+">平")
		print(ps.child[0].outset1, prefix+">左")
		print(ps.child[0].outset2, prefix+">右")
	}
}

func main() {
	begin := time.Now()
	fmt.Println("\nbegin:", begin.String())

	initdata()

	ret, ps := findstep(1, 0xfff, 0xfff)

	if ret {
		fmt.Println("\nsearch completed.\nbegin:", begin.String(), "\nnow:", time.Now().String(), "\nused:", time.Since(begin))
		fmt.Println("\n\n一级方案总数：", len(ps.child))

		print(ps, "")
	}

	fmt.Println("\nbegin:", begin.String(), "\nnow:", time.Now().String(), "\nused:", time.Since(begin))
}
