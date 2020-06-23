package engine

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/longbai/wiser-go/db"
	"github.com/longbai/wiser-go/encoding"
)

const NGram = 2

type phraseSearchCursor struct {
	positions []int /* 位置信息 */
	base      int   /* 词元在查询中的位置 */
	current   int   /* 当前的位置信息 */
}

type docSearchCursor struct {
	documents *encoding.PostingsList /* 文档编号的序列 */
	current   *encoding.PostingsList /* 当前的文档编号 */
}

func Search(query, compressMethod string, indexCount int, d *db.Db, enablePhraseSearch bool) error {
	if len(query) < NGram {
		fmt.Println("too short")
		return errors.New("query short than token length")
	}
	t, err := splitQueryToTokens(query, compressMethod, d)
	if err != nil {
		return err
	}
	r := t.searchDocs(enablePhraseSearch, indexCount)
	r.sort()
	r.print()
	return nil
}

func searchPhrase(tokens []*tokenIndexItems, docCursors []docSearchCursor) (phraseCount int){

	positions := 0
	/* 获取查询中词元的总数 */
	for _, v := range tokens {
		positions += v.positionsCount
	}

	/* 初始化游标 */
	cursors := make([]phraseSearchCursor, positions)
	cursorPos := 0
	for k, v := range tokens {
		for _, p := range v.postings.Positions {
			cursors[cursorPos].base = p
			cursors[cursorPos].positions = docCursors[k].current.Positions
			cursors[cursorPos].current = 0
		}
	}

	/* 检索短语 */
	for cursors[0].current < len(cursors[0].positions) {
		var relPosition, nextRelPosition int
		nextRelPosition = cursors[cursorPos].positions[cursors[0].current] - cursors[0].base
		relPosition = nextRelPosition
		/* 对于除词元A以外的词元，不断地向后读取其出现位置，直到其偏移量不小于词元A的偏移量为止 */
		for _, cur := range cursors[1:] {
			for ; cur.current < len(cur.positions) && (cur.positions[cur.current]-cur.base) < relPosition; cur.current++ {
			}
			if cur.current == len(cur.positions) {
				return
			}
			/* 对于除词元A以外的词元，若其偏移量不等于A的偏移量，就退出循环 */
			off := cur.positions[cur.current] - cur.base
			if off != relPosition {
				nextRelPosition = off
				break
			}
		}
		if nextRelPosition < relPosition {
			/* 不断向后读取，直到词元A的偏移量不小于next_rel_position为止 */
			for cursors[0].current < len(cursors[0].positions)&&
				cursors[0].positions[cursors[0].current]- cursors[0].base < nextRelPosition   {
				cursors[0].current++
			}
		}
	}
	return
}

func calcTfIdf(tokens []*tokenIndexItems, cursors []docSearchCursor, indexCount int) float64 {
	var score float64
	for k, v := range tokens {
		idf := math.Log2(float64(indexCount)/float64(v.docCount))
		score += float64(cursors[k].current.PositionsCount()) * idf
	}

	return score
}

func (t *TokenIndex) searchDocs(phrase bool, indexCount int) SearchResults {
	var ret []documentScore
	items := t.sortItems()
	tokenCount := len(items)
	cursors := make([]docSearchCursor, tokenCount)
	for k, v := range items {
		if v.tokenId == 0 {
			return nil
		}
		pl, _, err := t.fetchPostings(v.tokenId)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if pl == nil {
			fmt.Println("no postings list")
			return nil
		}
		cursors[k].documents = pl
		cursors[k].current = pl
	}

	for cursors[0].current != nil{
		var docId, nextDocId int
		/* 将拥有文档最少的词元称作A */
		docId = cursors[0].current.DocumentId
		/* 对于除词元A以外的词元，不断获取其下一个document_id，直到当前的document_id不小于词元A的document_id为止 */
		for _, cur := range cursors {
			for cur.current != nil && cur.current.DocumentId < docId {
				cur.current = cur.current.Next
			}
			if cur.current == nil {
				return nil
			}
			/* 对于除词元A以外的词元，如果其document_id不等于词元A的document_id，*/
			/* 那么就将这个document_id设定为next_doc_id */
			if cur.current.DocumentId != docId {
				nextDocId = cur.current.DocumentId
				break
			}
		}

		if nextDocId > 0{
			/* 不断获取A的下一个document_id，直到其当前的document_id不小于next_doc_id为止 */
			for cursors[0].current != nil && cursors[0].current.DocumentId < nextDocId {
				cursors[0].current = cursors[0].current.Next
			}
		} else {
			phraseCount := -1
			if phrase {
				phraseCount = searchPhrase(items, cursors)
			}
			if phraseCount != 0{
				 doubleScore := calcTfIdf(items, cursors, indexCount)
				 ret = append(ret, documentScore{
					 docId: docId,
					 score: doubleScore,
				 })
			}
			cursors[0].current = cursors[0].current.Next
		}
	}
	return ret
}

func (t *TokenIndex) positionsCount() int {
	count := 0
	for _, v := range t.index {
		count += v.positionsCount
	}
	return count
}

type ByCount []*tokenIndexItems

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].docCount < a[j].docCount }

func (t *TokenIndex) sortItems() ByCount {
	items := make([]*tokenIndexItems, len(t.index))
	i:= 0
	for _, v := range t.index {
		items[i] = v
		i++
	}
	b := ByCount(items)
	sort.Sort(b)
	return b
}

func splitQueryToTokens(query, compressMethod string, d *db.Db) (*TokenIndex, error) {
	index := NewTokenIndex(d, compressMethod)
	err := index.TextToPostingsLists(0, query)
	if err != nil {
		return nil, err
	}
	return index, nil
}

type documentScore struct {
	docId int
	score float64
}

type SearchResults []documentScore

func (a SearchResults) Len() int           { return len(a) }
func (a SearchResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SearchResults) Less(i, j int) bool { return a[i].score > a[j].score }


func (r SearchResults) print() {

}

func (r SearchResults) sort() {
	sort.Sort(r)
}
