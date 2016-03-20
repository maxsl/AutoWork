package route

const (
	StatusDefaultError  = 1
	StatusAuthFaild     = 600
	StatusUnknowAction  = 601
	StatusPostArgsError = 602
	StatusNotValidFile  = 603
)

var statusText = map[int]string{
	StatusDefaultError:  "Unknow Error",
	StatusAuthFaild:     "Authentication Failed!",
	StatusUnknowAction:  "Unknow Action",
	StatusPostArgsError: "Post Args Error",
	StatusNotValidFile:  "Not Is A Valid Zip File",
}

func StatusText(code int) string {
	str, ok := statusText[code]
	if ok {
		return str
	}
	return statusText[StatusDefaultError]
}
