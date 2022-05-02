package podio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Space struct {
	ID              int      `json:"space_id"`
	Name            string   `json:"name,omitempty"`
	URL             string   `json:"url,omitempty"`
	OrgID           int      `json:"org_id,omitempty"`
	Privacy         string   `json:"privacy,omitempty"`
	AutoJoin        bool     `json:"auto_join,omitempty"`
	URLLabel        string   `json:"url_label,omitempty"`
	PostOnNewApp    bool     `json:"post_on_new_app,omitempty"`
	PostOnNewMember bool     `json:"post_on_new_member,omitempty"`
	Rights          []string `json:"rights,omitempty"`
	Role            string   `json:"role,omitempty"`
	Subscribed      bool     `json:"subscribed,omitempty"`
	CreatedOn       string   `json:"created_on,omitempty"`
	CreatedBy       User     `json:"created_by,omitempty"`
}

func (c *Client) GetSpace(spaceID string) (*Space, error) {
	space := &Space{}
	err := c.get(fmt.Sprintf("/space/%s", spaceID), space)
	return space, err
}

func (c *Client) GetSpaceByURL(url string) (*Space, error) {
	space := &Space{}
	err := c.get(fmt.Sprintf("/space/url?url=%s", url), space)
	return space, err
}

type CreateSpaceParams struct {
	Name            string `json:"name,omitempty"`
	OrgID           int    `json:"org_id,omitempty"`
	Privacy         string `json:"privacy,omitempty"`
	PostOnNewApp    bool   `json:"post_on_new_app,omitempty"`
	PostOnNewMember bool   `json:"post_on_new_member,omitempty"`
	AutoJoin        bool   `json:"auto_join,omitempty"`
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

		return nil, fmt.Errorf("podio-go: failed to create space: status code %d; return payload: %s", resp.StatusCode, string(output))
	}

	data := &struct {
		ID  int    `json:"space_id"`
		Url string `json:"url"`
	}{}
	json.NewDecoder(resp.Body).Decode(data)

	space, err := c.GetSpace(strconv.Itoa(data.ID))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get space after creation: %w", err)
	}

	return space, nil
}

func (c *Client) DeleteSpace(spaceID string) error {
	return c.delete(fmt.Sprintf("/space/%s", spaceID))
}

func (c *Client) UpdateSpace(spaceID string, params CreateSpaceParams) (*Space, error) {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to update space request: %w", err)
	}

	url, err := url.Parse(fmt.Sprintf("/space/%s", spaceID))
	if err != nil {
		return nil, fmt.Errorf("podio-go: update space failed to create url: %w", err)
	}
	req := &http.Request{
		Method: "PUT",
		URL:    url,
		Body:   ioutil.NopCloser(body),
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("podio-go: update space failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		output, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(output))
		return nil, fmt.Errorf("podio-go: update space failed: status code %d", resp.StatusCode)
	}

	space, err := c.GetSpace(spaceID)
	return space, err
}
