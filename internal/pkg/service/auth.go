package service

import (
	"fmt"
	"net/http"
)

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked Login")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}

func (s *Service) Signup(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked Signup")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked Logout")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}
