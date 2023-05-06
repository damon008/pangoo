package license

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io/ioutil"
	"testing"
)

//根据授权json信息生成授权文件
func TestGeneratorLicence(t *testing.T) {
	/*data := `
{
    "appId": "***.com",
    "issuedTime": 1669082100,
    "notBefore": 1669082100,
    "notAfter": 1669084189,
    "customerInfo": "***公司",
    "authorization": "all,training,inference",
    "machineCodes": [],
    "nodeNum":0
}
`*/
	data := `{
    "appId":"***.com",
    "issuedTime":1669082100,
    "notBefore":1669082100,
    "notAfter":1764643298,
    "customerInfo":"***公司",
    "authorization":"all,training,inference",
    "machineCodes":[],
    "nodeNum":0
}`
	licence, err := GeneratorLicence("id_rsa.pub", "id_rsa", []byte(data))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(licence)

	err = ioutil.WriteFile("LICENSE", []byte(licence), 0644)
	if err != nil {
		hlog.Error("write into file of licence err %s", licence)
		return
	}
}

//拿到授权json信息
func TestGetContentData(t *testing.T) {
	byteLicence, err := VerifyLicence("licence")
	if err != nil {
		hlog.Error("verify the licence error", err)
		return
	}
	licenceContent := IssuedContent{}
	err = sonic.Unmarshal(byteLicence, &licenceContent)
	if err != nil {
		hlog.Error("unmarshal the licenceContent error", err)
		return
	}
	hlog.Info(licenceContent)
}


func TestVerifyLicence(t *testing.T) {
	data, err := VerifyLicence("licence")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func TestConvertObj(t *testing.T) {
	byteLicence, err := VerifyLicence("licence")
	licenceContent := IssuedContent{}
	err = sonic.Unmarshal(byteLicence, &licenceContent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", licenceContent)
}


/*func TestGeneratorLicence1(t *testing.T) {
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
	licence, err := GeneratorLicence("id_rsa1.pub", "id_rsa1", []byte(data))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(licence)
}*/
