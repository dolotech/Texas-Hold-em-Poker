package algorithm

import (
	"strings"
	"fmt"
)

func (this *Cards) Bytes() []byte {
	b := make([]byte, len(*this))
	for k, v := range *this {
		b[k] = byte(v)
	}
	return b
}
func (this *Cards) Len() int {
	return len(*this)
}
func (this *Cards) Take() byte {
	card := (*this)[0]
	(*this) = (*this)[1:]
	return card
}
func (this *Cards) Append(cards ...byte) Cards {
	cs := make([]byte, 0, len(cards)+len(*this))
	cs = append(cs, (*this)...)
	cs = append(cs, cards...)
	return cs
}

func (this *Cards) Equal(cards []byte) bool {
	if len(*this) != len(cards) {
		return false
	}
	for k, v := range *this {
		if cards[k] != v {
			return false
		}
	}
	return true
}
func Color(color byte) (char string) {
	switch color {
	case 0:
		char = "♦"
	case 1:
		char = "♣"
	case 2:
		char = "♥"
	case 3:
		char = "♠"
	}
	return
}

func String2Num(c byte) (n byte) {
	switch c {
	case '2':
		n = 2
	case '3':
		n = 3
	case '4':
		n = 4
	case '5':
		n = 5
	case '6':
		n = 6
	case '7':
		n = 7
	case '8':
		n = 8
	case '9':
		n = 9
	case 'T':
		n = 0xA
	case 'J':
		n = 0xB
	case 'Q':
		n = 0xC
	case 'K':
		n = 0xD
	case 'A':
		n = 0xE
	}
	return
}
func Num2String(n byte) (c byte) {
	switch n {
	case 2:
		c = '2'
	case 3:
		c = '3'
	case 4:
		c = '4'
	case 5:
		c = '5'
	case 6:
		c = '6'
	case 7:
		c = '7'
	case 8:
		c = '8'
	case 9:
		c = '9'
	case 0xA:
		c = 'T'
	case 0xB:
		c = 'J'
	case 0xC:
		c = 'Q'
	case 0xD:
		c = 'K'
	case 0xE:
		c = 'A'
	}
	return
}

func (this *Cards) SetByString(str string) {
	array := strings.Split(str, " ")
	*this = make([]byte, len(array))
	for k, v := range array {
		(*this)[k] = String2Num(byte(v[0]))
	}

}
func (this *Cards) String() (str string) {
	for k, v := range *this {
		color := Color(v)
		value := Num2String(v)
		str += string(color) + string(value)
		if k < len(*this)-1 {
			str += " "
		}
	}
	return
}

func (this *Cards) Hex() string {
	return fmt.Sprintf("%#v", *this)
}
