package server

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/skratchdot/open-golang/open"

	"github.com/konojunya/cli-oauth/auth"
	"github.com/konojunya/cli-oauth/twitter"
)

var (
	closeCh = make(chan bool, 1)
	addr    string
)

func init() {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	addr = l.Addr().String()
	l.Close()
}

func getRedirectURL() string {
	config := auth.GetOauthClient()
	rt, err := config.RequestTemporaryCredentials(nil, "http://"+addr+"/oauth", nil)
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
	srv := http.Server{Addr: addr}
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/oauth", callbackHandler)
	go func() {
		open.Run("http://" + addr)
		log.Println("listen and serve http://" + addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
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
