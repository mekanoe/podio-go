// Package podio provides a Podio API client and authentication.
package podio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is a Podio API client.
type Client struct {
	httpClient *http.Client
	options    ClientOptions
	apiURL     *url.URL
}

// ClientOptions are the options for a Podio API client.
type ClientOptions struct {
	ApiKey    string
	ApiSecret string
	ApiURL    string
	UserAgent string
}

// NewClient creates a new Podio API client with the given options.
func NewClient(options ClientOptions) *Client {
	if options.ApiURL == "" {
		options.ApiURL = "https://api.podio.com"
	}

	if options.ApiKey == "" || options.ApiSecret == "" {
		panic("podio-go: ApiKey and ApiSecret are required")
	}

	apiURL, err := url.Parse(options.ApiURL)
	if err != nil {
		panic(fmt.Errorf("podio-go: failed to parse API URL: %w", err))
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &authenticatedTransport{
				options: options,
				apiURL:  apiURL,
			},
		},
		options: options,
		apiURL:  apiURL,
	}
}

// AuthenticateWithCreditentials authenticates with the given credentials.
func (c *Client) AuthenticateWithCredentials(username, password string) error {
	oauth, err := c.doOAuthGrant(oAuth2Request{
		Username:     username,
		Password:     password,
		ClientID:     c.options.ApiKey,
		ClientSecret: c.options.ApiSecret,
		GrantType:    "password",
	})

	if err != nil {
		return err
	}

	c.httpClient.Transport = &authenticatedTransport{
		apiToken: oauth.AccessToken,
		options:  c.options,
		apiURL:   c.apiURL,
	}

	return nil
}

func (c *Client) get(path string, v interface{}) error {
	resp, err := c.httpClient.Get(path)
	if err != nil {
		return fmt.Errorf("podio-go: failed to GET %s: %w", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("podio-go: failed to GET %s: %s\nPayload: %s", path, resp.Status, string(output))
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("podio-go: failed to decode response: %w", err)
	}

	return nil
}

func (c *Client) delete(path string) error {
	url, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("podio-go: failed to parse URL: %w", err)
	}
	req := &http.Request{
		Method: "DELETE",
		URL:    url,
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("podio-go: failed to DELETE %s: %w", path, err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("podio-go: failed to DELETE %s: %s\nPayload: %s", path, resp.Status, string(output))
	}

	return nil
}

type oAuth2Request struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type oAuth2Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Ref          struct {
		Type string `json:"type"`
		ID   uint64 `json:"id"`
	} `json:"ref"`
}

func (c *Client) doOAuthGrant(credentials oAuth2Request) (*oAuth2Token, error) {
	body := url.Values{}
	body.Set("grant_type", credentials.GrantType)
	body.Set("username", credentials.Username)
	body.Set("password", credentials.Password)
	body.Set("client_id", credentials.ClientID)
	body.Set("client_secret", credentials.ClientSecret)

	buf := bytes.NewBufferString(body.Encode())

	resp, err := c.httpClient.Post("/oauth/token", "application/x-www-form-urlencoded", buf)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to request OAuth token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("podio-go: failed to request OAuth token: %s; body: %s", resp.Status, string(output))
	}

	token := &oAuth2Token{}
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to decode OAuth token: %w", err)
	}

	return token, nil
}

type authenticatedTransport struct {
	apiToken string
	apiURL   *url.URL
	options  ClientOptions
}

func (a *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if a.options.UserAgent != "" {
		req.Header.Add("user-agent", a.options.UserAgent)
	} else {
		req.Header.Add("user-agent", "podio-go")
	}

	if a.apiToken != "" {
		req.Header.Add("authorization", "OAuth2 "+a.apiToken)
	}

	if req.URL.Host == "" {
		req.URL.Scheme = a.apiURL.Scheme
		req.URL.Host = a.apiURL.Host
	}

	return http.DefaultTransport.RoundTrip(req)
}
