package engine

import (
	"errors"
	"github.com/longbai/wiser-go/db"
)

var ErrEmptyDocument error = errors.New("empty document")

type TokenPersistent interface {
	PersistToken(token string)(id, count int, err error)
	GetTokenId(token string)(id, count int, err error)
	GetToken(id int)(token string, err error)
}

type DocumentPersistent interface {
	GetDocumentId(title string) (int, error)
	GetDocumentTitle(id int) (string, error)
	PersistDocument(title, body string) (int, error)
	GetDocumentCount() (count int, err error)
}

type PostingsPersistent interface {
	GetPostings(id int) (count int, postings []byte, err error)
	UpdatePostings(tokenId int, docCount int, postings []byte) (err error)
}

type SettingPersistent interface {
	SetSettings(name, value string) error
	GetSettings(name string) (value string, err error)
}

type dbDocumentPersistent struct {
	*db.SqliteDb
}

func (d *dbDocumentPersistent) GetDocumentId(title string) (int, error){
	return d.SqliteDb.GetDocumentId(title)
}

func (d *dbDocumentPersistent) GetDocumentTitle(id int) (string, error){
	return d.SqliteDb.GetDocumentTitle(id)
}

func (d *dbDocumentPersistent) PersistDocument(title, body string) (did int, err error){
	if title == "" || body == "" {
		return -1, ErrEmptyDocument
	}
	did, err = d.SqliteDb.AddDocument(title, body)
	return
}

func (d *dbDocumentPersistent) GetDocumentCount() (count int, err error){
	return d.SqliteDb.GetDocumentCount()
}

type dbTokenPersistent struct {
	*db.SqliteDb
}

func (d *dbTokenPersistent)PersistToken(token string)(id, count int, err error) {
	return d.SqliteDb.GetTokenId(token, true)
}

func (d *dbTokenPersistent)GetTokenId(token string)(id, count int, err error) {
	return d.SqliteDb.GetTokenId(token, false)
}

func (d *dbTokenPersistent)GetToken(id int)(token string, err error) {
	return d.SqliteDb.GetToken(id)
}
