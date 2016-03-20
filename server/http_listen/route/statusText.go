package route

const (
	StatusDefaultError   = 1
	StatusAuthFaild      = 600
	StatusReadBodyError  = 601
	StatusPostArgsError  = 602
	StatusUnmarshalError = 603
)

var statusText = map[int]string{
	StatusDefaultError:   "Unknow Error",
	StatusAuthFaild:      "Authentication Failed!",
	StatusReadBodyError:  "Read Body Error",
	StatusPostArgsError:  "Post Args Error",
	StatusUnmarshalError: "Msg Unmarshal Error",
}

func StatusText(code int) string {
	str, ok := statusText[code]
	if ok {
		return str
	}
	return statusText[StatusDefaultError]
}
