package license

import (
	"fmt"
	"testing"
)

func TestGoRSA_Sign(t *testing.T) {
	NewRSAFile("id_rsa1.pub", "id_rsa1", 4096)
	rsa, _ := NewGoRSA("id_rsa1.pub", "id_rsa1")
	//data := "XXX"
	/*data := `{
    "appId": "bitahub.com",
    "issuedTime": 1595951714,
    "notBefore": 1538671712,
    "notAfter": 1540966400,
    "customerInfo": "XXXXXX公司",
    "authorization": "all,training,inference",
    "machineCode": ["GPU-bed3da9b-952f-e16e-46e0-b0b9a6ec5b7d","GPU-bed3da9b-952f-e16e-46e0-b0b9a6ec5b7e"]
}`*/
	data := `{
    "appId": "bitahub.com",
    "issuedTime": 1595951714,
    "notBefore": 1538671712,
    "notAfter": 1540966400,
    "customerInfo": "XXXXXX公司",
    "authorization": "all,training,inference",
    "machineCodes": [],
    "nodeNum":0
}`
	sign, _ := rsa.Sign(data)
	fmt.Println(sign)
	err := rsa.Verify(data, sign)
	fmt.Println(err)

}
