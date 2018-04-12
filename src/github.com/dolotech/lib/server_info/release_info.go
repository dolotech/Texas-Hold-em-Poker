package serverInfo

import (
	"time"
	"github.com/dolotech/lib/pse"
	"runtime"
	"github.com/labstack/gommon/bytes"
	"encoding/json"
)

// 服务器版本信息
type ServerInfo struct {
	GitCommit     string    `json:"git_commit"`     // Git上面的commit版本号
	CompileTime   string    `json:"compile_time"`   // 服务器二进制文件编译时间
	ServerVersion string    `json:"server_version"` // 服务器当前产品开发版本
	StarupTime    time.Time `json:"starup_time"`    // 服务器启动时间
	GoVersion     string    `json:"go_version"`     // 编译服务器的golang版本
	Memery        string    `json:"memery"`
	Cores         int       `json:"cores"`
	CPU           float64   `json:"cpu"`
}

var info ServerInfo

func GetInfo() ServerInfo {
	updateUsage(&info)
	return info
}

func GetInfoBytes() []byte {
	o:=GetInfo()
	b,_:=json.Marshal(o)
	return b
}

func SetInfo(commit, compile, version string) {
	info.Cores = runtime.NumCPU()
	info.GoVersion = runtime.Version()
	info.GitCommit = commit
	info.CompileTime = compile
	info.ServerVersion = version
	info.StarupTime = time.Now()
}

func updateUsage(v *ServerInfo) {
	var rss, vss int64
	var pcpu float64
	pse.ProcUsage(&pcpu, &rss, &vss)
	v.Memery = bytes.Format(rss)
	v.CPU = pcpu
}
