package main

import (
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

func copyHeader(dst, src http.Header) {
	for headerKey, headerValue := range src {
		for _, headerSegment := range headerValue {
			dst.Add(headerKey, headerSegment)
		}
	}
}
func (s *HttpHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	log.Printf("REQUEST: %s %s", request.Method, request.RequestURI)
	dialer := socks.DialSocksProxy(s.SocksProto, s.SocksAddr)
	if request.Method == "CONNECT" {
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
		outconn, err := dialer("tcp", request.Host)
		if err != nil {
			log.Printf("dial to %s error, %s", request.Host, err)
			return
		}
		go Copy(conn, outconn)
		go Copy(outconn, conn)
		return
	}
	tr := &http.Transport{
		Dial:               dialer,
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	outreq := *request
	outreq.RequestURI = ""
	copyHeader(outreq.Header, request.Header)
	resp, err := client.Do(&outreq)
	if err != nil {
		log.Printf("request socks: %s", err)
		return
	}
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)

}
