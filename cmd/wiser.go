package main

import (
	"flag"
	"fmt"
	"github.com/longbai/wiser-go/db"
	"github.com/longbai/wiser-go/search"
	"github.com/longbai/wiser-go/source"
	"github.com/longbai/wiser-go/util"
	"os"
)

//"usage: %s [options] db_file\n"
//"options:\n"
//"  -c compress_method            : compress method for postings list\n"
//"  -x wikipedia_dump_xml         : wikipedia dump xml path for indexing\n"
//"  -q search_query               : query for search\n"
//"  -m max_index_count            : max count for indexing document\n"
//"  -t ii_buffer_update_threshold : inverted index buffer merge threshold\n"
//"  -s                            : don't use tokens' positions for search\n"
//"compress_methods:\n"
//"  none   : don't compress.\n"
//"  golomb : Golomb-Rice coding(default).\n",

func main() {
	compressMethod := flag.String("c", "golomb", "compress method for postings list(none   : don't compress;golomb : Golomb-Rice,default)")
	wikipediaDump := flag.String("x", "", "wikipedia dump xml path for indexing")
	queryStr := flag.String("q", "", "query for search")
	maxIndexCount := flag.Int("m", -1, "max count for indexing document")
	iibuThreshold := flag.Int("t", 2048, "inverted index buffer merge threshold")
	enablePhraseSearch := flag.Bool("s", true, "enable phrase search")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.PrintDefaults()
		return
	}
	dbPath := args[len(args)-1]
	if *wikipediaDump != "" {
		_, e :=os.Stat(dbPath)
		if e == nil {
			fmt.Println(dbPath, "is already exists!")
			return
		}
	}

	database, err := db.Open(dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer database.Close()
	defer util.PrintTimeDiff()
	util.PrintTimeDiff()

	if *wikipediaDump != "" {
		err = buildIndex(database, compressMethod, wikipediaDump, maxIndexCount, iibuThreshold)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *queryStr != "" {
		query(database, queryStr, enablePhraseSearch)
	}
}

func query(database *db.Db, query *string, enablePhraseSearch *bool) {
	cm, _ := database.GetSettings("compress_method")
	indexCount := database.GetDocumentCount()
	search.Search(*query, cm, indexCount, database, *enablePhraseSearch)
}

func buildIndex(database *db.Db, compressMethod *string, wikipediaDump *string, maxIndexCount, iibuThreshold *int) (err error){
	err = database.SetSettings("compress_method", *compressMethod)
	if err != nil {
		return err
	}
	database.Begin()
	if err = source.LoadWiki(*wikipediaDump, *maxIndexCount, func(title, body string) error {
		fmt.Println(title, body)
		return database.AddDocument(title, body)
	});err != nil {
		//add doc finish
		database.Commit()
	} else {
		database.Rollback()
	}
	return
}
