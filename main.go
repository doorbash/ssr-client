package main

import (
	"log"

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

	Dns string `long:"dns" description:"custom dns" required:"false" default:"8.8.8.8:53"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	parser.Usage = "[OPTIONS] address"

	_, err := parser.Parse()

	if err != nil {
		log.Fatalln(err)
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

	ssrClient.ListenAndServe()
}
