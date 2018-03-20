package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/garyburd/go-oauth/oauth"

	"github.com/konojunya/cli-oauth/server"
	"github.com/konojunya/cli-oauth/twitter"
	"github.com/urfave/cli"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Tweet(c *cli.Context) {
	found := exists(".token.json")
	if !found {
		fmt.Println("you are authorize yet :)\n$ go run main.go login")
		return
	}
	loadToken()
	err := twitter.Tweet("portを自由に変えれるのかテスト")
	if err != nil {
		log.Fatal(err)
	}
}

func loadToken() {
	raw, err := ioutil.ReadFile(".token.json")
	if err != nil {
		log.Fatal(err)
	}

	var credentials *oauth.Credentials
	err = json.Unmarshal(raw, &credentials)
	if err != nil {
		log.Fatal(err)
	}

	twitter.SetAPI(credentials)
}

func Login(c *cli.Context) {
	server.Listen()
}
