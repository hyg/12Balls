using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace _12balls
{
    class Program
    {
        static void Main(string[] args)
        {
            int i,j;
            int SolutionNum = 1 ;
            string[] Position = { "空","左","右"};
            string[] Result1 = { "平", "左", "右" };
            string[] Result2 = { "平", "右", "左" };
            int[] Image = { 0,2,1,6,8,7,3,5,4,18,20,19,24,26,25,21,23,22,9,11,10,15,17,16,12,14,13};
            int[] Place,Ball;
            int[,] B;
            int BallNum = 0;
            int PlaceNum = 1;
            bool Balance = true;
            string PrintBuffer;

            Place = new int[27];
            Ball = new int[12];
            B = new int[3, 27];

            for (i = 0; i < 27; i++)
            {
                B[0, i] = i % 3;
                B[1, i] = ((i - B[0, i]) / 3) % 3;
                B[2, i] = ((i - B[0, i] - B[1, i] * 3) / 9) % 3;
            }

            for (i = 0; i < 27; i++)
                Place[i] = 0;

            while (BallNum >= 0)
            {
                while (PlaceNum < 27 && Place[PlaceNum] != 0 )
                {
                    PlaceNum++;
                }
 
                if(PlaceNum < 27)
                {
                    //find a empty place
                    Ball[BallNum] = PlaceNum;

                    Place[PlaceNum] = 1;
                    Place[Image[PlaceNum]] = 2;

                    if (BallNum < 11)
                    {
                        // continue
                        BallNum++;
                    }
                    else
                    {
                        //the last ball, the solution finished!
                        Balance = true;
                        PrintBuffer = string.Empty;

                        for (i = 0; i < 3; i++)
                        {
                            string Left, Right, Free;
                            int LCnt, RCnt;

                            Left = string.Empty;
                            Right = string.Empty;
                            Free = string.Empty;

                            LCnt = 0;
                            RCnt = 0;

                            for (j = 0; j < 12; j++)
                            {

                                    switch (B[i, Ball[j]])
                                    {
                                        case 0:
                                            Free += (j+1).ToString() + " ";
                                            break;
                                        case 1:
                                            Left += (j+1).ToString() + " ";
                                            LCnt++;
                                            break;
                                        case 2:
                                            Right += (j+1).ToString() + " ";
                                            RCnt++;
                                            break;
                                    }
                            }

                            PrintBuffer += string.Format("第{0}次称重，左边：{1}，右边：{2}，空闲：{3}.\r\n", i+1, Left, Right, Free);


                            if (LCnt != RCnt)
                            {
                                Balance = false;
                            }

                        }

                        if (Balance)
                        {
                            Console.WriteLine("第{0}套方案：\r\n", SolutionNum++);
                            Console.WriteLine(PrintBuffer);

                            PrintBuffer = string.Empty;

                            for (int k = 0; k < 12; k++)
                            {
                                PrintBuffer += string.Format("{0}号球重：{1}{2}{3}.\r\n", k+1, Result1[B[0, Ball[k]]], Result1[B[1, Ball[k]]], Result1[B[2, Ball[k]]]);
                                PrintBuffer += string.Format("{0}号球轻：{1}{2}{3}.\r\n", k+1, Result2[B[0, Ball[k]]], Result2[B[1, Ball[k]]], Result2[B[2, Ball[k]]]);
                            }
                            Console.WriteLine(PrintBuffer);
                        }

                        Place[PlaceNum] = 0;
                        Place[Image[PlaceNum]] = 0;
                        PlaceNum++;
                    }
                }
                else
                {
                    BallNum--;
                    if (BallNum >= 0)
                    {
                        PlaceNum = Ball[BallNum];

                        Place[PlaceNum] = 0;
                        Place[Image[PlaceNum]] = 0;
                        PlaceNum++;

                        Ball[BallNum] = 0;
                    }

                }
            }



            Console.ReadKey();
            return ;
        }
    }
}
