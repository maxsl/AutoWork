package http_listen

import (
	"net/http"

	"github.com/AutoWork/server/http_listen/route/handler"
)

func RunHttpServer(ip string) {
	http.HandleFunc("/receive", handler.Receive)
	http.HandleFunc("/run", handler.Run)
	err := http.ListenAndServe(ip, nil)
	if err != nil {
		println(err.Error())
	}
}
