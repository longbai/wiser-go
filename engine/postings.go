package engine

import (
	"errors"
	"fmt"
	"github.com/longbai/wiser-go/db"
	"github.com/longbai/wiser-go/encoding"
	"github.com/longbai/wiser-go/util"
)

type PostingsManager struct {
	index    map[int]*tokenIndexItems
	database *db.SqliteDb
	compress string
}

func NewPostingsManager(d *db.SqliteDb, compressMethod string) *PostingsManager {
	return &PostingsManager{
		index:    make(map[int]*tokenIndexItems),
		database: d,
		compress: compressMethod,
	}
}

/* 存储在缓冲区中的文档数量达到了指定的阈值时，更新存储器上的倒排索引 */
func (p *PostingsManager) Flush(threshold int) {
	l := len(p.index)
	if l <= threshold {
		return
	}
	util.PrintTimeDiff()

	for k, v := range p.index {
		p.updatePostings(k, v)
	}

	p.index = make(map[int]*tokenIndexItems)
	fmt.Println("index flushed", l)
	util.PrintTimeDiff()
}

/**
 * 从数据库中获取关联到指定词元上的倒排列表
 * @param[in] env 存储着应用程序运行环境的结构体
 * @param[in] token_id 词元编号
 * @param[out] postings 获取到的倒排列表
 * @param[out] postings_len 获取到的倒排列表中的元素数
 */
func (p *PostingsManager) fetchPostings(tokenId int) (pl *encoding.PostingsList, length int, err error) {
	count, postings, err := p.database.GetPostings(tokenId)
	if err != nil || count == 0 || len(postings) == 0 {
		return nil, 0, err
	}

	pl, length, err = encoding.DecodePostings(postings, p.compress)
	if err != nil {
		fmt.Println("postings list decode error", err)
		return
	}
	if count != length {
		err = errors.New(fmt.Sprintf("postings list decode error: stored:%d decoded:%d.", count, length))
		return nil, 0, err
	}
	return
}

/**
 * 将内存上（小倒排索引中）的倒排列表与存储器上的倒排列表合并后存储到数据库中
 * @param[in] env 存储着应用程序运行环境的结构体
 * @param[in] p 含有倒排列表的倒排索引中的索引项
 */
func (p *PostingsManager) updatePostings(tokenId int, items *tokenIndexItems) {
	oldPostings, length, err := p.fetchPostings(tokenId)
	if err != nil {
		fmt.Printf("cannot fetch old postings list of token(%d) for update. %s\n", tokenId, err.Error())
		return
	}

	if length != 0{
		items.postings = encoding.Merge(items.postings, oldPostings)
		items.docCount += length
	}
	if items.docCount != items.postings.Length(){
		fmt.Println("length miss", tokenId, items.docCount, items.postings.Length(), length, items.docCount - length)
	}
	data := encoding.EncodePostings(items.postings, p.compress, items.docCount)
	p.database.UpdatePostings(tokenId, items.docCount, data)
}

func (p *PostingsManager) Merge(index map[int]*tokenIndexItems) {
	if len(p.index) == 0 {
		p.index = index
		return
	}
	for k, v := range index {
		if v2, ok := p.index[k]; ok {
			v2.merge(v)
		} else {
			p.index[k] = v
		}
	}
}
