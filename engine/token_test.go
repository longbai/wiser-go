package engine

import (
	"fmt"
	"testing"
)

func TestTokenSplit(t *testing.T) {
	s :=`止损是指当某一投资出现的亏损达到预定数额时，及时斩仓出局，以避免形成更大的亏损。其目的就在于投资失误时把损失限定在较小的范围内。股票投资与赌博的一个重要区别就在于前者可通过止损把损失限制在一定的范围之内，同时又能够最大限度地获取成功的报酬，换言之，止损使得以较小代价博取较大利益成为可能。股市中无数血的事实表明，一次意外的投资错误足以致命，但止损能帮助投资 者化险为夷。

　　止损既是一种理念，也是一个计划，更是一项操作。止损理念是指投资者必须 从战略高度认识止损在股市投资中的重要意义，因为在高风险的股市中，首先是要生存下去，才谈得上进一步的发展，止损的关键作用就在于能让投资者更好的生存下来。可以说，止损是股市投资中最关键的理念之一。止损计划是指在一项重要的投资决策实施之前，必须相应地制定如何止损的计划，止损计划中最重要的一步是根据各种因素（如重要的技术位，或资金状况等）来决定具体的止损位。止损操作是止损计划的实施，是股市投资中具有重大意义的一个步骤，倘若止损计划不能化 为实实在在的止损操作，止损仍只是纸上谈兵。`
	biGramSplit(s, func(s string, i int) error {
		fmt.Println(s, i)
		return nil
	})
}
