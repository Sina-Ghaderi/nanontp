package getter

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)
const (
	LI_NO_WARNING      = 0
	LI_ALARM_CONDITION = 3
	VN_FIRST           = 1
	VN_LAST            = 4
	MODE_CLIENT        = 3
	FROM_1900_TO_1970  = 2208988800
)

func Serve(req []byte, addr net.Addr) ([]byte, error) {
	if validFormat(req) {
		res := generate(req, addr)
		return res, nil
	}
	return []byte{}, errors.New("invalid format.")
}

func validFormat(req []byte) bool {
	var l = req[0] >> 6
	var v = (req[0] << 2) >> 5
	var m = (req[0] << 5) >> 5
	if (l == LI_NO_WARNING) || (l == LI_ALARM_CONDITION) {
		if (v >= VN_FIRST) && (v <= VN_LAST) {
			if m == MODE_CLIENT {
				return true
			}
		}
	}
	return false
}
func unix2ntp(u int64) int64 {
	return u + FROM_1900_TO_1970
}

func ntp2unix(n int64) int64 {
	return n - FROM_1900_TO_1970
}

func int2bytes(i int64) []byte {
	var b = make([]byte, 4)
	h1 := i >> 24
	h2 := (i >> 16) - (h1 << 8)
	h3 := (i >> 8) - (h1 << 16) - (h2 << 8)
	h4 := byte(i)
	b[0] = byte(h1)
	b[1] = byte(h2)
	b[2] = byte(h3)
	b[3] = byte(h4)
	return b
}

func NTPPickup(upntp []string, ClientIP net.Addr) time.Time {
	for _, addr := range upntp {
		log.Println("\033[1;34mrequest\033[0m ---> asking for time from", ClientIP)
		log.Println("\033[1;36maccess\033[0m ----> trying to ntp server:", addr)
		custom, err := ntpaccessclient(addr)
		if err == nil {
			log.Println("\033[1;32msuccess\033[0m ---> time received from:", addr)
			log.Println("\033[1;32mresponse\033[0m --> answering to the client", ClientIP)
			return custom
		}
		log.Printf("\033[1;31merror\033[0m -----> ntp address: %v (%v)\n", addr, err)

	}
	log.Println("\033[1;31merror\033[0m -----> can not get anything from ntp servers!")
	log.Println("\033[1;33mwarining\033[0m --> passing system time to client (might be wrong!)")
	return time.Now()
}

var ArgNTP []string

func generate(req []byte, addr net.Addr) []byte {
	custom := NTPPickup(ArgNTP, addr)
	var second = unix2ntp(custom.Unix())
	var fraction = unix2ntp(int64(custom.Nanosecond()))
	var res = make([]byte, 48)
	var vn = req[0] & 0x38
	res[0] = vn + 4
	res[1] = 1
	res[2] = req[2]
	res[3] = 0xEC
	res[12] = 0x4E
	res[13] = 0x49
	res[14] = 0x43
	res[15] = 0x54
	copy(res[16:20], int2bytes(second)[0:])
	copy(res[24:32], req[40:48])
	copy(res[32:36], int2bytes(second)[0:])
	copy(res[36:40], int2bytes(fraction)[0:])
	copy(res[40:48], res[32:40])
	return res
}

type mode byte

const (
	reserved mode = 0 + iota
	symmetricActive
	symmetricPassive
	client
	server
	broadcast
	controlMessage
	reservedPrivate
)

type ntpTime struct {
	Seconds  uint32
	Fraction uint32
}

func (t ntpTime) UTC() time.Time {
	nsec := uint64(t.Seconds)*1e9 + (uint64(t.Fraction) * 1e9 >> 32)
	return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nsec))
}

type msg struct {
	LiVnMode       byte
	Stratum        byte
	Poll           byte
	Precision      byte
	RootDelay      uint32
	RootDispersion uint32
	ReferenceId    uint32
	ReferenceTime  ntpTime
	OriginTime     ntpTime
	ReceiveTime    ntpTime
	TransmitTime   ntpTime
}

func (m *msg) SetVersion(v byte) {
	m.LiVnMode = (m.LiVnMode & 0xc7) | v<<3
}

func (m *msg) SetMode(md mode) {
	m.LiVnMode = (m.LiVnMode & 0xf8) | byte(md)
}

func ntpaccessclient(host string) (time.Time, error) {
	raddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return time.Now(), err
	}

	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return time.Now(), err
	}
	defer con.Close()
	con.SetDeadline(time.Now().Add(5 * time.Second))

	m := new(msg)
	m.SetMode(client)
	m.SetVersion(4)

	err = binary.Write(con, binary.BigEndian, m)
	if err != nil {
		return time.Now(), err
	}

	err = binary.Read(con, binary.BigEndian, m)
	if err != nil {
		return time.Now(), err
	}

	t := m.ReceiveTime.UTC().Local()
	return t, nil
}
