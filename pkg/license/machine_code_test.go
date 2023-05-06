package license

import (
	"fmt"
	"testing"
)

// must on linux os
// sudo dmidecode -t 1 | grep Serial
// sudo dmidecode -s system-serial-number
func TestGetCPUInfo(t *testing.T) {
	info, err := GetSystemSerialNum()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
