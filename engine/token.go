package engine

import (
	"github.com/longbai/wiser-go/encoding"
)

func ignoreChar(r rune) bool {
	switch r {
	case ' ':
		fallthrough
	case '\f':
		fallthrough
	case '\n':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
		fallthrough
	case '\v':
		fallthrough
	case '!':
		fallthrough
	case '"':
		fallthrough
	case '#':
		fallthrough
	case '$':
		fallthrough
	case '%':
		fallthrough
	case '&':
		fallthrough
	case '\'':
		fallthrough
	case '(':
		fallthrough
	case ')':
		fallthrough
	case '*':
		fallthrough
	case '+':
		fallthrough
	case ',':
		fallthrough
	case '-':
		fallthrough
	case '.':
		fallthrough
	case '/':
		fallthrough
	case ':':
		fallthrough
	case ';':
		fallthrough
	case '<':
		fallthrough
	case '=':
		fallthrough
	case '>':
		fallthrough
	case '?':
		fallthrough
	case '@':
		fallthrough
	case '[':
		fallthrough
	case '\\':
		fallthrough
	case ']':
		fallthrough
	case '^':
		fallthrough
	case '_':
		fallthrough
	case '`':
		fallthrough
	case '{':
		fallthrough
	case '|':
		fallthrough
	case '}':
		fallthrough
	case '~':
		fallthrough
	case '　': /* 全角空格 */
		fallthrough
	case '、':
		fallthrough
	case '。':
		fallthrough
	case '（':
		fallthrough
	case '）':
		fallthrough
	case '！':
		fallthrough
	case '，':
		fallthrough
	case '：':
		fallthrough
	case '；':
		fallthrough
	case '？':
		return true
	default:
		return false
	}
}

// expect ngram is 2
func (p *TokenIndex)TextToPostingsLists(documentId int, body string)(err error){
	var lastRune rune
	ignoreLast := false
	p2 := &TokenIndex{
		index:    make(map[int]*TokenIndexItems),
		database: p.database,
	}
	for i, r := range body {
		if ignoreChar(r) {
			ignoreLast = true
			continue
		}
		if !ignoreLast && i != 0 {
			s := string([]rune{lastRune, r})
			err = p2.TokenToPostingsList(documentId, s, i-1)
			if err != nil {
				return
			}
		}
		lastRune = r
	}
	p.Merge(p2)
	return
}

func (p *TokenIndex)TokenToPostingsList(documentId int, token string, position int) error {
	id, count, err := p.database.GetTokenId(token, documentId>0)
	if err != nil {
		return err
	}
	if documentId == 0 {
		count = 1
	}
	entry, ok := p.index[id]
	if !ok {
		entry = &TokenIndexItems{
			DocCount:       count,
			PositionsCount: 0,
			Postings:       &encoding.PostingsList{
				DocumentId:     documentId,
				Positions:      nil,
				Next:           nil,
			},
		}
		p.index[id] = entry

	}
	entry.Postings.Positions = append(entry.Postings.Positions, position)
	entry.PositionsCount++
	return nil
}
