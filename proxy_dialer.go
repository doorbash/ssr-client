package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	C "github.com/kooroshh/bridge/constant"
)

type ProxyDialer struct {
	dnsCache *DNSCache
	proxy    C.Proxy
	resolver *net.Resolver
}

func (p *ProxyDialer) resolve(address string) (net.IP, error) {
	ips, err := p.dnsCache.Get(address)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ips, err := p.resolver.LookupIP(ctx, "ip", address)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no address associated with this domain %s", address)
		}
		p.dnsCache.Put(address, ips)
		return ips[0], nil
	}
	ip := ips[0]
	log.Printf("Cache: %s is available in cache : %s\n", address, ip.String())
	return ip, nil
}

func (p *ProxyDialer) Dial(network, addr string) (c net.Conn, err error) {
	log.Printf("Dial: network: %s, addr: %s\n", network, addr)

	colonIndex := strings.LastIndex(addr, ":")

	if colonIndex == -1 {
		return nil, errors.New("bad address")
	}

	address := addr[:colonIndex]
	port := addr[colonIndex+1:]
	ip := net.ParseIP(address)

	if ip == nil {
		if ip, err = p.resolve(address); err != nil {
			return nil, err
		}
	}

	metadata := &C.Metadata{
		NetWork: C.TCP,
		DstPort: port,
	}

	ipv4 := ip.To4()

	if ipv4 != nil {
		metadata.DstIP = ipv4
		metadata.AddrType = C.ATypIPv4
	} else {
		metadata.DstIP = ip
		metadata.AddrType = C.ATypIPv6
	}

	_, err = strconv.Atoi(port)

	if err != nil {
		return nil, err
	}

	conn, err := p.proxy.Dial(metadata)
	return conn, err
}

func (p *ProxyDialer) DialUDP(network, addr string) (pc net.PacketConn, writeTo net.Addr, err error) {
	log.Printf("DialUDP: network: %s, addr: %s\n", network, addr)

	colonIndex := strings.LastIndex(addr, ":")

	if colonIndex == -1 {
		return nil, nil, errors.New("bad address")
	}

	address := addr[:colonIndex]
	port := addr[colonIndex+1:]
	ip := net.ParseIP(address)

	if ip == nil {
		if ip, err = p.resolve(address); err != nil {
			return nil, nil, err
		}
	}

	metadata := &C.Metadata{
		NetWork: C.UDP,
		DstPort: port,
	}

	ipv4 := ip.To4()

	if ipv4 != nil {
		metadata.DstIP = ipv4
		metadata.AddrType = C.ATypIPv4
	} else {
		metadata.DstIP = ip
		metadata.AddrType = C.ATypIPv6
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

func NewProxyDialer(p C.Proxy, Dns string, dnsCache *DNSCache) (*ProxyDialer, error) {
	c := strings.LastIndex(Dns, ":")
	if c == -1 {
		return nil, errors.New("bad dns address")
	}
	dnsAddr := Dns[:c]
	dp := Dns[c+1:]

	dnsIp := net.ParseIP(dnsAddr)
	if dnsIp == nil {
		return nil, errors.New("bad dns ip")
	}

	dnsPort, err := strconv.Atoi(dp)
	if err != nil {
		return nil, errors.New("bad dns port")
	}

	pd := &ProxyDialer{
		proxy:    p,
		dnsCache: dnsCache,
	}

	pd.resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			md := &C.Metadata{
				NetWork: C.UDP,
				DstIP:   dnsIp,
				DstPort: dp,
			}

			pk, err := pd.proxy.DialUDP(md)

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			addr := &net.UDPAddr{
				IP:   md.DstIP,
				Port: dnsPort,
			}

			return UdpConn{
				pk,
				addr,
			}, nil
		},
	}

	return pd, nil
}
