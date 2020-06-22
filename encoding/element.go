package encoding

import "os"

type PostingsList struct {
	DocumentId     int
	Positions      []int
	Next           *PostingsList
}

func (p *PostingsList)PositionsCount() int {
	return len(p.Positions)
}

func Merge(pa, pb *PostingsList) (ret *PostingsList) {
	var p *PostingsList
	/* 用pa和pb分别遍历base和to_be_added（参见函数TokenIndex.Merge）中的倒排列表中的元素， */
	/* 将二者连接成按文档编号升序排列的链表 */
	for pa != nil || pb != nil {
		var e *PostingsList
		if pb == nil || (pa != nil && pa.DocumentId <= pb.DocumentId) {
			e = pa
		} else if pa == nil || pa.DocumentId >= pb.DocumentId {
			e = pb
			pb = pb.Next
		} else {
			os.Exit(0) // abort
		}
		e.Next = nil
		if ret == nil {
			ret = e
		} else {
			p.Next = e
		}
		p = e
	}
	return
}

