package main

import (
	"fmt"
	"net"
	"time"

	"github.com/doorbash/bridge/adapter/outbound"
	C "github.com/doorbash/bridge/constant"
	"github.com/nadoo/glider/log"
	"github.com/nadoo/glider/proxy"
	"github.com/nadoo/glider/proxy/http"
	"github.com/nadoo/glider/proxy/socks5"
)

type SSRClient struct {
	httpServer   proxy.Server
	socks5server proxy.Server
}

func (s *SSRClient) ListenAndServe() {
	go s.httpServer.ListenAndServe()
	s.socks5server.ListenAndServe()
}

func NewSSRClient(
	serverAddr string,
	serverPort int,
	localAddr string,
	localSocksPort int,
	localHttpPort int,
	password string,
	method string,
	obfs string,
	obfsParam string,
	protocol string,
	protocolParam string,
	forwardProxy string,
	socketTimeout time.Duration,
	verboseMode bool,
) (*SSRClient, error) {
	if verboseMode {
		log.F = log.Debugf
	}

	ssrNode := map[string]interface{}{
		"name":           "ssr",
		"type":           "ssr",
		"server":         serverAddr,
		"port":           serverPort,
		"password":       password,
		"cipher":         method,
		"obfs":           obfs,
		"obfs-param":     obfsParam,
		"protocol":       protocol,
		"protocol-param": protocolParam,
		"udp":            true,
	}

	p, _ := outbound.ParseProxy(ssrNode)

	var ps C.Proxy
	if forwardProxy != "" {
		host, port, err := net.SplitHostPort(forwardProxy)
		if err != nil {
			return nil, err
		}
		socks5Node := make(map[string]interface{})
		socks5Node["name"] = "socks"
		socks5Node["type"] = "socks5"
		socks5Node["server"] = host
		socks5Node["port"] = port
		socks5Node["udp"] = true
		socks5Node["skip-cert-verify"] = true

		ps, _ = outbound.ParseProxy(socks5Node)

		p.SetDialer(ps)
	}

	pr, err := NewProxyDialer(p, socketTimeout)

	if err != nil {
		return nil, err
	}

	ssrProxy := &SSRProxy{
		dialer: pr,
	}

	client := &SSRClient{}

	client.socks5server, _ = socks5.NewSocks5Server(fmt.Sprintf("socks://%s:%d", localAddr, localSocksPort), ssrProxy)

	client.httpServer, _ = http.NewHTTPServer(fmt.Sprintf("http://%s:%d", localAddr, localHttpPort), ssrProxy)

	return client, nil
}
