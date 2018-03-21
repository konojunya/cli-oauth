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
	closeCh  = make(chan bool, 1)
	listener net.Listener
)

func getRedirectURL() string {
	config := auth.GetOauthClient()
	rt, err := config.RequestTemporaryCredentials(nil, "http://"+listener.Addr().String()+"/oauth", nil)
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

func tcpListen() {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	listener = l
}

func Listen() {
	tcpListen()
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/oauth", callbackHandler)
	go func() {
		open.Run("http://" + listener.Addr().String())
		log.Println("listen and serve http://" + listener.Addr().String())
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal(err)
		}
	}()

	stop := <-closeCh
	if stop {
		close(closeCh)
		time.Sleep(time.Second * 3)
		if err := listener.Close(); err != nil {
			log.Print(err)
		}
	}
}
