package engine

import (
	"fmt"
	"unicode/utf8"
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
	case '“':
		fallthrough
	case '”':
		fallthrough
	case '？':
		return true
	default:
		return false
	}
}
// 2-gram
func biGramSplit(body string, f func(string, int) error){
	ignoreLast := true
	var lastRune rune
	charPos := 0
	for i, r := range body {
		if ignoreChar(r) {
			if !ignoreLast && i != 0{
				s := string([]rune{lastRune})
				err:= f(s, charPos-1)
				if err != nil {
					return
				}
			}
			ignoreLast = true
			continue
		}else {
			charPos++
		}

		if !ignoreLast && i != 0 {
			s := string([]rune{lastRune, r})
			err:= f(s, charPos-2)
			if err != nil {
				return
			}
		}
		ignoreLast = false
		lastRune = r
	}
}

func (p *TokenIndex) TextToPostingsLists(documentId int, body string) (err error) {
	p2 := &TokenIndex{
		index:    make(map[int]*tokenIndexItems),
		database: p.database,
	}
	fmt.Println("body size", utf8.RuneCountInString(body))
	count := 0
	biGramSplit(body, func(s string, pos int) error {
		count++
		return p2.tokenToPostingsList(documentId, s, pos)
	})
	fmt.Println("token count", count)
	p.Merge(p2)
	return
}

func (p *TokenIndex) tokenToPostingsList(documentId int, token string, position int) error {
	id, count, err := p.database.GetTokenId(token, documentId != QueryDocId)
	if err != nil {
		return err
	}
	if documentId == QueryDocId || count == 0{
		count = 1
	}
	entry, ok := p.index[id]
	if !ok {
		entry = &tokenIndexItems{
			tokenId:        id,
			token: token,
			docCount:       count,
			positionsCount: 0,
			postings: &encoding.PostingsList{
				DocumentId: documentId,
				Positions:  nil,
				Next:       nil,
			},
		}
		p.index[id] = entry

	}
	entry.postings.Positions = append(entry.postings.Positions, position)
	entry.positionsCount++
	return nil
}
