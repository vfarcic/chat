package main

// TODO: Change hardcoded auth URLs in login.html

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	"github.com/stretchr/objx"
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
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", logoutHandler)
	// TODO: Remove
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle(
		"/bower_components/",
		http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components"))))
	http.Handle(
		"/components/",
		http.StripPrefix("/components/", http.FileServer(http.Dir("components"))))
	go r.run()
	log.Println("Starting the server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

