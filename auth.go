package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//if not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		fmt.Println(err.Error())
	} else {
		//if success
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("Error when trying to get begin auth URL", provider, "-", err)
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		// log.Println("TODO handle login for", provider)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Authorized action %s not supported", action)
	}
}
