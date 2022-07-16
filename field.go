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

type Field struct {
	// "field_id": The id of the field,
	FieldID int `json:"field_id,omitempty"`
	// "type": The type of the field (see area for more information),
	Type string `json:"type,omitempty"`
	// "external_id": External id automatically generated that will never change,
	ExternalID string `json:"external_id,omitempty"`
	// "config": The configuration of the field,
	Config FieldConfig `json:"config,omitempty"`
}

type FieldConfig struct {
	// "label": The label of the field, which is what the users will see,
	Label string `json:"label,omitempty"`
	// "description": The description of the field, shown to the user when inserting and editing,
	Description string `json:"description,omitempty"`
	// "delta": An integer indicating the order of the field compared to other fields,
	Delta int `json:"delta,omitempty"`
	// "settings": The settings of the field which depends on the type of the field (see area for more information),
	Settings interface{} `json:"settings,omitempty"`
	// "mapping": The mapping of the field, one of "meeting_time", "meeting_participants", "meeting_agenda" and "meeting_location" for type="meeting" and "contact_name","contact_job_title","contact_organization","contact_email","contact_phone","contact_address","contact_website","contact_notes","contact_image" for type="contact"
	Mapping string `json:"mapping,omitempty"`
	// "required": True if the field is required when creating and editing items, false otherwise
	Required bool `json:"required,omitempty"`
}

type FieldDelete struct {
	// "field_id": The id of the field,
	FieldID int `json:"field_id,omitempty"`
	// "delete_values": True if the values for the field should be deleted, false otherwise
	DeleteValues bool `json:"delete_values,omitempty"`
}

type CreateFieldParams struct {
	Type   string      `json:"type,omitempty"`
	Config FieldConfig `json:"config,omitempty"`
}

func (c *Client) CreateField(appID string, params CreateFieldParams) (*Field, error) {
	var err error

	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("podio-go: could not encode params: %w", err)
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("/app/%s/field/", appID), "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to create application: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("podio-go: failed to create field: %s\nPayload: %s", resp.Status, string(output))
	}

	//This line defines a struct type and simultaneously creates an empty struct
	data := &struct {
		FieldID int `json:"field_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("podio-go: failed to decode field creation response: %w", err)
	}

	return c.GetField(appID, strconv.Itoa(data.FieldID))
}

func (c *Client) GetField(appID string, fieldID string) (*Field, error) {
	field := &Field{}
	err := c.get(fmt.Sprintf("/app/%s/field/%s", appID, fieldID), field)
	return field, err
}

func (c *Client) UpdateField(appID string, fieldID string, params FieldConfig) (*Field, error) {
	var err error

	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("podio-go: could not encode params: %w", err)
	}

	url, err := url.Parse(fmt.Sprintf("/app/%s/field/%s/", appID, fieldID))
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
		return nil, fmt.Errorf("podio-go: failed to create field: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		output, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("podio-go: failed to update field: %s\nPayload: %s", resp.Status, string(output))
	}

	return c.GetField(appID, fieldID)
}

func (c *Client) DeleteField(appID string, fieldID string, deleteValues bool) error {
	deleteValuesStr := "false"
	if deleteValues {
		deleteValuesStr = "true"
	}
	return c.delete(fmt.Sprintf("/app/%s/field/%s?deleteValues=%s", appID, fieldID, deleteValuesStr))
}
