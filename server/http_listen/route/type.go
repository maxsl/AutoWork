package route

type log interface {
	PrintfI(formate string, v ...interface{})
	PrintfW(formate string, v ...interface{})
	PrintfE(formate string, v ...interface{})
	PrintfF(formate string, v ...interface{})
}
