package utils

import (
	"math/rand"
	"time"
	"sync"
	//crand "crypto/rand"
)

var o *rand.Rand = rand.New(rand.NewSource(TimestampNano()))
var random_mux_ sync.Mutex

func RandInt64() (r int64) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Int63())
}

func RandInt32() (r int32) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Int31())
}

func RandUint32() (r uint32) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Uint32())
}

func RandInt64N(n int64) (r int64) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Int63n(n))
}

func RandInt32N(n int32) (r int32) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Int31n(n))
}

var randomChan chan uint32

func randUint32() {
	randomChan = make(chan uint32, 1024)
	go func() {
		var numstr uint32
		for {
			numstr = RandUint32()
			select {
			case randomChan <- numstr:
			}
			<-time.After(time.Millisecond * 100)
		}
	}()
}

func GetRandUint32() uint32 {
	return <-randomChan
}


func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(length int) string {

	//rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


var letters_ = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZ123456789")

func RandomString_(length int) string {

	//rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters_[rand.Intn(len(letters_))]
	}
	return string(b)
}