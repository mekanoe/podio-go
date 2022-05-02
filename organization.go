package podio

import (
	"fmt"
)

type Organization struct {
	ID           int      `json:"org_id,omitempty"`
	Name         string   `json:"name,omitempty"`
	URL          string   `json:"url,omitempty"`
	URLLabel     string   `json:"url_label,omitempty"`
	Status       string   `json:"status,omitempty"`
	Premium      bool     `json:"premium,omitempty"`
	Logo         int      `json:"logo,omitempty"`
	UserLimit    int      `json:"user_limit,omitempty"`
	Type         string   `json:"type,omitempty"`
	Segment      string   `json:"segment,omitempty"`
	Tier         string   `json:"tier,omitempty"`
	SalesAgentID int      `json:"sales_agent_id,omitempty"`
	Domains      []string `json:"domains,omitempty"`
	Role         string   `json:"role,omitempty"`
	CreatedOn    string   `json:"created_on,omitempty"`
	CreatedBy    User     `json:"created_by,omitempty"`
	Image        Image    `json:"image,omitempty"`
	Rights       []string `json:"rights,omitempty"`
}

func (c *Client) GetOrganization(orgId string) (*Organization, error) {
	org := &Organization{}
	err := c.get(fmt.Sprintf("/org/%s", orgId), org)
	return org, err
}

func (c *Client) GetOrganizationBySlug(orgSlug string) (*Organization, error) {
	org := &Organization{}
	err := c.get(fmt.Sprintf("/org/url?url=https://podio.com/%s", orgSlug), org)
	return org, err
}
