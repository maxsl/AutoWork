package tcp_listen

import (
	"fmt"
	"net"
	"testing"
)

func Test_client(t *testing.T) {
	go t_server()

	con, err := ClientTls("127.0.0.1:1789", "../certs/client.pem", "../certs/client.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	_, err = con.Write([]byte("Hello World"))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 20)
	n, err := con.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf[:n]))
}
func Benchmark_cleint(b *testing.B) {
	con, err := ClientTls("127.0.0.1:1789", "../certs/client.pem", "../certs/client.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for i := 0; i < b.N; i++ {
		_, err := con.Write([]byte("Hello World"))
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, 20)
		_, err = con.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(buf[:n]))
	}
}
func t_server() {
	lis, err := Servertls(":1789", "../certs/server.pem", "../certs/server.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		handler(con)
	}
}

func handler(con net.Conn) {
	defer con.Close()
	for {
		buf := make([]byte, 20)
		_, err := con.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(buf[:n]))
		con.Write(buf)
	}
}
