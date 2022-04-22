// Package podio provides a Podio API client and authentication.
package podio

type Client struct{}

type ClientOptions struct {
	ApiKey    string
	ApiSecret string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) AuthenticateWithAuthCode(authCode, redirectURI string) ([]byte, error) {
	return nil, nil
}

func (c *Client) AuthenticateWithCredentials(username, password string) ([]byte, error) {
	return nil, nil
}
