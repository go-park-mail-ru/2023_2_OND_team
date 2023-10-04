package main

import (
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("upload"))

	s := &http.Server{
		Addr:    ":8081",
		Handler: http.StripPrefix("/upload/", fs),
	}
	s.ListenAndServe()
}
