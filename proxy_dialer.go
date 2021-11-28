package main

import (
	"log"
	"net"
	"strconv"

	C "github.com/doorbash/bridge/constant"
)

type ProxyDialer struct {
	proxy C.Proxy
}

func (p *ProxyDialer) Dial(network, addr string) (c net.Conn, err error) {
	log.Printf("Dial: network: %s, addr: %s\n", network, addr)

	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		return nil, err
	}

	metadata := &C.Metadata{
		NetWork: C.TCP,
		DstPort: port,
	}

	ip := net.ParseIP(host)

	if ip == nil {
		metadata.Host = host
		metadata.AddrType = C.ATypDomainName
	} else {
		ipv4 := ip.To4()
		if ipv4 != nil {
			metadata.DstIP = ipv4
			metadata.AddrType = C.ATypIPv4
		} else {
			metadata.DstIP = ip
			metadata.AddrType = C.ATypIPv6
		}
	}

	_, err = strconv.Atoi(port)

	if err != nil {
		return nil, err
	}

	return p.proxy.Dial(metadata)
}

func (p *ProxyDialer) DialUDP(network, addr string) (pc net.PacketConn, writeTo net.Addr, err error) {
	log.Printf("DialUDP: network: %s, addr: %s\n", network, addr)

	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		return nil, nil, err
	}

	metadata := &C.Metadata{
		NetWork: C.UDP,
		DstPort: port,
	}

	ip := net.ParseIP(host)

	if ip == nil {
		metadata.Host = host
		metadata.AddrType = C.ATypDomainName
	} else {
		ipv4 := ip.To4()
		if ipv4 != nil {
			metadata.DstIP = ipv4
			metadata.AddrType = C.ATypIPv4
		} else {
			metadata.DstIP = ip
			metadata.AddrType = C.ATypIPv6
		}
	}

	prt, err := strconv.Atoi(port)

	if err != nil {
		return nil, nil, err
	}

	uaddr := &net.UDPAddr{
		IP:   metadata.DstIP,
		Port: prt,
	}

	pc, e := p.proxy.DialUDP(metadata)

	return pc, uaddr, e
}

func (p *ProxyDialer) Addr() string {
	return ""
}

func NewProxyDialer(p C.Proxy) (*ProxyDialer, error) {
	return &ProxyDialer{
		proxy: p,
	}, nil
}
