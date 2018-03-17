package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/konojunya/cli-oauth/auth"
)

type Client struct {
	credentials *oauth.Credentials
}

var (
	oauthClient *oauth.Client
	api         *Client
)

func init() {
	oauthClient = auth.GetOauthClient()
}

func SetAPI(at *oauth.Credentials) {
	api = &Client{
		credentials: at,
	}
}

func SetUpClient(oauthToken, oauthVerifier string) error {
	at, err := auth.GetAccessToken(&oauth.Credentials{
		Token: oauthToken,
	}, oauthVerifier)
	if err != nil {
		return err
	}

	outputJSON, err := json.Marshal(&at)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(".token.json", []byte(outputJSON), os.ModePerm)

	SetAPI(at)

	return nil
}

func Tweet(text string) error {
	if api == nil {
		return fmt.Errorf("api is not authorized")
	}
	return api.Tweet(text)
}
