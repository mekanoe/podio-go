package podio

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Organization struct {
	ID           int       `json:"org_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	URL          string    `json:"url,omitempty"`
	URLLabel     string    `json:"url_label,omitempty"`
	Status       string    `json:"status,omitempty"`
	Premium      bool      `json:"premium,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	UserLimit    int       `json:"user_limit,omitempty"`
	Type         string    `json:"type,omitempty"`
	Segment      string    `json:"segment,omitempty"`
	Tier         string    `json:"tier,omitempty"`
	SalesAgentID int       `json:"sales_agent_id,omitempty"`
	Domains      []string  `json:"domains,omitempty"`
	Role         string    `json:"role,omitempty"`
	CreatedOn    string    `json:"created_on,omitempty"`
	CreatedBy    CreatedBy `json:"created_by,omitempty"`
	Image        Image     `json:"image,omitempty"`
	Rights       []string  `json:"rights,omitempty"`
}

func (c *Client) GetOrganization(orgId string) (*Organization, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("/org/%s", orgId))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get organization: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("podio-go: failed to get organization: status code %d", resp.StatusCode)
	}

	org := &Organization{}
	json.NewDecoder(resp.Body).Decode(org)
	return org, nil
}

func (c *Client) GetOrganizationBySlug(orgSlug string) (*Organization, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("/org/url?url=https://podio.com/%s", orgSlug))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to get organization: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("podio-go: failed to get organization: status code %d", resp.StatusCode)
	}

	org := &Organization{}
	json.NewDecoder(resp.Body).Decode(org)
	return org, nil
}
