package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//type AesEncrypt struct {
//	Key []byte
//}

func  SetKey(key []byte) ([]byte,error){
	keyLen := len(key)
	if keyLen < 16 {
		//panic("res key 长度不能小于16")
		return nil,errors.New("res key 长度不能小于16")
	}
	if keyLen >= 32 {
		//取前32个字节
		return  key[:32],nil
	}
	if keyLen >= 24 {
		//取前24个字节
		return key[:24],nil
	}
	//取前16个字节
	return key[:16],nil
}

//加密字符串
func  AesEncrypt(key []byte,strMesg []byte) ([]byte, error) {
	key,err:= SetKey(key)
	encrypted := make([]byte, len(strMesg))
	if err!= nil{
		return encrypted,err
	}
	var iv = key[:aes.BlockSize]

	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, strMesg)
	return encrypted, nil
}

//解密字符串
func AesDecrypt(key []byte,src []byte) (strDesc []byte, err error) {
	key,err= SetKey(key)
	if err!= nil{
		return
	}

	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return decrypted,err
}
