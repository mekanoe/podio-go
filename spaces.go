package podio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Space struct {
	ID              int      `json:"space_id"`
	Name            string   `json:"name,omitempty"`
	URL             string   `json:"url,omitempty"`
	Privacy         string   `json:"privacy,omitempty"`
	AutoJoin        bool     `json:"auto_join,omitempty"`
	URLLabel        string   `json:"url_label,omitempty"`
	PostOnNewApp    bool     `json:"post_on_new_app,omitempty"`
	PostOnNewMember bool     `json:"post_on_new_member,omitempty"`
	Rights          []string `json:"rights,omitempty"`
	Role            string   `json:"role,omitempty"`
	Subscribed      bool     `json:"subscribed,omitempty"`
	CreatedOn       string   `json:"created_on,omitempty"`
	CreatedBy       struct {
		UserID string `json:"user_id,omitempty"`
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	} `json:"created_by,omitempty"`
}

// GetSpace returns the space/workspace with the given ID.
func (c *Client) GetSpace(spaceID string) (*Space, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("/space/%s", spaceID))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get space: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("podio-go: failed to get space: status code %d", resp.StatusCode)
	}

	space := &Space{}
	json.NewDecoder(resp.Body).Decode(space)
	return space, nil
}

func (c *Client) GetSpaceByURL(url string) (*Space, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("/space/url?url=%s", url))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get space: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("podio-go: failed to get space: status code %d", resp.StatusCode)
	}

	space := &Space{}
	json.NewDecoder(resp.Body).Decode(space)
	return space, nil
}

type CreateSpaceParams struct {
	Name            string `json:"name"`
	OrgID           int    `json:"org_id"`
	Privacy         string `json:"privacy"`
	PostOnNewApp    bool   `json:"post_on_new_app"`
	PostOnNewMember bool   `json:"post_on_new_member"`
	AutoJoin        bool   `json:"auto_join"`
}

func (c *Client) CreateSpace(params CreateSpaceParams) (*Space, error) {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to create space request: %w", err)
	}

	resp, err := c.httpClient.Post("/space/", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to create space: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(output))
		return nil, fmt.Errorf("podio-go: failed to create space: status code %d", resp.StatusCode)
	}

	data := &struct {
		Url string `json:"url"`
	}{}
	json.NewDecoder(resp.Body).Decode(data)

	space, err := c.GetSpaceByURL(data.Url)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get space after creation: %w", err)
	}

	return space, nil
}

func (c *Client) DeleteSpace(spaceID string) error {
	url, err := url.Parse(fmt.Sprintf("/space/%s", spaceID))
	if err != nil {
		return fmt.Errorf("podio-go: delete space failed to create url: %w", err)
	}
	req := &http.Request{
		Method: "DELETE",
		URL:    url,
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("podio-go: delete space failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		output, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(output))
		return fmt.Errorf("podio-go: delete space failed: status code %d", resp.StatusCode)
	}

	return nil
}
