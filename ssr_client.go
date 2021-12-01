package main

import (
	"fmt"
	"net"
	"time"

	"github.com/doorbash/bridge/adapter/outbound"
	C "github.com/doorbash/bridge/constant"
	"github.com/nadoo/glider/log"
	"github.com/nadoo/glider/proxy"
	"github.com/nadoo/glider/proxy/mixed"
)

type SSRClient struct {
	mixedServer proxy.Server
}

func (s *SSRClient) ListenAndServe() {
	s.mixedServer.ListenAndServe()
}

func NewSSRClient(
	serverAddr string,
	serverPort int,
	localAddr string,
	localPort int,
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

	client := &SSRClient{}

	client.mixedServer, err = mixed.NewMixedServer(fmt.Sprintf("mixed://%s:%d", localAddr, localPort), &SSRProxy{
		dialer: pr,
	})

	return client, err
}
