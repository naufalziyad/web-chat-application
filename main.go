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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, r)

}

func main() {
	fmt.Println("---Web Chat Application---")
	var addr = flag.String("addr", ":8999", "The address of the application.")
	flag.Parse()
	//setup gomniauth
	gomniauth.SetSecurityKey("key")
	gomniauth.WithProviders(
		facebook.New("key", "key",
			URLDomain+"/auth/callback/facebook"),
		github.New("key", "",
			URLDomain+"/auth/callback/github"),
		google.New(googleClientID, googleClientSecret,
			URLDomain+"/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.run()

	//start web server
	log.Println("Starting web server on :", *addr)
	err := http.ListenAndServe(":8999", nil)
	if err != nil {
		fmt.Printf("Erorr nih")
	}
}
