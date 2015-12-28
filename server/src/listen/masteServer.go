package listen

import (
	"alltype"
	"fmt"
	"net"
	"os"
	ftime "public_tool/time"
)

func init() {
	lis, err := net.Listen("tcp", alltype.Config.IP+":"+alltype.Config.Port)
	if err != nil {
		fmt.Println(ftime.LineTime(), " ", err)
		os.Exit(1)
	}
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		go HandleMessage(con)
	}
}

func HandleMessage(con net.Conn) {

}
func ReadMessage(con net.Conn) ([]byte, error) {

}
