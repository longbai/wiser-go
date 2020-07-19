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

func (p *PostingsList)Length() int {
	if p == nil {
		return 0
	}
	c := p
	count := 0
	for c != nil {
		count++
		c = c.Next
	}
	return count
}

func Merge(base, added *PostingsList) (ret *PostingsList) {
	var p *PostingsList
	/* 将二者连接成按文档编号升序排列的链表 */
	for base != nil || added != nil {
		e := new(PostingsList)
		if added == nil || (base != nil && base.DocumentId <= added.DocumentId) {
			*e = *base
			base = base.Next
		} else if base == nil || base.DocumentId >= added.DocumentId {
			*e = *added
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

func EncodePostings(postings *PostingsList, compress string, count int) []byte {
	switch compress {
	case "none":
		return EncodePostingsNone(postings, count)
	case "golomb":
		return EncodePostingsGolomb(postings, count)
	default:
		util.Abort()
	}
	return nil
}

func DecodePostings(data []byte, compress string) (list *PostingsList, count int, err error) {
	switch compress {
	case "none":
		return DecodePostingsNone(data)
	case "golomb":
		return DecodePostingsGolomb(data)
	default:
		util.Abort()
	}
	return
}