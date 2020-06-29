package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/longbai/wiser-go/db"
	"github.com/longbai/wiser-go/engine"
	"github.com/longbai/wiser-go/source"
	"github.com/longbai/wiser-go/util"
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

var scan bool
func main() {
	compressMethod := flag.String("c", "golomb", "compress method for postings list(none   : don't compress;golomb : Golomb-Rice,default)")
	wikipediaDump := flag.String("x", "", "wikipedia dump xml path for indexing")
	queryStr := flag.String("q", "", "query for search")
	maxIndexCount := flag.Int("m", -1, "max count for indexing document")
	iibuThreshold := flag.Int("t", 2048, "inverted index buffer merge threshold")
	enablePhraseSearch := flag.Bool("s", true, "enable phrase search")
	trace := flag.Bool("scan", false, "trace wiki")
	flag.Parse()
	scan = *trace
	args := flag.Args()
	if len(args) == 0 {
		flag.PrintDefaults()
		return
	}
	dbPath := args[len(args)-1]
	if *wikipediaDump != "" {
		_, e := os.Stat(dbPath)
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
		err = construct(database, *compressMethod, *wikipediaDump, *maxIndexCount, *iibuThreshold)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *queryStr != "" {
		query(database, *queryStr, *enablePhraseSearch)
	}
}

func query(database *db.Db, query string, enablePhraseSearch bool) {
	cm, _ := database.GetSettings("compress_method")
	indexCount, _ := database.GetDocumentCount()
	engine.Search(query, cm, indexCount, database, enablePhraseSearch)
}

func construct(database *db.Db, compressMethod string, wikipediaDump string, maxIndexCount, iibuThreshold int) (err error) {
	err = database.SetSettings("compress_method", compressMethod)
	if err != nil {
		return err
	}
	database.Begin()

	engine1 := engine.NewEngine(database, compressMethod, iibuThreshold)
	if err = source.LoadWiki(wikipediaDump, maxIndexCount, func(title, body string) (err error) {
		if scan {
			return
		}
		err = engine1.BuildPostings(title, body)
		engine1.Flush()
		return
	}); err != nil {
		engine1.Flush()
		database.Commit()
	} else {
		database.Rollback()
	}
	return
}
