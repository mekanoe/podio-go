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
