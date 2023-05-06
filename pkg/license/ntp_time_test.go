package license

import (
	"fmt"
	"testing"
)

func TestGetNTPTime(t *testing.T) {

	ntime, err := GetNTPTime()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ntime)
	fmt.Println(ntime.Unix())
}
