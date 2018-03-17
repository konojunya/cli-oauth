package server

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/konojunya/cli-oauth/auth"
	"github.com/konojunya/cli-oauth/twitter"
)

var closeCh = make(chan bool, 1)

func getRedirectURL() string {
	config := auth.GetOauthClient()
	rt, err := config.RequestTemporaryCredentials(nil, "http://127.0.0.1:3000/oauth", nil)
	if err != nil {
		panic(err)
	}

	url := config.AuthorizationURL(rt, nil)
	return url
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	oauthToken := r.URL.Query().Get("oauth_token")
	oauthVerifier := r.URL.Query().Get("oauth_verifier")

	closeCh <- true

	err := twitter.SetUpClient(oauthToken, oauthVerifier)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Listen() {
	srv := &http.Server{Addr: ":3000"}
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/oauth", callbackHandler)
	go func() {
		log.Println("listen and serve http://localhost:3000")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := <-closeCh
	if stop {
		close(closeCh)
		time.Sleep(time.Second * 5)
		if err := srv.Shutdown(nil); err != nil {
			log.Print(err)
		}
	}
}
