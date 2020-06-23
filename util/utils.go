package util

import (
	"fmt"
	"log"
	"time"
)

var preTime *time.Time

func PrintTimeDiff() {
	t := time.Now()
	if preTime == nil{
		preTime = &t
		fmt.Println("[time]", t.String())
	} else {
		diff := t.UnixNano() - (*preTime).UnixNano()
		fmt.Printf("[time] %s (diff %d)", t.String(), diff)
	}
}

func Abort(){
	log.Panicln("abort")
}
