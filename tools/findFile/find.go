package findFile

import (
	"regexp"
	"strings"
	"time"
)

//import (
//	"os"
//	"path/filepath"
//)

//func find() {

//}

func NewMatch(date int64, reg string) (m match) {
	if reg != "" {
		if strings.Index(reg, "*") == 0 {
			reg = "." + reg
		} else {
			reg = "^" + reg
		}
		reg += "$"
		reg, err := regexp.Compile(reg)
		if err == nil {
			m.reg = reg
		}
	}
	if date != 0 {
		if date < 0 {
			m.unixtime = time.Now().Unix() + date*24*60*60
		} else {
			m.unixtime = time.Now().Unix() - date*24*60*60
			m.less = true
		}
	}
	return
}
