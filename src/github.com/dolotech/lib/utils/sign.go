package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

//生成签名
func LoginSign(gameid string, device string) (string, error) {
	parra := fmt.Sprintf("gameid=%s&deviceId=%s", gameid, device)
	parameters := strings.Split(parra, "&")
	sort_parameters := sort.StringSlice(parameters)
	sort.Sort(sort_parameters)
	parameter := ""
	for i := 0; i < len(sort_parameters); i++ {
		parameter += "&" + sort_parameters[i]
	}
	parameter = string([]byte(parameter)[1:])

	h := md5.New()

	_, err := h.Write([]byte(parameter + "&key=" + "19a87399fde1bccbf04a5eaa018ea0df"))
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}
