package license

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	//key的长度必须是16、24或者32字节，分别用于选择AES-128, AES-192, or AES-256
	var aeskey = RandStringBytesMaskImpr(32)
	fmt.Println("key:" + string(aeskey))
	data := `{
    "appId": "bitahub.com",
    "issuedTime": 1595951714,
    "notBefore": 1538671712,
    "notAfter": 1640966400,
    "customerInfo": "XXX公司",
    "authorization": "all,training,inference",
    "machineCode": "XXX",
	}`
	pass := []byte(data)
	xpass, err := AesEncrypt(pass, aeskey)
	if err != nil {
		fmt.Println(err)
		return
	}

	pass64 := base64.StdEncoding.EncodeToString(xpass)
	fmt.Printf("加密后:%v\n", pass64)

	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	if err != nil {
		fmt.Println(err)
		return
	}

	tpass, err := AesDecrypt(bytesPass, aeskey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("解密后:%s\n", tpass)
}