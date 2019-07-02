package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"./trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("public", t.filename)))
	})
	t.templ.Execute(w, r)

}

func main() {
	fmt.Println("---Web Chat Application---")
	var addr = flag.String("addr", ":8999", "The address of the application.")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()

	//start web server
	log.Println("Starting web server on :", *addr)
	err := http.ListenAndServe(":8999", nil)
	if err != nil {
		fmt.Printf("Erorr nih")
	}
}
