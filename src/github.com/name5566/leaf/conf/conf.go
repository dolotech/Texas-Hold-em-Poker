package conf

var (
	LenStackBuf = 4096

	// glog
/*	LogLevel string
	LogPath  string
	LogFlag  int*/

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)
