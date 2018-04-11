package server

import (
	"testing"
	"time"
	"log"
)
var writeChan= make(chan int,10)
func Test_ws(t *testing.T)  {
	go test()
	writeChan <- 1
	writeChan <-0
	<- time.After(time.Hour)
}

func test()  {

	log.Print("hello1")
	for b := range writeChan {
		if b == 0 {
			break
		}

		log.Print(b)
	}
	log.Print("hello")
}
