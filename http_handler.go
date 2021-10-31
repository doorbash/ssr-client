package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/nadoo/glider/proxy"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("http_handler: copy connection error, %s", err)
	}
}

type HttpHandler struct {
	addr  string
	proxy proxy.Proxy
}

func (s *HttpHandler) dial(network, addr string) (net.Conn, error) {
	c, _, err := s.proxy.Dial(network, addr)
	return c, err
}

func (s *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST: %s %s", r.Method, r.RequestURI)

	if r.Method == "CONNECT" {
		hj, ok := rw.(http.Hijacker)
		if !ok {
			log.Printf("http_handler: can't cast to Hijacker")
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			log.Printf("http_handler: can't hijack the connection")
			return
		}
		_, err = bufrw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			log.Printf("http_handler: write CONNECT response header error, %s", err)
			return
		}
		err = bufrw.Flush()
		if err != nil {
			log.Printf("http_handler: flush error, %s", err)
			return
		}
		outconn, err := s.dial("tcp", r.Host)
		if err != nil {
			log.Printf("http_handler: dial to %s error, %s", r.Host, err)
			return
		}
		go copy(outconn, conn)
		copy(conn, outconn)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: s.dial,
		},
	}
	r.RequestURI = ""
	resp, err := client.Do(r)
	if err != nil {
		http.Error(rw, "Server Error", http.StatusInternalServerError)
		log.Printf("http_handler: request: %s", err)
		return
	}
	defer resp.Body.Close()

	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		log.Printf("http_handler: write response: %s", err)
	}
}

func (s *HttpHandler) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s)
}

func NewHttpHandler(addr string, p proxy.Proxy) *HttpHandler {
	return &HttpHandler{
		addr:  addr,
		proxy: p,
	}
}
