package network

import (
	"net"
	"strconv"
	"strings"

	"github.com/sina-ghaderi/nanontp/engine"
	"github.com/sina-ghaderi/nanontp/getter"
)

var handler *Handler

type Handler struct {
	engine.UdpHandler
}

func GetHandler() *Handler {
	if handler == nil {
		handler = new(Handler)
	}
	return handler
}

var UDPRemoteAddr *net.UDPAddr

func (p *Handler) DatagramReceived(data []byte, addr net.Addr) {
	res, err := getter.Serve(data, addr)
	if err == nil {
		ip, port := spliteAddr(addr.String())
		p.UdpWrite(string(res), ip, port)
	}
}
func spliteAddr(addr string) (string, int) {
	ip := strings.Split(addr, ":")[0]
	port := strings.Split(addr, ":")[1]
	p, _ := strconv.Atoi(port)
	return ip, p
}
