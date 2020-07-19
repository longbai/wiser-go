package engine

import (
	"github.com/longbai/wiser-go/db"
)

const QueryDocId = -1

type Engine struct {
	DocumentPersistent
	PostingsPersistent
	SettingPersistent
	buffer *PostingsManager
	t *TextProcessing
}

func NewEngine(d *db.SqliteDb, compress string, flushThreshold int) *Engine{
	return &Engine{
		DocumentPersistent: &dbDocumentPersistent{d},
		PostingsPersistent: d,
		SettingPersistent:  d,
		buffer:             NewPostingsManager(d, compress),
		t:                  &TextProcessing{&dbTokenPersistent{d}},
	}
}


func (e *Engine)BuildPostings(title, body string)(err error) {
	if title == "" || body == "" {
		return
	}
	var did int
	did, err = e.PersistDocument(title, body)
	if err != nil {
		return err
	}
	v, err := e.t.TextToPostingsLists(did, body)
	if err != nil {
		return err
	}
	e.buffer.Merge(v)
	return
}

func (e *Engine)Flush(flushThreshold int){
	e.buffer.Flush(flushThreshold)
}
