package encoding

import (
	"encoding/binary"
)

func calcBufferSize(postings *PostingsList) int {
	var pl *PostingsList
	count := 0
	for pl = postings; pl != nil; pl = pl.Next {
		count += 8 + len(pl.Positions)*4
	}
	return count
}

func EncodePostingsNone(postings *PostingsList, count int) []byte {
	var pl *PostingsList
	c := 0
	buf := make([]byte, calcBufferSize(postings))
	pos := 0
	for pl = postings; pl != nil; pl = pl.Next {
		binary.LittleEndian.PutUint32(buf[pos:], uint32(pl.DocumentId))
		pos += 4
		binary.LittleEndian.PutUint32(buf[pos:], uint32(len(pl.Positions)))
		pos += 4
		for _, v := range pl.Positions {
			binary.LittleEndian.PutUint32(buf[pos:], uint32(v))
			pos += 4
		}
		c++
	}
	if c != count {
		//fmt.Println("encoding miss", c, count)
	}
	return buf
}

func DecodePostingsNone(data []byte) (list *PostingsList, count int, err error) {
	var pl *PostingsList
	length := len(data)
	pos := 0
	var prev *PostingsList
	for pos < length {
		pl = new(PostingsList)
		if pos == 0 {
			list = pl
			prev = pl
		}
		count++

		pl.DocumentId = int(binary.LittleEndian.Uint32(data[pos:]))
		pos += 4
		length := int(binary.LittleEndian.Uint32(data[pos:]))
		pos += 4
		pl.Positions = make([]int, length)
		for i := 0; i < length; i++ {
			pl.Positions[i] = int(binary.LittleEndian.Uint32(data[pos:]))
			pos += 4
		}
		if prev != pl {
			prev.Next = pl
			prev = pl
		}

	}
	return
}
