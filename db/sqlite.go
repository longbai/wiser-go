package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db *sql.DB
}

func Open(path string)(*Db, error){
	d,err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	ret := Db{}
	ret.db = d
	return &ret, nil
}

func (db *Db)Close(){
	db.db.Close()
}

func (db *Db)GetDocumentId(title string)(int32, error){
	return 0, nil
}

func (db *Db)GetDocumentTitle(id int32)(string, error){
	return "", nil
}

func (db *Db)AddDocument(title, body string) error {
	return nil
}

func (db *Db)GetTokenId() {

}

func (db *Db)GetToken() {

}

func (db *Db)GetPostings() {

}

func (db *Db)UpdatePostings(tokenId int32) {

}

func (db *Db)SetSettings(name, value string) {

}

func (db *Db)GetSettings(name string) string {
	return ""
}

func (db *Db)GetDocumentCount() int {
	return 0
}

func (db *Db)Begin() {

}

func (db *Db)Commit() {

}

func (db *Db)Rollback() {

}
