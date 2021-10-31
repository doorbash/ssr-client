package main

import (
	"compress/gzip"
	"h12.io/socks"
	"io"
	"log"
	"net/http"
)

// Borrowed from https://github.com/dworld/http2socks/

type HttpHandler struct {
	SocksAddr  string // socks proxy address
	SocksProto int    // socks proxy protocol type
}

// Copy check return value of io.Copy
func Copy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("copy connection error, %s", err)
	}
}

func (s *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST: %s %s", r.Method, r.RequestURI)
	dialer := socks.DialSocksProxy(s.SocksProto, s.SocksAddr)
	if r.Method == "CONNECT" {
		hj, ok := rw.(http.Hijacker)
		if !ok {
			log.Printf("can't cast to Hijacker")
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			log.Printf("can't hijack the connection")
			return
		}
		_, err = bufrw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			log.Printf("write CONNECT response header error, %s", err)
			return
		}
		err = bufrw.Flush()
		if err != nil {
			log.Printf("flush error, %s", err)
			return
		}
		outconn, err := dialer("tcp", r.Host)
		if err != nil {
			log.Printf("dial to %s error, %s", r.Host, err)
			return
		}
		go Copy(conn, outconn)
		go Copy(outconn, conn)
		return
	}
	tr := &http.Transport{
		Dial: dialer,
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	outreq := *r
	outreq.RequestURI = ""
	resp, err := client.Do(&outreq)
	if err != nil {
		log.Printf("request socks: %s", err)
		return
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	_, err = io.Copy(rw, reader)
	if err != nil {
		log.Printf("write response: %s", err)
		return
	}
}
