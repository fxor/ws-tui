package wssv

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	httpServer    *http.Server
	upgrader      websocket.Upgrader
	clientSockets map[string]*ClientConn
}

func NewServer(addr string) *Server {
	sv := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clientSockets: map[string]*ClientConn{},
	}

	r := http.NewServeMux()
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

	})
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		sv.loginHandler(w, r)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sv.upgrade(w, r)
	})
	sv.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return sv
}

func (s *Server) Start() error {
	fmt.Printf("Server listening at %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
