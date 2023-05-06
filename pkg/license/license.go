package license

import (
	"encoding/base64"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io/ioutil"
	"strconv"
)

// License = AesKey32 + AesEnc(data).length + AesEnc(data) + Base64(RSAPubKey).length + Base64(RSAPubKey) + RsaSign(AesEnc(data));
func GeneratorLicence(pubKeyFileName, priKeyFileName string, data []byte) (string, error) {
	var licence string

	var aeskey = RandStringBytesMaskImpr(32)
	xpass, err := AesEncrypt(data, aeskey)
	if err != nil {
		return "", err
	}
	pass64 := base64.StdEncoding.EncodeToString(xpass)

	rsa, err := NewGoRSA("", priKeyFileName)
	if err != nil {
		return "", err
	}
	sign, err := rsa.Sign(pass64)
	if err != nil {
		return "", err
	}
	publicKey, err := ioutil.ReadFile(pubKeyFileName)
	if err != nil {
		return "", err
	}
	rsaPubKey := base64.StdEncoding.EncodeToString(publicKey)
	pass64Lenth := strconv.Itoa(len(pass64))
	licence = string(aeskey) + strconv.Itoa(len(pass64Lenth)) + strconv.Itoa(len(pass64)) + pass64 + strconv.Itoa(len(rsaPubKey)) + rsaPubKey + sign
	return licence, nil
}

func VerifyLicence(licenceFileName string) ([]byte, error) {
	licenceData, err := ioutil.ReadFile(licenceFileName)
	if err != nil {
		return nil, err
	}
	if len(licenceData) == 0 {
		return nil, errors.New("licence is empty")
	}
	licence := string(licenceData)
	aesKey := licence[0:32]
	passLenth, err := strconv.Atoi(licence[32:33])
	if err != nil {
		return nil, err
	}
	encDataLength, err := strconv.Atoi(licence[33 : 33+passLenth])
	if err != nil {
		return nil, err
	}
	encData := licence[33+passLenth : 33+passLenth+encDataLength]

	rsaIndex := 33 + passLenth + encDataLength
	rsaPubLength, err := strconv.Atoi(licence[rsaIndex : rsaIndex+4])
	if err != nil {
		return nil, err
	}
	rsaPubKey := licence[rsaIndex+4 : rsaIndex+4+rsaPubLength]

	sign := licence[rsaIndex+4+rsaPubLength:]
	pubKey, err := base64.StdEncoding.DecodeString(rsaPubKey)
	if err != nil {
		return nil, err
	}
	rsa, err := NewGoRSAPub(pubKey)
	if err != nil {
		return nil, err
	}
	err = rsa.Verify(encData, sign)
	if err != nil {
		return nil, err
	}

	bytesPass, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return nil, err
	}

	data, err := AesDecrypt(bytesPass, []byte(aesKey))
	if err != nil {
		return nil, err
	}
	hlog.Info(string(data))
	return data, nil
}
