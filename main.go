package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/StinkyPeach/bridge/adapter/outbound"
	C "github.com/StinkyPeach/bridge/constant"
	"github.com/jessevdk/go-flags"
	"github.com/nadoo/glider/proxy"
	"github.com/nadoo/glider/proxy/socks5"
)

type Options struct {
	ServerAddr    string `short:"s" description:"server address" required:"true"`
	ServerPort    int    `short:"p" description:"server port" required:"false" default:"8388"`
	LocalAddr     string `short:"b" description:"local binding address" required:"false" default:"127.0.0.1"`
	LocalPort     int    `short:"l" description:"local port" required:"false" default:"1080"`
	Password      string `short:"k" description:"password" required:"true"`
	Method        string `short:"m" description:"encryption method" required:"false" default:"aes-256-cfb"`
	Obfs          string `short:"o" description:"obfsplugin" required:"false" default:"http_simple"`
	ObfsParam     string `long:"op" description:"obfs param" required:"false"`
	Protocol      string `short:"O" description:"protocol" required:"false" default:"origin"`
	ProtocolParam string `long:"Op" description:"protocol param" required:"false"`

	Dns string `long:"dns" description:"custom dns" required:"false" default:"8.8.8.8:53"`
}

var opts Options
var r *net.Resolver

type SSRProxy struct {
	dialer *ProxyDialer
}

type ProxyDialer struct {
	proxy C.Proxy
}

type UdpConn struct {
	net.PacketConn
	addr *net.UDPAddr
}

func (u UdpConn) Read(b []byte) (int, error) {
	n, _, err := u.ReadFrom(b)
	return n, err
}

func (u UdpConn) Write(b []byte) (int, error) {
	return u.WriteTo(b, u.RemoteAddr())
}

func (u UdpConn) RemoteAddr() net.Addr {
	return u.addr
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
		ctx := context.Background()
		ips, err := r.LookupIP(ctx, "ip", address)
		if err != nil {
			return nil, err
		}
		ip = ips[0]
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
		ctx := context.Background()
		ips, err := r.LookupIP(ctx, "ip", address)
		if err != nil {
			return nil, nil, err
		}
		ip = ips[0]
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

func (p SSRProxy) Dial(network, addr string) (net.Conn, proxy.Dialer, error) {
	conn, err := p.dialer.Dial(network, addr)
	return conn, p.dialer, err
}

func (p SSRProxy) DialUDP(network, addr string) (net.PacketConn, proxy.UDPDialer, net.Addr, error) {
	conn, ad, err := p.dialer.DialUDP(network, addr)
	return conn, p.dialer, ad, err
}

func (p SSRProxy) NextDialer(dstAddr string) proxy.Dialer {
	log.Printf("NextDialer: dstAddr: %s\n", dstAddr)
	return p.dialer
}

func (p SSRProxy) Record(dialer proxy.Dialer, success bool) {
	// log.Printf("Record: success: %v\n", success)
}

func StartSocksServer(addr string, proxy proxy.Proxy) {
	server, _ := socks5.NewSocks5Server(addr, SSRProxy{})
	server.ListenAndServe()
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	parser.Usage = "[OPTIONS] address"

	_, err := parser.Parse()

	if err != nil {
		log.Fatalln(err)
	}

	c := strings.LastIndex(opts.Dns, ":")
	if c == -1 {
		log.Fatalln("bad dns address")
	}
	dnsAddr := opts.Dns[:c]
	dp := opts.Dns[c+1:]
	dnsPort, err := strconv.Atoi(dp)

	dnsIp := net.ParseIP(dnsAddr)

	if dnsIp == nil {
		log.Fatalln("bad dns ip")
	}

	if err != nil {
		log.Fatalln("bad dns port")
	}

	ssrNode := make(map[string]interface{})
	ssrNode["name"] = "ssr"
	ssrNode["type"] = "ssr"
	ssrNode["server"] = opts.ServerAddr
	ssrNode["port"] = opts.ServerPort
	ssrNode["password"] = opts.Password
	ssrNode["cipher"] = opts.Method
	ssrNode["obfs"] = opts.Obfs
	ssrNode["obfs-param"] = opts.ObfsParam
	ssrNode["protocol"] = opts.Protocol
	ssrNode["protocol-param"] = opts.ProtocolParam
	ssrNode["udp"] = true

	p, _ := outbound.ParseProxy(ssrNode)

	pr := &ProxyDialer{
		proxy: p,
	}

	server, _ := socks5.NewSocks5Server(fmt.Sprintf("socks://%s:%d", opts.LocalAddr, opts.LocalPort), SSRProxy{
		dialer: pr,
	})

	r = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			md := &C.Metadata{
				NetWork: C.UDP,
				DstIP:   dnsIp,
				DstPort: dp,
			}

			pk, err := pr.proxy.DialUDP(md)

			if err != nil {
				fmt.Println(err)
			}

			addr := &net.UDPAddr{
				IP:   md.DstIP,
				Port: dnsPort,
			}

			return UdpConn{
				pk,
				addr,
			}, err

		},
	}

	server.ListenAndServe()
}
