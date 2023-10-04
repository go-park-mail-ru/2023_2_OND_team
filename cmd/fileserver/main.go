package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("upload"))

	s := &http.Server{
		Addr:    ":8081",
		Handler: http.StripPrefix("/upload/", fs),
	}
	log.Println("fileserver start")
	defer log.Panicln("fileserver finish")
	s.ListenAndServeTLS("/home/ond_team/cert/fullchain.pem", "/home/ond_team/cert/privkey.pem")
}
