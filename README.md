# Podio-go

Unofficial, semi-internally maintained Go client for Podio.

The only goal for this is to track https://github.com/kayteh/terraform-provider-oidop. This will never have perfect feature parity unless it's picked up by the Podio team for real.

# Usage

## Using this SDK in your Go project

Import:
```
import(
 "github.com/kayteh/podio-go"
)
```
Initialize the client:
```
var client *podio.Client
client = podio.NewClient(podio.ClientOptions{
		ApiKey:    "...",
		ApiSecret: "...",
		UserAgent: "...(anything can go in here)",
	})
```
Client is now ready to use, example:
```
 var appId := "1234"
 var app := client.GetApplication(appId) //get the details of a Podio App
 ...
```

## For developers: testing the SDK via the CLI (see `./cmd/podio-cli`)

Developers can add to the CLI the calls that they want to test out. Right now, the CLI only tests the CRUD operations over Podio 'Spaces'.

Change directory to the CLI program and build it:
```
cd ./cmd/podio-cli
go build
```

Set the following environment variables (code below works in bash):
```
export PODIO_USERNAME=abcd
export PODIO_PASSWORD=abcd1234
export PODIO_CLIENT_ID=something
export PODIO_CLIENT_SECRET=something
```

Use the CLI:

Get all organizations:
```
podio-cli organizations
```

Get all workspaces within an organization:
* `<orgID>` is a number
```
podio-cli workspaces <orgID>
```

Get all applications within a space:
* `<spaceID>` is a number
```
podio-cli applications <spaceID>
```

Read an application:
* `<appID>` is a number
```
podio-cli app <appID>
```

Create a field within an app
* `<appID>` is a number, the id of the app you wish to add the field to
* `<fieldType>` is the type of field you wish to create: text, date, location, phone, etc.
* `<label>` is the label you wish to give the field
```
podio-cli create-field <appID> <fieldType> <label>
```

Read a field
* `<appID>` is a number
* `<fieldID>` is a number
```
podio-cli field <appID> <fieldID>
```

Update a field
* `<appID>` is a number, the id of the app the field belongs to
* `<fieldID>` is a number, the id of the field you wish to update
* `<label>` is the new label you wish to give the field
```
podio-cli update-field <appID> <fieldID> <label>
```

Delete a field
* `<appID>` is a number, the id of the app the field belongs to
* `<fieldID>` is a number, the id of the field you wish to delete
* `<deleteValues>` can be set to either `true` or `false`, or it can be omitted (defaults to `false`). When `true`, the values in the field for existing items will be deleted.
```
podio-cli delete-field <appID> <fieldID> <deleteValues>
```

Read a Space:
* `<spaceID>` is a number
```
podio-cli space <spaceID>
```

Create a Space:
* `<orgID>` is a number
* `<spaceName>` is the name you want the space to have
```
podio-cli create-space <orgID> <spaceName>
```

Rename a Space:
* `<spaceID>` is a number
* `<newName>` is the new name you want the space to have
```
podio-cli rename-space <spaceID> <newName>
```

Delete a Space:
* `<spaceID>` is a number
```
podio-cli delete-space <spaceID>
```