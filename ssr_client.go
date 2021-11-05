package main

import (
	"fmt"

	"github.com/kooroshh/bridge/adapter/outbound"
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
	dns string,
) (*SSRClient, error) {
	client := &SSRClient{}

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

	pr, err := NewProxyDialer(p, dns)

	if err != nil {
		return nil, err
	}

	ssrProxy := &SSRProxy{
		dialer: pr,
	}

	client.socks5server, _ = socks5.NewSocks5Server(fmt.Sprintf("socks://%s:%d", localAddr, localSocksPort), ssrProxy)

	client.httpServer, _ = http.NewHTTPServer(fmt.Sprintf("http://%s:%d", localAddr, localHttpPort), ssrProxy)

	return client, nil
}
