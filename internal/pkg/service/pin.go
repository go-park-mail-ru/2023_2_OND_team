package service

import (
	"fmt"
	"net/http"
)

func (s *Service) GetPins(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked GetPins")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}

func (s *Service) GetPinByID(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked GetPinByID")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}
