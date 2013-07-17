// balls project main.go
package main

import (
	"fmt"
)

// 0: 000
// 1: 001
// 2: 002
// 3: 010
// 4: 011
// 5: 012
// 6: 020
// 7: 021
// 8: 022
// 9: 100
//10: 101
//11: 102
//12: 110
//13: 111
//14: 112
//15: 120
//16: 121
//17: 122
//18: 200
//19: 201
//20: 202
//21: 210
//22: 211
//23: 212
//24: 220
//25: 221
//26: 222

func main() {
	SolutionNum := 1 //解决方案总数
	BallCnt := 12
	Position := [3]string{"空", "左", "右"}
	ResultHeavy := [3]string{"平", "左", "右"}
	ResultLight := [3]string{"平", "右", "左"}

	var BallPlaceNo [14]int
	//只使用1~13编号，表示各号球在三次称重时的位置编号，是0~26的整数。
	//位置编号是三进制的三位数，分别为三次称重的：0-闲置 1-左天平 2-右天平
	//如果要能测出每种情况，每个球的位置编号应该不重复。
	//且为了剔除重复方案，只记录各球位置编号是递增的方案。

	var PlaceBallNo [27]int
	//正数表示占用该位置编号的球编号，负数表示占用反位置的球编号的负数。

	Image := []int{0, 2, 1, 6, 8, 7, 3, 5, 4, 18, 20, 19, 24, 26, 25, 21, 23, 22, 9, 11, 10, 15, 17, 16, 12, 14, 13}
	//反位置，或轻球的{位置编号->结果编号}映射
	//重球在s1方案与轻球在s2将呈现相同的结果，这s1与s2是反位置。
	//在同一个方案中，互反位置只能有一个被使用。
	//结果编号的三位数分别为三次称重结果为：0-天平持平 1-天平左侧 2-天平右侧

	var Bit [3][27]int
	// 记录0~26的各位的三进制数字

	for i := 0; i < 27; i++ {
		Bit[0][i] = i % 3
		Bit[1][i] = ((i - Bit[0][i]) / 3) % 3
		Bit[2][i] = ((i - Bit[0][i] - Bit[1][i]*3) / 9) % 3

		PlaceBallNo[i] = 0
	}

	//遍历所有符合条件的方案。
	CurBallNo := 1
	CurPlaceNo := 1 //0方案无法识别轻重，直接放弃。
	for CurBallNo > 0 {
		for CurPlaceNo < 27 && PlaceBallNo[CurPlaceNo] != 0 {
			CurPlaceNo++
		}

		if CurPlaceNo < 27 {
			BallPlaceNo[CurBallNo] = CurPlaceNo
			PlaceBallNo[CurPlaceNo] = CurBallNo
			PlaceBallNo[Image[CurPlaceNo]] = -CurBallNo

			if CurBallNo < BallCnt {
				CurBallNo++
			} else {
				//所有球都有位置，方案完整。
				Balance := true
				PlaceCnt := []int{0, 0, 0} //0-闲置 1-左天平 2-右天平   的总球数

				PrintBuf := fmt.Sprintf("第%d套方案\r\n", SolutionNum)
				strBuf := [3]string{"", "", ""} //0-闲置 1-左天平 2-右天平   的输出字符串

				for turn := 0; turn <= 2; turn++ {
					strBuf[turn] = fmt.Sprintf("第%d次称重：", turn+1)
					PlaceCnt[0] = 0
					PlaceCnt[1] = 0
					PlaceCnt[2] = 0
					strBuf[0] = ""
					strBuf[1] = ""
					strBuf[2] = ""

					for ballno := 1; ballno <= BallCnt; ballno++ {
						PlaceCnt[Bit[turn][BallPlaceNo[ballno]]]++
						strBuf[Bit[turn][BallPlaceNo[ballno]]] += fmt.Sprintf("%d ", ballno)
					}

					if PlaceCnt[1] != PlaceCnt[2] {
						Balance = false
						break
					} else {
						PrintBuf += fmt.Sprintf("第%d次称重，左边：%s，右边：%s，空闲：%s.\r\n", turn, strBuf[1], strBuf[2], strBuf[0])
					}

				}

				if Balance {
					SolutionNum++

					for ballno := 1; ballno <= BallCnt; ballno++ {
						PrintBuf += fmt.Sprintf("%d号球位置：%s,%s,%s\t如果它偏重:%s,%s,%s\t如果它偏轻:%s,%s,%s\r\n",
							ballno,
							Position[Bit[0][BallPlaceNo[ballno]]],
							Position[Bit[1][BallPlaceNo[ballno]]],
							Position[Bit[2][BallPlaceNo[ballno]]],
							ResultHeavy[Bit[0][BallPlaceNo[ballno]]],
							ResultHeavy[Bit[1][BallPlaceNo[ballno]]],
							ResultHeavy[Bit[2][BallPlaceNo[ballno]]],
							ResultLight[Bit[0][BallPlaceNo[ballno]]],
							ResultLight[Bit[1][BallPlaceNo[ballno]]],
							ResultLight[Bit[2][BallPlaceNo[ballno]]])
					}

					fmt.Println(PrintBuf)
				}
				PlaceBallNo[CurPlaceNo] = 0
				PlaceBallNo[Image[CurPlaceNo]] = 0
				BallPlaceNo[CurBallNo] = 0
				CurPlaceNo++
			}
		} else {
			//CurPlaceNo == 27
			//当前球可选的位置编号已经用尽，回溯上一个球
			CurBallNo--

			if CurBallNo > 0 {
				//取上一个球的当前位置
				CurPlaceNo = BallPlaceNo[CurBallNo]
				//清除上一个球的位置
				PlaceBallNo[CurPlaceNo] = 0
				PlaceBallNo[Image[CurPlaceNo]] = 0
				BallPlaceNo[CurBallNo] = 0
				//位置编号加一，作为起点重新为上一个球找位置方案
				CurPlaceNo++
			}

		}

	}
}
