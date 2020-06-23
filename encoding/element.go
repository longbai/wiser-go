package encoding

import (
	"github.com/longbai/wiser-go/util"
)

type PostingsList struct {
	DocumentId     int
	Positions      []int
	Next           *PostingsList
}

func (p *PostingsList)PositionsCount() int {
	return len(p.Positions)
}

func Merge(base, added *PostingsList) (ret *PostingsList) {
	var p *PostingsList
	/* 将二者连接成按文档编号升序排列的链表 */
	for base != nil || added != nil {
		var e *PostingsList
		if added == nil || (base != nil && base.DocumentId <= added.DocumentId) {
			e = base
		} else if base == nil || base.DocumentId >= added.DocumentId {
			e = added
			added = added.Next
		} else {
			util.Abort()
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

