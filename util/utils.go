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
		fmt.Println("[time]", t.String())
	} else {
		diff := t.UnixNano() - (*preTime).UnixNano()
		d := float64(diff)/1e9
		fmt.Printf("[time] %s (diff %f)\n", t.String(), d)
	}
	preTime = &t
}

func Abort(){
	log.Panicln("abort")
}
