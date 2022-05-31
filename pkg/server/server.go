package server

import (
	"fmt"
	"net/http"

	diag "github.com/juazasan/opentelemetry-go-basics/pkg/diagnostics"
)

type Server struct {
	port    string
	router  *http.ServeMux
	counter int
}

func NewServer() Server {
	r := http.NewServeMux()
	newServer := Server{
		port:    ":9000",
		router:  r,
		counter: 0,
	}
	r.HandleFunc("/", newServer.getHandler())
	return newServer
}

func (s *Server) StartNotBlocking() {
	go func() {
		http.ListenAndServe(s.port, diag.HTTPTraceMiddleware(s.router))
	}()
}

func (s *Server) getHandler() func(w http.ResponseWriter, req *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		s.doStuff()
		response := fmt.Sprintf("Hello World count %d", s.counter)
		w.Write([]byte(response))
	}
}

func (s *Server) doStuff() {
	s.counter++
}
