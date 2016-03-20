package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func receive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, StatusText(StatusReadBodyError), StatusReadBodyError)
		return
	}
	var rst JobResult
	err = json.Unmarshal(buf, &rst)
	if err != nil {
		http.Error(w, StatusText(StatusUnmarshalError), StatusUnmarshalError)
		return
	}
	// 可以在这里处理上报的执行结果.
}
