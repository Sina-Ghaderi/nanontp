package engine

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"
)

type Transport interface {
	Write(data string, addr string, port int)
}

type udpTransport struct {
	conn *net.UDPConn
}

func (p *udpTransport) setConn(conn *net.UDPConn) {
	p.conn = conn
}

func (p *udpTransport) Write(data string, addr string, port int) {
	laddr, err := net.ResolveUDPAddr("udp", addr+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("resolve addr err: ", err)
		return
	}
	_, er := p.conn.WriteTo([]byte(data), laddr)
	if er != nil {
		fmt.Println("resolve addr err: ", err)
		return
	}
}

type _reactor struct {
	udp_listeners map[int]UdpClient
	udp_conn      map[int]*net.UDPConn
	timer         []*LaterCalling
}

func (p *_reactor) ListenUdp(addaport string, udp UdpClient) {
	laddr, err := net.ResolveUDPAddr("udp", addaport)
	if err == nil {
		p.listenUdp(laddr, udp)
	} else {
		log.Println("resolve addr err: ", err)
		return
	}
}

type UdpClient interface {
	DatagramReceived(data []byte, addr net.Addr)
	SetUdpTransport(Transport)
}

func (p *_reactor) listenUdp(addr *net.UDPAddr, udp UdpClient) {
	p.initReactor()
	p.udp_listeners[addr.Port] = udp
	log.Println("ntp server listening on (UDP)", addr.String())
	c, erl := net.ListenUDP("udp", addr)
	if erl != nil {
		log.Fatal("listen error: ", erl)
	} else {
		p.udp_conn[addr.Port] = c
	}
	transport := new(udpTransport)
	transport.setConn(c)
	udp.SetUdpTransport(transport)
}

func (p *_reactor) initReactor() {
	if p.udp_listeners == nil {
		p.udp_listeners = make(map[int]UdpClient)
	}
	if p.udp_conn == nil {
		p.udp_conn = make(map[int]*net.UDPConn)
	}
}

var (
	Reactor        = new(_reactor)
	listening_chan chan int
)

type LaterCalling struct {
	millisecond int
	call        func()
}

type reactor interface {
	ListenUdp(port int, client UdpClient)
	CallLater(microsecond int, latercaller func())
	Run()
}

func (p *_reactor) CallLater(millisecond int, lc func()) {
	calling := new(LaterCalling)
	calling.millisecond = millisecond
	calling.call = lc
	p.timer = append(p.timer, calling)
}

func (p *_reactor) Run() {
	runtime.GOMAXPROCS(len(p.udp_conn))
	for port, l := range p.udp_conn {
		go handleUdpConnection(l, p.udp_listeners[port])
	}
	for len(p.timer) > 0 {
		caller := p.timer[0]
		p.timer = p.timer[1:]
		selectTimer(caller)
	}
	for {
		fmt.Println("------------------ Logs and Errors ------------------")
		select {
		case <-listening_chan:
			fmt.Println("-----------------------------------------------------")

		}
	}
}

func selectTimer(caller *LaterCalling) {
	select {
	case <-time.After(time.Duration(caller.millisecond) * time.Millisecond):
		caller.call()
	}
}

func handleUdpConnection(conn *net.UDPConn, client UdpClient) {
	for {
		data := make([]byte, 512)
		read_length, remoteAddr, err := conn.ReadFromUDP(data[0:])
		if err != nil {
			return
		} else {
		}
		if read_length > 0 {
			go panicWrapping(func() {
				client.DatagramReceived(data[0:read_length], remoteAddr)
			})
		}
	}
}

func panicWrapping(f func()) {
	defer func() {
		recover()
	}()
	f()
}

type UdpHandler struct {
	udptransport Transport
}

func (p *UdpHandler) SetUdpTransport(transport Transport) {
	p.udptransport = transport
}

func (p *UdpHandler) UdpWrite(data string, addr string, port int) {
	p.udptransport.Write(data, addr, port)
}
