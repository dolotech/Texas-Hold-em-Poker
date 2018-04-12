/**
 * Created by Michael on 2015/8/4.
 */
package utils

import (
	"encoding/json"
	"fmt"
	"testing"


)




func Test_StringAdd(t *testing.T) {
	t.Log(StringAdd("123"))
	t.Log(StringAdd("789786654567"))
	t.Log(StringAdd("sadfasdf"))
	t.Log(StringAdd("123944"))
	t.Log(StringAdd("1111111111"))
	t.Log(StringAdd("999999999"))
	t.Log(StringAdd("0"))
	t.Log(StringAdd("099"))

}
/*
func Test_ran(t *testing.T) {
	//for i := 0; i < 100; i++ {
	//	a := RandInt64()
	//	t.Log(Conver10to62(a), len(Conver10to62(a)))

	//}
	//str := "12345678901"
	str := Conver10to62(RandInt64())
	t.Log(str, len(str))
	if len(str) < 12 {
		var buf = make([]byte, 12-len(str))
		for i := 0; i < 12-len(str); i++ {
			buf[i] = 48
		}
		str = string(buf) + str
	}
	t.Log(len(str), str)
}*/
func Test_copy(t *testing.T) {
	a := AA{A: 999}
	b := AA{}
	//err :=Clone(b, a)
	t.Log(a, b)
}
func Test_AES(t *testing.T) {
	aesEnc := AesEncrypt{}
	aesEnc.SetKey([]byte("aalk;lkasjd;lkfj;alk"))
	doc := []byte("abcde号。")
	arrEncrypt, err := aesEnc.Encrypt(doc)
	glog.Infoln(string(arrEncrypt))
	if err != nil {
		glog.Infoln(string(arrEncrypt))
		return
	}
	strMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		glog.Infoln(string(arrEncrypt))
		return
	}
	glog.Infoln(string(strMsg))
}
func Test_XXTEA(t *testing.T) {
	/*str := "Hello World! 你好，中国！"
	key := "1234567890"
	encrypt_data := Encrypt([]byte(str), []byte(key))
	//glog.Infoln(base64.StdEncoding.EncodeToString(encrypt_data))
	decrypt_data := Decrypt(encrypt_data, []byte(key))*/

}
func TestPWD(t *testing.T) {
	t.Log(AalidataPwd("dolo0425"))
}

func TestPhone(t *testing.T) {
	t.Log(PhoneRegexp("8601593533372"))
}

type AA struct {
	CC
	A int `json:"a"`
}

type BB interface {
	Decode(b *[]byte) error
	Encode() (*[]byte, error)
}

type CC struct{}

func (this *CC) Decode(b *[]byte) error {
	return json.Unmarshal(*b, this)
}

func (this *CC) Encode() (*[]byte, error) {
	data, err := json.Marshal(this)
	return &data, err
}
