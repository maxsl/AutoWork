package intTobyte

func ToInt(b []byte) int {
	var intIP int
	for k, v := range b {
		intIP = intIP | int(v)<<uint(8*(3-k))
	}
	return intIP
}
func ToByte(num int) []byte {
	if num > 4294967295 {
		return []byte{}
	}
	var buf []byte = make([]byte, 4)
	buf[3] = byte(num & 255)
	buf[2] = byte(num >> 8 & 255)
	buf[1] = byte(num >> 16 & 255)
	buf[0] = byte(num >> 24 & 255)
	return buf
}
