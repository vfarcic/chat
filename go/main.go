package main

import (
	"log"
	"net/http"
	"flag"
	phttp "github.com/pikanezi/http"
)

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application. Default value is ':8080'")
	flag.Parse()
	r := phttp.NewRouter()
	r.SetCustomHeader(phttp.Header{
		"Access-Control-Allow-Origin": "*",
	})
	room := newRoom()
	r.Handle("/room", room)
	r.PathPrefix("/bower_components/").Handler(
		http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components"))))
	r.PathPrefix("/components/").Handler(
		http.StripPrefix("/components/", http.FileServer(http.Dir("components"))))
	go room.Run()
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Println("Starting the server on", *addr)

}

