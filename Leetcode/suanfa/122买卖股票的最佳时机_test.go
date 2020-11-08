package suanfa

import (
	"fmt"
	"testing"
)

var price =[]int{7,6,4,3,1}




func maxProfit(prices []int) int {
	if len(prices)<2{
		return 0
	}
	if len(prices)==2 && prices[0]>=prices[1]{
		return 0
	}

	if prices[0]<prices[1]{
		return prices[1]-prices[0]+maxProfit(prices[1:])
	}else {
		return maxProfit(prices[1:])
	}
}

func Test_suanfa2(t *testing.T) {
	//
	fmt.Println(fmt.Sprintf("最多可以赚 %d 元", 	maxProfit(price)))
}

//
//https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock-ii/
//给定一个数组，它的第 i 个元素是一支给定股票第 i 天的价格。
//
//设计一个算法来计算你所能获取的最大利润。你可以尽可能地完成更多的交易（多次买卖一支股票）。
//
//注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
//
// 
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock-ii
//著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。