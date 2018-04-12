package conf

import (
	"io/ioutil"
	//"leaf/glog"
	"encoding/json"
	"github.com/golang/glog"
)

var Server struct {
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
	DBMaxConnNum int
	DBUrl        string
}

func init() {
	data, err := ioutil.ReadFile("server.json")
	if err != nil {
		glog.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		glog.Fatal("%v", err)
	}
}
