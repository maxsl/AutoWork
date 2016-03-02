package encode

import "github.com/axgle/mahonia"

func GBKtoUTF8() mahonia.Decoder {
	str := mahonia.NewDecoder("gb18030")
	return str
}
