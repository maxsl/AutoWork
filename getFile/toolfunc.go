package getFile

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type cfg struct {
	Path     string `json:path`
	Host     string `json:host`
	ListenIP string `json:listenip`
	Port     string `json:port`
	Debug    bool   `json:debug`
	TempPath string `json:temppath`
	Log      string `json:log`
}

var (
	configFile string
)

var config cfg = cfg{Path: "/tmp",
	Host:     "http://127.0.0.1/getFile/finished",
	ListenIP: "0.0.0.0",
	Port:     "1789",
	Debug:    true,
	TempPath: "tmp/",
	Log:      ""}

var Log *log.Logger

func init() {
	flag.StringVar(&configFile, "c", "config.json", "-c config-file-name")
	flag.Parse()
	buf, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(buf, &config)
	}
	File, err := os.OpenFile(config.Log, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		println("usage default out.")
		Log = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		Log = log.New(File, "", log.LstdFlags)
	}
	if err != nil {
		Log.Println("init config faild:", err)
		Log.Println("usage default:", config)
	} else {
		Log.Println("config:", config)
	}
	if !strings.HasSuffix(config.Path, "/") {
		config.Path += "/"
	}
}

const (
	unDefindError  = 500
	notFoundFile   = 600
	marshalError   = 601
	unmarshalError = 602
	readBodyError  = 603
	WhitelistError = 604
)

var errorCodeMap map[int]string = map[int]string{unDefindError: "undefind error",
	notFoundFile:   "can't found file",
	marshalError:   "marshal error",
	readBodyError:  "read body error",
	WhitelistError: "remoteaddr not in white list"}

func ErrorCode(code int) string {
	str, ok := errorCodeMap[code]
	if ok {
		return str
	}
	return errorCodeMap[unDefindError]
}

func client(j string) (int, error) {
	resp, err := http.Get(config.Host + "?JobId=" + j + "&" + "Port=" + config.Port)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func size(path string, l []string) int64 {
	var s int64
	for _, file := range l {
		info, err := os.Lstat(path + file)
		if err != nil {
			continue
		}
		s += info.Size()
	}
	return s
}

func walkDir(path string) ([]string, int64) {
	path = config.Path + path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	var l []string
	var s int64
	filepath.Walk(path, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		root = strings.TrimPrefix(filepath.ToSlash(root), path)
		l = append(l, root)
		s += info.Size()
		return nil
	})
	return l, s
}

func SplitPath(path string) (string, string) {
	path = filepath.ToSlash(path)
	if filepath.IsAbs(path) {
		list := strings.Split(path, "/")
		if len(list) <= 1 {
			return "", ""
		}
		path = strings.Join(list[1:], "/")
	}
	dir := filepath.ToSlash(filepath.Dir(path))
	base := filepath.ToSlash(filepath.Base(path))
	if dir == "." || dir == base {
		dir = ""
	}
	return dir, base
}
