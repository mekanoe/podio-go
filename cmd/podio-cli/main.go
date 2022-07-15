package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kayteh/podio-go"
)

func main() {
	username := os.Getenv("PODIO_USERNAME")
	password := os.Getenv("PODIO_PASSWORD")
	clientID := os.Getenv("PODIO_CLIENT_ID")
	clientSecret := os.Getenv("PODIO_CLIENT_SECRET")

	if username == "" || password == "" || clientID == "" || clientSecret == "" {
		fmt.Println("PODIO_USERNAME, PODIO_PASSWORD, PODIO_CLIENT_ID and PODIO_CLIENT_SECRET must be set")
		os.Exit(1)
	}

	client := podio.NewClient(podio.ClientOptions{
		ApiKey:    clientID,
		ApiSecret: clientSecret,
		UserAgent: "podio-cli",
	})

	err := client.AuthenticateWithCredentials(username, password)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		os.Exit(1)
	}

	outputEncoder := json.NewEncoder(os.Stdout)
	outputEncoder.SetIndent("", " ")

	if os.Args[1] == "space" {
		spaceID := os.Args[2]
		var space *podio.Space
		var err error
		if strings.HasPrefix(spaceID, "http") {
			space, err = client.GetSpaceByURL(spaceID)
		} else {
			space, err = client.GetSpace(spaceID)
		}
		if err != nil {
			fmt.Println("Failed to get space:", err)
			os.Exit(1)
		}
		outputEncoder.Encode(space)
	}

	if os.Args[1] == "field" {
		appID := os.Args[2]
		fieldID := os.Args[3]
		var field *podio.Field
		var err error
		field, err = client.GetField(appID, fieldID)
		if err != nil {
			fmt.Println("Failed to get field:", err)
			os.Exit(1)
		}
		outputEncoder.Encode(field)
	}

	if os.Args[1] == "create-field" {
		appID := os.Args[2]
		fieldType := os.Args[3]
		label := os.Args[4]

		fieldConfig := podio.FieldConfig{
			Label: label,
		}

		var params = podio.CreateFieldParams{
			Type:   fieldType,
			Config: fieldConfig,
		}

		outputEncoder.Encode(params)

		field, err := client.CreateField(appID, params)

		if err != nil {
			fmt.Println("Failed to create field: ", err)
		}

		outputEncoder.Encode(field)
	}

	if os.Args[1] == "update-field" {
		appID := os.Args[2]
		fieldID := os.Args[3]
		label := os.Args[4]

		fieldConfig := podio.FieldConfig{
			Label: label,
		}

		field, err := client.UpdateField(appID, fieldID, fieldConfig)

		if err != nil {
			fmt.Println("Failed to create field: ", err)
		}

		outputEncoder.Encode(field)
	}

	if os.Args[1] == "delete-field" {
		appID := os.Args[2]
		fieldID := os.Args[3]
		deleteValues := false

		if len(os.Args) == 5 && os.Args[4] == "true" {
			deleteValues = true
		}

		err := client.DeleteField(appID, fieldID, deleteValues)
		if err != nil {
			fmt.Println("Error deleting field: ", err)
		}
	}

	if os.Args[1] == "app" {
		appID := os.Args[2]
		var application *podio.Application
		var err error
		application, err = client.GetApplication(appID)
		if err != nil {
			fmt.Println("Failed to get app: ", err)
			os.Exit(1)
		}
		outputEncoder.Encode(application)
	}

	if os.Args[1] == "applications" {
		spaceID := os.Args[2]
		var applications *[]podio.Application
		var err error
		applications, err = client.GetApplications(spaceID)
		if err != nil {
			fmt.Println("Failed to get applications: ", err)
			os.Exit(1)
		}
		outputEncoder.Encode(applications)
	}

	if os.Args[1] == "workspaces" {
		orgID := os.Args[2]
		var spaces *[]podio.Space
		var err error
		spaces, err = client.GetWorkSpaces(orgID)
		if err != nil {
			fmt.Println("Failed to get space: ", err)
			os.Exit(1)
		}
		outputEncoder.Encode(spaces)
	}

	if os.Args[1] == "organizations" {
		var organizations *[]podio.Organization
		var err error
		organizations, err = client.GetOrganizations()
		if err != nil {
			fmt.Println("Failed to get organization:", err)
			os.Exit(1)
		}
		outputEncoder.Encode(organizations)
	}

	if os.Args[1] == "create-space" {
		orgID, _ := strconv.Atoi(os.Args[2])
		spaceName := os.Args[3]
		space, err := client.CreateSpace(podio.CreateSpaceParams{
			OrgID:           orgID,
			Name:            spaceName,
			Privacy:         "closed",
			AutoJoin:        false,
			PostOnNewApp:    false,
			PostOnNewMember: false,
		})

		if err != nil {
			fmt.Println("Failed to create space:", err)
			os.Exit(1)
		}

		outputEncoder.Encode(space)
	}

	if os.Args[1] == "delete-space" {
		err := client.DeleteSpace(os.Args[2])
		if err != nil {
			fmt.Println("Failed to delete space:", err)
			os.Exit(1)
		}

		fmt.Println("Space deleted")
	}

	if os.Args[1] == "rename-space" {
		spaceID := os.Args[2]
		newName := os.Args[3]

		space, err := client.UpdateSpace(spaceID, podio.CreateSpaceParams{
			Name: newName,
		})
		if err != nil {
			fmt.Println("Failed to update space:", err)
			os.Exit(1)
		}

		outputEncoder.Encode(space)
	}

}
