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

type Application struct {
	// "app_id": The id of the app,
	AppID int `json:"app_id,omitempty"`
	// "original": The original app that this app was a copy of,
	OriginalID int `json:"original,omitempty"`
	// "original_revision": The revision of the original app at the time of copy,
	OriginalRevision int `json:"original_revision,omitempty"`
	// "status": The status of the app, either "active", "inactive" or "deleted".
	Status string `json:"status,omitempty"`
	// "space_id": The id of the space on which the app is placed,
	SpaceID int `json:"space_id,omitempty"`
	// "token": The app token to use when logging in as an app,
	Token string `json:"token,omitempty"`
	// "mailbox": The mailbox to use when mailing to the app,
	Mailbox string `json:"mailbox,omitempty"`
	// "owner": The owner of the app, which has special access to the app
	Owner User `json:"owner,omitempty"`
	// "config": The current configuration of the app,
	Config AppConfig `json:"config,omitempty"`
	// "integration": The current status of the integration, if any,
	Integration AppIntegration `json:"integration,omitempty"`
	// "rights": The array of rights the active user has on the app,
	Rights []string `json:"rights,omitempty"`
	Fields []Field  `json:"fields,omitempty"`
}

type CreateApplicationParams struct {
	Config AppConfig `json:"config,omitempty"`
	Fields []Field   `json:"fields,omitempty"`

	// Only applicable for creation
	SpaceID int `json:"space_id,omitempty"`

	// Only applicable for updates
	FieldsToDelete []FieldDelete `json:"fields_to_delete,omitempty"`
}

type AppConfig struct {
	// "type": The type of the app, either "standard", "meeting" or "contact",
	Type string `json:"type,omitempty"`
	// "name": The name of the app,
	Name string `json:"name,omitempty"`
	// "item_name": The name of each item in an app,
	ItemName string `json:"item_name,omitempty"`
	// "description": The description of the app,
	Description string `json:"description,omitempty"`
	// "usage": Description of how the app should be used,
	Usage string `json:"usage,omitempty"`
	// "external_id": The external id of the app. This can be used to store an id from an external system on the app,
	ExternalID string `json:"external_id,omitempty"`
	// "icon": The name of the icon used to represent the app,
	Icon string `json:"icon,omitempty"`
	// "allow_edit": True if other members are allowed to edit items from the app, false otherwise,
	AllowEdit bool `json:"allow_edit,omitempty"`
	// "default_view": The default view of the app items on the app main page (see area for more information),
	DefaultView string `json:"default_view,omitempty"`
	// "allow_attachments": True if attachment of files to an item is allowed, false otherwise,
	AllowAttachments bool `json:"allow_attachments,omitempty"`
	// "allow_comments": True if members can make comments on an item, false otherwise,
	AllowComments bool `json:"allow_comments,omitempty"`
	// "silent_creates": True if item creates should not be posted to the stream, false otherwise,
	SilentCreates bool `json:"silent_creates,omitempty"`
	// "silent_edits": True if item edits should not be posted to the stream, false otherwise,
	SilentEdits bool `json:"silent_edits,omitempty"`
	// "fivestar": True if fivestar rating is enabled on an item, false otherwise,
	FiveStar bool `json:"fivestar,omitempty"`
	// "fivestar_label": If fivestar rating is enabled, this is the label that will be presented to the users,
	FiveStarLabel string `json:"fivestar_label,omitempty"`
	// "approved": True if an item can be approved, false otherwise,
	Approved bool `json:"approved,omitempty"`
	// "thumbs": True if an item can have a thumbs up or thumbs down, false otherwise,
	Thumbs bool `json:"thumbs,omitempty"`
	// "thumbs_label": If thumbs ratings are enabled, this is the label that will be presented to the users,
	ThumbsLabel string `json:"thumbs_label,omitempty"`
	// "rsvp": True if RSVP is enabled, false otherwise,
	Rsvp bool `json:"rsvp,omitempty"`
	// "rsvp_label": If RSVP is enabled, this is the label that will be presented to the users,
	RsvpLabel string `json:"rsvp_label,omitempty"`
	// "yesno": True if yes/no rating is enabled, false otherwise,
	YesNo bool `json:"yesno,omitempty"`
	// "yesno_label": If yes/no is enabled, this is the label that will be presented to the users,
	YesNoLabel string `json:"yesno_label,omitempty"`
	// "tasks": The automatic tasks that are created when an item in the app is created,
	Tasks []AppTask `json:"tasks,omitempty"`
}

type AppTask struct {
	// "text": The text to use for the task,
	Text string `json:"text,omitempty"`
	// "responsible": The list of responsible for the task,
	Responsible []User `json:"responsible,omitempty"`
}

type AppIntegration struct {
	// "status": The status of the integration, either "inactive", "active", "disabled", or "error".
	Status string `json:"status,omitempty"`
	// "type": The type of the integration, see the integration area for details,
	Type string `json:"type,omitempty"`
	// "updating": True if the integration is currently updating, false otherwise,
	Updating bool `json:"updating,omitempty"`
	// "last_updated_on": The date and time of the last update,
	LastUpdatedOn string `json:"last_updated_on,omitempty"`
	// "next_refresh_on": The date and time when the integration will be refreshed
	NextRefreshOn string `json:"next_refresh_on,omitempty"`
}

func (c *Client) GetApplication(appID string) (*Application, error) {
	app := &Application{}
	err := c.get(fmt.Sprintf("/app/%s", appID), app)
	return app, err
}

func (c *Client) CreateApplication(spaceID string, params CreateApplicationParams) (*Application, error) {
	var err error
	params.SpaceID, err = strconv.Atoi(spaceID)
	if err != nil {
		return nil, fmt.Errorf("podio-go: invalid space id, must parse to int: %s", spaceID)
	}

	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("podio-go: could not encode params: %w", err)
	}

	resp, err := c.httpClient.Post("/app", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to create application: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("podio-go: failed to create application: %s\nPayload: %s", resp.Status, string(output))
	}

	//This line defines a struct type and simultaneously creates an empty struct
	data := &struct {
		AppID int `json:"app_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to decode application response: %w", err)
	}

	return c.GetApplication(strconv.Itoa(data.AppID))
}

func (c *Client) UpdateApplication(appID string, params CreateApplicationParams) (*Application, error) {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(params)

	url, err := url.Parse(fmt.Sprintf("/app/%s", appID))
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to generate url: %w", err)
	}

	req := &http.Request{
		Method: "PUT",
		URL:    url,
		Body:   ioutil.NopCloser(body),
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to update application: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		output, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("podio-go: failed to update application: %s\nPayload: %s", resp.Status, string(output))
	}

	return c.GetApplication(appID)
}

func (c *Client) DeleteApplication(appID string) error {
	return c.delete(fmt.Sprintf("/app/%s", appID))
}

func (c *Client) GetApplications(spaceID string) (*[]Application, error) {
	apps := &[]Application{}
	err := c.get(fmt.Sprintf("/app/space/%s/?include_inactive=false", spaceID), apps)
	return apps, err
}
