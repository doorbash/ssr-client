package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/nadoo/glider/proxy"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type HttpHandler struct {
	addr  string
	proxy proxy.Proxy
}

func copy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("http_handler: copy connection error, %s", err)
	}
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
	delHopHeaders(r.Header)

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		appendHostToXForwardHeader(r.Header, clientIP)
	}
	resp, err := client.Do(r)
	if err != nil {
		http.Error(rw, "Server Error", http.StatusInternalServerError)
		log.Printf("http_handler: request: %s", err)
		return
	}
	defer resp.Body.Close()
	delHopHeaders(resp.Header)

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
