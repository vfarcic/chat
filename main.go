package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	"github.com/stretchr/objx"
	phttp "github.com/pikanezi/http"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application. Default value is ':8080'")
	flag.Parse()
	r := phttp.NewRouter()
	r.SetCustomHeader(phttp.Header{
		"Access-Control-Allow-Origin": "*",
	})
	room := newRoom()
	r.Handle("/test", &templateHandler{filename: "test.html"})
	r.Handle("/room", room)
	r.PathPrefix("/bower_components/").Handler(
		http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components"))))
	r.PathPrefix("/components/").Handler(
		http.StripPrefix("/components/", http.FileServer(http.Dir("components"))))
	go room.run()
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Println("Starting the server on", *addr)

}

