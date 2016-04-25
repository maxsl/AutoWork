package tcp_listen

import (
	"crypto/rand"
	"crypto/tls"
	"net"
)

func Server(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func Servertls(addr, crt, key string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", addr, &config)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
