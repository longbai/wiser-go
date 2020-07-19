package engine

import (
	"fmt"
	"github.com/longbai/wiser-go/encoding"
	"unicode/utf8"
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


type TextProcessing struct {
	TokenPersistent
}

func (p *TextProcessing) TextToPostingsLists(documentId int, body string) (index map[int]*tokenIndexItems, err error) {
	fmt.Println("body size", utf8.RuneCountInString(body))
	count := 0
	index = make(map[int]*tokenIndexItems)
	biGramSplit(body, func(s string, pos int) error {
		count++
		return p.tokenToPostingsList(index, documentId, s, pos)
	})
	fmt.Println("split token count", count)
	return
}

func (p *TextProcessing) tokenToPostingsList(index map[int]*tokenIndexItems, documentId int, token string, position int) error {
	var pt = p.PersistToken
	if documentId == QueryDocId{
		pt = p.GetTokenId
	}
	id, _, err := pt(token)
	if err != nil {
		return err
	}

	entry, ok := index[id]
	if !ok {
		entry = &tokenIndexItems{
			tokenId:        id,
			token: token,
			docCount:       1,
			positionsCount: 0,
			postings: &encoding.PostingsList{
				DocumentId: documentId,
				Positions:  nil,
				Next:       nil,
			},
		}
		index[id] = entry

	}
	entry.postings.Positions = append(entry.postings.Positions, position)
	entry.positionsCount++
	return nil
}

type tokenIndexItems struct {
	token		string // for debug
	tokenId        int
	docCount       int
	positionsCount int
	postings       *encoding.PostingsList
}

func (p *tokenIndexItems) merge(other *tokenIndexItems) {
	p.docCount += other.docCount
	p.postings = encoding.Merge(p.postings, other.postings)
}

