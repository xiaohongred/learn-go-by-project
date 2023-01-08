package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

func Address(s *simpleServer) string {
	return s.addr
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
	return
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port           string
	rounRobinCount int
	servers        []Server
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.rounRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.rounRobinCount++
		server = lb.servers[lb.rounRobinCount%len(lb.servers)]
	}
	lb.rounRobinCount++
	return server
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:           port,
		servers:        servers,
		rounRobinCount: 0,
	}
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.163.com/"),
		newSimpleServer("http://www.bing.com"),
		newSimpleServer("http://baidu.com"),
	}

	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("serving requests at localhost:%s\n", lb.port)
	err := http.ListenAndServe(":"+lb.port, nil)
	log.Fatal(err)
}
