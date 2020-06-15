package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db                    *sql.DB
	get_document_id_st    *sql.Stmt
	get_document_title_st *sql.Stmt
	insert_document_st    *sql.Stmt
	update_document_st    *sql.Stmt
	get_token_id_st       *sql.Stmt
	get_token_st          *sql.Stmt
	store_token_st        *sql.Stmt
	get_postings_st       *sql.Stmt
	update_postings_st    *sql.Stmt
	get_settings_st       *sql.Stmt
	set_settings_st   *sql.Stmt
	get_document_count_st *sql.Stmt

	begin_st    *sql.Stmt
	commit_st   *sql.Stmt
	rollback_st *sql.Stmt

	tx *sql.Tx
}

func Open(path string) (*Db, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	ret := Db{}
	ret.db = d

	d.Exec("CREATE TABLE settings  ( key TEXT PRIMARY KEY, value TEXT);")
	d.Exec("CREATE TABLE documents ( id INTEGER PRIMARY KEY, title TEXT NOT NULL, body TEXT NOT NULL);")
	d.Exec("CREATE TABLE tokens    ( id INTEGER PRIMARY KEY, token TEXT NOT NULL, docs_count INT NOT NULL, postings  BLOB NOT NULL);")

	d.Exec("CREATE UNIQUE INDEX token_index ON tokens(token);")
	d.Exec("CREATE UNIQUE INDEX title_index ON documents(title);")

	ret.get_document_id_st, _ = d.Prepare("SELECT id FROM documents WHERE title = ?;")
	ret.get_document_title_st, _ = d.Prepare("SELECT title FROM documents WHERE id = ?;")
	ret.insert_document_st, _ = d.Prepare("INSERT INTO documents (title, body) VALUES (?, ?);")
	ret.update_document_st, _ = d.Prepare("UPDATE documents set body = ? WHERE id = ?;")

	ret.get_token_id_st, _ = d.Prepare("SELECT id, docs_count FROM tokens WHERE token = ?;")
	ret.get_token_st, _ = d.Prepare("SELECT token FROM tokens WHERE id = ?;")
	ret.store_token_st, _ = d.Prepare("INSERT OR IGNORE INTO tokens (token, docs_count, postings) VALUES (?, 0, ?);")
	ret.get_postings_st, _ = d.Prepare("SELECT docs_count, postings FROM tokens WHERE id = ?;")
	ret.update_postings_st, _ = d.Prepare("UPDATE tokens SET docs_count = ?, postings = ? WHERE id = ?;")
	ret.get_settings_st, _ = d.Prepare("SELECT value FROM settings WHERE key = ?;")
	ret.set_settings_st, _ = d.Prepare("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?);")
	ret.get_document_count_st, _ = d.Prepare("SELECT COUNT(*) FROM documents;")

	ret.begin_st, _ = d.Prepare("BEGIN;")
	ret.commit_st, _ = d.Prepare("COMMIT;")
	ret.rollback_st, _ = d.Prepare("ROLLBACK;")

	return &ret, nil
}

func (db *Db) Close() {
	db.db.Close()
}

func (db *Db) GetDocumentId(title string) (int32, error) {
	return 0, nil
}

func (db *Db) GetDocumentTitle(id int32) (string, error) {
	return "", nil
}

func (db *Db) AddDocument(title, body string) error {
	return nil
}

func (db *Db) GetTokenId() {

}

func (db *Db) GetToken() {

}

func (db *Db) GetPostings() {

}

func (db *Db) UpdatePostings(tokenId int32) {

}

func (db *Db) SetSettings(name, value string) error{
	_, err := db.set_settings_st.Exec(name, value)
	return err
}

func (db *Db) GetSettings(name string) (string, error) {
	var value string
	err := db.get_settings_st.QueryRow(name).Scan(&value)
	return value, err
}

func (db *Db) GetDocumentCount() int {
	return 0
}

func (db *Db) Begin() {
	db.tx, _ = db.db.Begin()
}

func (db *Db) Commit() {
	db.tx.Commit()
}

func (db *Db) Rollback() {
	db.tx.Rollback()
}
