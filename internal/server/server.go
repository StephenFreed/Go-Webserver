package server

import (
    "net/http"
)


type Server struct {
    listenAddr string
}

func NewServer(listenAddr string) *Server {
    return &Server{
        listenAddr: listenAddr,
    }
}

//Server start method listens on Server listenAddr port.
//All paths are passed to the custom router.
func (s *Server) Start() error {
    http.HandleFunc("/", s.router)
    return http.ListenAndServe(s.listenAddr, nil)
}
