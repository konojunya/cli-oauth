package twitter

import (
	"net/url"
)

func (api *Client) Tweet(text string) error {
	v := url.Values{}
	v.Set("status", text)
	res, err := oauthClient.Post(nil, api.credentials, "https://api.twitter.com/1.1/statuses/update.json", v)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
