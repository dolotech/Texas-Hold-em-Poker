package algorithm
/*

type Card byte

func (this *Card) Equal(card byte) bool {
	return byte(*this) == card
}
func (this *Card) Byte() byte {
	return byte(*this)
}
func (this *Card) Value() byte {
	return byte(*this) & 0xF
}

func (this *Card) Color() byte {
	return byte(*this) >> SUITSIZE
}

const NilCard = Card(0)

*/

func  ToV(v byte) byte {
	return v& 0xF
}

func  ToC(v byte) byte {
	return v >> SUITSIZE
}