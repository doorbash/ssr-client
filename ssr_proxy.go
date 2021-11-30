package main

import (
	"net"

	"github.com/nadoo/glider/proxy"
)

type SSRProxy struct {
	dialer *ProxyDialer
}

func (p SSRProxy) Dial(network, addr string) (net.Conn, proxy.Dialer, error) {
	conn, err := p.dialer.Dial(network, addr)
	return conn, p.dialer, err
}

func (p SSRProxy) DialUDP(network, addr string) (net.PacketConn, proxy.UDPDialer, net.Addr, error) {
	conn, ad, err := p.dialer.DialUDP(network, addr)
	return conn, p.dialer, ad, err
}

func (p SSRProxy) NextDialer(dstAddr string) proxy.Dialer {
	return p.dialer
}

func (p SSRProxy) Record(dialer proxy.Dialer, success bool) {
}
