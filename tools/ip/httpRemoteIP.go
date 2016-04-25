package ip

import (
	"net/http"
	"strings"
)

func RemoteIP(r *http.Request) string {
	l := strings.Split(r.RemoteAddr, ":")
	if len(l) != 2 {
		return ""
	}
	return l[0]
}
