package license

import (
	"encoding/binary"
	"net"
	"time"
)

const ntpEpochOffset = 2208988800

type packet struct {
	Settings       uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimeSec     uint32
	RefTimeFrac    uint32
	OrigTimeSec    uint32
	OrigTimeFrac   uint32
	RxTimeSec      uint32
	RxTimeFrac     uint32
	TxTimeSec      uint32
	TxTimeFrac     uint32
}

func GetNTPTime() (time.Time, error) {
	systemTime := time.Now()
	host := "ntp.aliyun.com:123"

	conn, err := net.Dial("udp", host)
	if err != nil {
		return systemTime, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return systemTime, err
	}

	req := &packet{Settings: 0x1B}

	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		return systemTime, err
	}

	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		return systemTime, err
	}

	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32

	showtime := time.Unix(int64(secs), nanos)

	return showtime, nil
}

func GetTimeNow() time.Time {
	ntpTime, err := GetNTPTime()
	if err != nil {
		return time.Now()
	}
	return ntpTime
}
