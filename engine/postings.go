package engine

import (
	"errors"
	"fmt"
	"github.com/longbai/wiser-go/db"
	"github.com/longbai/wiser-go/encoding"
	"github.com/longbai/wiser-go/util"
	"os"
)

type TokenIndexItems struct {
	DocCount int
	PositionsCount int
	Postings *encoding.PostingsList
}

func (p *TokenIndexItems) merge(other *TokenIndexItems) {
	p.DocCount += other.DocCount
	p.Postings = encoding.Merge(p.Postings, other.Postings)
}

type TokenIndex struct {
	index    map[int]*TokenIndexItems
	database *db.Db
	compress string
}

func NewTokenIndex(d *db.Db, compressMethod string) *TokenIndex {
	return &TokenIndex{
		index:    make(map[int]*TokenIndexItems),
		database: d,
		compress: compressMethod,
	}
}


/* 存储在缓冲区中的文档数量达到了指定的阈值时，更新存储器上的倒排索引 */
func (p *TokenIndex)Flush(threshold int) {
	if len(p.index) < threshold {
		return
	}
	util.PrintTimeDiff()

	for k, v := range p.index {
		p.updatePostings(k, v)
	}

	p.index = make(map[int]*TokenIndexItems)
	fmt.Println("index flushed")
	util.PrintTimeDiff()
}

/**
 * 从数据库中获取关联到指定词元上的倒排列表
 * @param[in] env 存储着应用程序运行环境的结构体
 * @param[in] token_id 词元编号
 * @param[out] postings 获取到的倒排列表
 * @param[out] postings_len 获取到的倒排列表中的元素数
 */
func (p *TokenIndex)fetchPostings(tokenId int) (pl *encoding.PostingsList, length int, err error) {
	count, postings, err := p.database.GetPostings(tokenId)
	if err != nil || count == 0 || len(postings) == 0 {
		return nil, 0, err
	}

	pl, length, err = p.decodePostings(postings)
	if err != nil {
		fmt.Println("postings list decode error", err)
		return
	}
	if count != count {
		err = errors.New(fmt.Sprintf("postings list decode error: stored:%d decoded:%d.", count, length))
		return nil, 0, err
	}
	return
}

func (p *TokenIndex)encodePostings(postings *encoding.PostingsList, count int) []byte {
	switch p.compress {
	case "":
		return encoding.EncodePostingsNone(postings)
	case "golomb":
		c, _ := p.database.GetDocumentCount()
		return encoding.EncodePostingsGolomb(postings, c)
	default:
		os.Exit(0)
	}
	return nil
}

func (p *TokenIndex)decodePostings(data []byte)(list *encoding.PostingsList, count int, err error)  {
	switch p.compress {
	case "":
		encoding.DecodePostingsNone(data)
	case "golomb":
		encoding.DecodePostingsGolomb(data)
	default:
		os.Exit(0)
	}
	return
}


/**
 * 将内存上（小倒排索引中）的倒排列表与存储器上的倒排列表合并后存储到数据库中
 * @param[in] env 存储着应用程序运行环境的结构体
 * @param[in] p 含有倒排列表的倒排索引中的索引项
 */
func (p *TokenIndex)updatePostings(tokenId int, items *TokenIndexItems) {
	oldPostings, length, err := p.fetchPostings(tokenId)
	if err != nil {
		fmt.Printf("cannot fetch old postings list of token(%d) for update.\n", tokenId)
		return
	}

	if length != 0 {
		encoding.Merge(items.Postings, oldPostings)
		items.DocCount += length
	}

	data := p.encodePostings(items.Postings, items.DocCount)

	p.database.UpdatePostings(tokenId, items.DocCount, data)
}

func (p *TokenIndex)Merge(other *TokenIndex) {
	if len(p.index) == 0 {
		p.index = other.index
		return
	}
	for k, v := range other.index {
		if v2, ok := p.index[k]; ok {
			v2.merge(v)
		} else {
			p.index[k] = v
		}
	}
}

