package engine

import (
	"github.com/longbai/wiser-go/db"
)

const QueryDocId = -1

type Engine struct {
	TokenPersistent
	DocumentPersistent
	PostingsPersistent
	SettingPersistent
	buffer *TokenIndex
}

func NewEngine(d *db.Db, compress string, flushThreshold int) *Engine{
	return &Engine{
		TokenPersistent:    &dbTokenPersistent{d},
		DocumentPersistent: &dbDocumentPersistent{d},
		PostingsPersistent: d,
		SettingPersistent:  d,
		buffer: NewTokenIndex(d, compress),
	}
}


func (e *Engine)BuildPostings(title, body string)(err error) {
	if title == "" || body == "" {
		return
	}
	var did int
	did, err = e.PersistDocument(title, body)
	err = e.buffer.TextToPostingsLists(did, body)
	return
}

func (e *Engine)Flush(flushThreshold int){
	e.buffer.Flush(flushThreshold)
}
