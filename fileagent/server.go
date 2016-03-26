package fileagent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var FileInfo chan *FilesInfo = make(chan *FilesInfo, 5)

func backresolve() {
	for {
		info := <-FileInfo
		go func(info *FilesInfo) {
			if info.Host == "" {
				Log.Printf("Backresolve Error: callBack host can't null")
				return
			}
			list := info.Run()
			buf, err := json.Marshal(list)
			if err != nil {
				Log.Printf("Backresolve Error: %v\n", err)
				return
			}
			b := bytes.NewReader(buf)
			//把执行的结果推送给Host.尝试次数为3.
			for i := 0; i < 3; i++ {
				resp, err := http.Post(info.Host, BodyType, b)
				if err == nil && resp.StatusCode == 200 {
					break
				}
				Log.Printf("Backresolve Error: %v\n", err)
			}
		}(info)
	}
}

func Server(ip string) {
	go backresolve()
	http.HandleFunc("/download", download)
	http.HandleFunc("/ack", ack)
	http.HandleFunc("/", index)
	http.ListenAndServe(ip, nil)
}

//index方法是接收客户端发的请求.根据请求判断用户的行为.是查看目录下的文件还是copy打包要下载的
//文件.不同的Tag标签响应的内容不同.
func index(w http.ResponseWriter, r *http.Request) {
	Log.Printf("Access Info: %v\n", r.URL.String())
	defer r.Body.Close()
	var request *Job = new(Job)
	err := resolve(w, r, request)
	if err != nil {
		Log.Printf("Error: index func %v\n", err)
		return
	}
	Log.Printf("Receive Job: %v\n", *request)
	if request.Tag == "printfiles" {
		result := request.PrintFiles()
		Log.Printf("Job:%v run result,%v\n", request.JobId, result)
		buf, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "MarshalError", MarshalError)
			return
		}
		w.Write(buf)
		return
	}
	result := request.GetFilesInfo()
	Log.Printf("Job:%v run result,%v\n", request.JobId, result)
	buf, err := json.Marshal(*result)
	if err != nil {
		http.Error(w, "MarshalError", MarshalError)
		return
	}
	w.Write(buf)
	return
}

//此函数是将Body解析成FilesInfo,主要是为了安全,访问Index函数后响应给用户确认
//然后执行接收文件copy等操作的
func ack(w http.ResponseWriter, r *http.Request) {
	Log.Printf("Access Info: %v\n", r.URL.String())
	defer r.Body.Close()
	var t *FilesInfo = new(FilesInfo)
	err := resolve(w, r, t)
	if err != nil {
		Log.Printf("Error: ack func %v\n", err)
		return
	}
	Log.Printf("Receive FilesInfo: %v\n", *t)
	FileInfo <- t
}

func download(w http.ResponseWriter, r *http.Request) {
	Log.Printf("Access Info: %v\n", r.URL.String())
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !ipIsLanIP(ip) {
		http.Error(w, "ForbiddenAccess", ForbiddenAccess)
		return
	}
	path := r.FormValue("file")
	info, err := os.Stat(path)
	//判断请求的文件路径是不是在临时文件夹下,不是的话就认为是非法请求,返回404
	if err != nil || info.IsDir() || strings.Index(path, tempDir) != 0 {
		http.Error(w, "", 404)
		return
	}
	http.ServeFile(w, r, path)
}

func resolve(w http.ResponseWriter, r *http.Request, t interface{}) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !ipIsLanIP(ip) {
		http.Error(w, "ForbiddenAccess", ForbiddenAccess)
		return err
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ReadBodyError", ReadBodyError)
		return err
	}
	err = json.Unmarshal(buf, t)
	if err != nil {
		http.Error(w, "UnmarshalError", UnmarshalError)
		return err
	}
	return nil
}
