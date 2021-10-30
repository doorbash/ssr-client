package main

import (
	"fmt"
	"h12.io/socks"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
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
	Dns           string `long:"dns" description:"custom dns" required:"false" default:"8.8.8.8:53"`
	LocalHttpPort int    `short:"r" description:"http relay port" required:"false" default:"0"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()

	if err != nil {
		os.Exit(1)
	}

	ssrClient, err := NewSSRClient(
		opts.ServerAddr,
		opts.ServerPort,
		opts.LocalAddr,
		opts.LocalPort,
		opts.Password,
		opts.Method,
		opts.Obfs,
		opts.ObfsParam,
		opts.Protocol,
		opts.ProtocolParam,
		opts.Dns,
	)

	if err != nil {
		log.Fatalln(err)
	}
	if opts.LocalHttpPort != 0 {
		httpHandler := HttpHandler{
			SocksAddr:  fmt.Sprintf("127.0.0.1:%d",opts.LocalPort),
			SocksProto: socks.SOCKS5,
		}
		go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d",opts.LocalHttpPort), &httpHandler)
	}
	ssrClient.ListenAndServe()
}
