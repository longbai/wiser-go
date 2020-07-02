package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
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

func Open(path string) (*SqliteDb, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	ret := SqliteDb{}
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
	ret.get_document_count_st, _ = d.Prepare("SELECT COUNT(*) FROM documents;")

	ret.get_token_id_st, _ = d.Prepare("SELECT id, docs_count FROM tokens WHERE token = ?;")
	ret.get_token_st, _ = d.Prepare("SELECT token FROM tokens WHERE id = ?;")
	ret.store_token_st, _ = d.Prepare("INSERT OR IGNORE INTO tokens (token, docs_count, postings) VALUES (?, 0, ?);")

	ret.get_postings_st, _ = d.Prepare("SELECT docs_count, postings FROM tokens WHERE id = ?;")
	ret.update_postings_st, _ = d.Prepare("UPDATE tokens SET docs_count = ?, postings = ? WHERE id = ?;")

	ret.get_settings_st, _ = d.Prepare("SELECT value FROM settings WHERE key = ?;")
	ret.set_settings_st, _ = d.Prepare("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?);")


	ret.begin_st, _ = d.Prepare("BEGIN;")
	ret.commit_st, _ = d.Prepare("COMMIT;")
	ret.rollback_st, _ = d.Prepare("ROLLBACK;")

	return &ret, nil
}

func (db *SqliteDb) Close() {
	db.db.Close()
}

func (db *SqliteDb) GetDocumentId(title string) (int, error) {
	var value int
	err := db.get_document_id_st.QueryRow(title).Scan(&value)
	return value, err
}

func (db *SqliteDb) GetDocumentTitle(id int) (string, error) {
	var value string
	err := db.get_document_title_st.QueryRow(id).Scan(&value)
	return value, err
}

func (db *SqliteDb) AddDocument(title, body string) (did int, err error) {
	did, err = db.GetDocumentId(title)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if did != 0 {
		_, err = db.update_document_st.Exec(body, did)
	} else {
		var r sql.Result
		r, err = db.insert_document_st.Exec(title, body)
		var did1 int64
		did1, err = r.LastInsertId()
		did = int(did1)
	}
	return
}

func (db *SqliteDb) GetDocumentCount() (count int, err error) {
	err = db.get_document_count_st.QueryRow().Scan(&count)
	return
}

func (db *SqliteDb) GetTokenId(token string, insert bool)(id, count int, err error) {
	if insert {
		_, err = db.store_token_st.Exec(token, "")
		if err != nil {
			return
		}
	}
	err = db.get_token_id_st.QueryRow(token).Scan(&id, &count)
	return
}

func (db *SqliteDb) GetToken(id int)(token string, err error) {
	err = db.get_token_st.QueryRow(id).Scan(&token)
	return
}

func (db *SqliteDb) GetPostings(id int) (count int, postings []byte, err error){
	err = db.get_postings_st.QueryRow(id).Scan(&count, &postings)
	return
}

func (db *SqliteDb) UpdatePostings(tokenId int, docCount int, postings []byte) (err error){
	_, err = db.update_postings_st.Exec(docCount, postings, tokenId)
	return
}

func (db *SqliteDb) SetSettings(name, value string) error{
	_, err := db.set_settings_st.Exec(name, value)
	return err
}

func (db *SqliteDb) GetSettings(name string) (value string, err error) {
	err = db.get_settings_st.QueryRow(name).Scan(&value)
	return
}

func (db *SqliteDb) Begin() {
	db.tx, _ = db.db.Begin()
}

func (db *SqliteDb) Commit() {
	db.tx.Commit()
}

func (db *SqliteDb) Rollback() {
	db.tx.Rollback()
}
