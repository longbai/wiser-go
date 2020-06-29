package engine

import (
	"fmt"
	"github.com/longbai/wiser-go/db"
)

type Engine struct {
	TokenPersistent
	DocumentPersistent
	PostingsPersistent
	SettingPersistent
	buffer *TokenIndex
	flushThreshold int
}

func NewEngine(d *db.Db, compress string, flushThreshold int) *Engine{
	return &Engine{
		TokenPersistent:    &dbTokenPersistent{d},
		DocumentPersistent: &dbDocumentPersistent{d},
		PostingsPersistent: d,
		SettingPersistent:  d,
		buffer: NewTokenIndex(d, compress),
		flushThreshold: flushThreshold,
	}
}


func (e *Engine)BuildPostings(title, body string)(err error) {
	if title == "" || body == "" {
		return
	}
	var did int
	did, err = e.PersistDocument(title, body)
	err = e.buffer.TextToPostingsLists(did, body)
	fmt.Println("text", err)
	return
}

func (e *Engine)Flush(){
	e.buffer.Flush(e.flushThreshold)
}
