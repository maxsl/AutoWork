package file

import "io"

func SendFile(con io.Writer, sFile io.Reader) (bool, error) {
	buf := make([]byte, 1024)
	for {
		n, err := sFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return false, err
			}
			break
		}
		_, err = con.Write(buf[:n])
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
