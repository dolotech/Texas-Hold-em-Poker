package algorithm

import (
	"time"
	"math/rand"
)



//var n int64
//var a int64= 1<<62
// 洗牌
func (this *Cards) Shuffle() {
	*this = make([]Card, TOTAL)
	copy(*this, CARDS)
	source := rand.NewSource(time.Now().UnixNano() )
	//n ++
	//n %=a
	r := rand.New(source)
	for i := TOTAL - 1; i > 0; i-- {
		index := r.Int() % i
		(*this)[i], (*this)[index] = (*this)[index], (*this)[i]
	}

}
