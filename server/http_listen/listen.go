package http_listen

import (
	"net/http"

	"github.com/AutoWork/server/http_listen/route"
)

func RunHttpServer(ip string) {
	http.HandleFunc("/", route.Router)
	err := http.ListenAndServe(ip, nil)
	if err != nil {
		println(err.Error())
	}
}
