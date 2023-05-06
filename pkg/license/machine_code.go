package license

import (
	"errors"
	"github.com/yumaojun03/dmidecode"
)

// sudo dmidecode -t 1 | grep Serial
func GetSystemSerialNum() (string, error) {
	decoder, err := dmidecode.New()
	if err != nil {
		return "", err
	}
	systemInfos, err := decoder.System()

	if len(systemInfos) <= 0 {
		return "", errors.New("dmidecode get system info failed")
	}

	return systemInfos[0].SerialNumber, nil
}
