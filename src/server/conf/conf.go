package conf

import (
	"log"
	"time"
)

var (
	// glog conf
	LogFlag = log.LstdFlags

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)



var Server struct {
	WSAddr       string
	CertFile     string
	KeyFile      string
	TCPAddr      string
	MaxConnNum   int
	ConsolePort  int
	ProfilePath  string
	DBMaxConnNum int
	DBUrl        string
}