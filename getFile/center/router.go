package center

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func index(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Printf("index read sourceJob faild:", err)
		w.Write([]byte("index read sourceJob faild: " + err.Error()))
		return
	}
	var SJ SourceJob
	err = json.Unmarshal(buf, &SJ)
	if err != nil {
		Log.Println("Unmarshal sourceJob faild:", err)
		w.Write([]byte("Unmarshal sourceJob faild: " + err.Error()))
		return
	}
	if config.Debug {
		Log.Println(SJ)
	}
	if len(SJ.IPlist) <= 0 {
		w.Write([]byte("IPlist can not null."))
		return
	}
	SJ.Job.JobId = getJobId()
	buf, _ = json.Marshal(SJ.Job)
	result := forwardrequst(SJ.IPlist, buf)
	Log.Println(result)
	buf, err = json.Marshal(result)
	if err != nil {
		Log.Println("Marshal result faild:", err)
		http.Error(w, "Marshal result faild", 601)
		return
	}
	w.Write(buf)
}

func run(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Println("Run read sourceJob faild:", err)
		http.NotFound(w, r)
		return
	}
	m := make(map[string]Exec)
	err = json.Unmarshal(buf, &m)
	if err != nil {
		Log.Println("Run unmarshal msg faild:", err)
		http.NotFound(w, r)
		return
	}
	request(m)
	for _, v := range m {
		if len(v.JobId) > 0 {
			w.Write([]byte("http://" + r.Host + "/getFile/download?&JobId=" + v.JobId))
			break
		}
	}
}

func finished(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ipl := strings.Split(r.RemoteAddr, ":")
	if len(ipl) != 2 {
		Log.Println("Finished Func:", r.RemoteAddr)
		return
	}
	ip, ok := ServerRelationship[ipl[0]]
	if !ok {
		ip = ipl[0]
	}
	port := r.FormValue("Port")
	id := r.FormValue("JobId")
	os.MkdirAll(config.TmpDir+id, 0644)
	u := "http://" + ipl[0] + ":" + port + "/getFile/download?JobId=" + id
	if config.Debug {
		Log.Println("Request:", u)
	}
	resp, err := http.Get(u)
	if err != nil {
		Log.Println("Finished Func:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Log.Println("Finished Func: error code", resp.StatusCode)
		return
	}
	File, err := os.Create(config.TmpDir + id + "/" + ip + ".zip")
	if err != nil {
		Log.Println("Finished Func:", err)
		return
	}
	io.Copy(File, resp.Body)
	File.Close()
	sendResult(id)
}

func download(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.FormValue("JobId")
	file := r.FormValue("file")
	if len(id) > 0 {
		var list []string
		if len(file) <= 0 {
			l, err := ioutil.ReadDir(config.TmpDir + id)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			for _, v := range l {
				list = append(list, v.Name())
			}
			writeIndex(w, id, list)
		} else {
			http.ServeFile(w, r, config.TmpDir+id+"/"+file)
		}
	} else {
		http.NotFound(w, r)
	}
}

func writeIndex(w io.Writer, JobId string, list []string) {
	w.Write([]byte("<html><pre>"))
	for _, v := range list {
		str := "<a href=download?JobId=" + JobId + "&file=" + v + ">" + v + "</a><br>"
		w.Write([]byte(str))
	}
	w.Write([]byte("</pre></html>"))
}

type msg struct {
	mode     string `json:mode`
	contacts string `json:contacts`
	ip       string `json:ip`
}
type sendmsg struct {
	lock *sync.RWMutex
	m    map[string]msg
}

var Contacts *sendmsg = &sendmsg{new(sync.RWMutex), make(map[string]msg)}

func (self *sendmsg) add(id string, m msg) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.m[id] = m
}
func (self *sendmsg) del(id string) (msg, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	m, ok := self.m[id]
	delete(self.m, id)
	return m, ok
}

func contacts(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	mode := r.FormValue("mode")
	cts := r.FormValue("cts")
	if len(id) == 0 || len(cts) == 0 {
		http.Error(w, "Jobid mode or contacts can't empty.", 605)
		return
	}
	Contacts.add(id, msg{mode, cts, r.Host})
}
