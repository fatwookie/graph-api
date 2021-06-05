package main

/*
  This uses an unofficial Azure SDK auth package. Get it using:

		go get github.com/yaegashi/msgraph.go/msauth

		export AZURE_TENANT_ID="xxx"
		export AZURE_CLIENT_ID="xxx"
		export AZURE_CLIENT_SECRET="xxx"

	It dumps all JSON info from the API to stdout. Normally this would
	be ingested by a log shipper.
*/

import (
	"context"
	"io"
	"log"
	"os"

	msauth "github.com/yaegashi/msgraph.go/msauth"
	"golang.org/x/oauth2"
)

var (
	tenantID     = os.Getenv("AZURE_TENANT_ID")
	clientID     = os.Getenv("AZURE_CLIENT_ID")
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
)

const (
	usersAPI      = "https://graph.microsoft.com/beta/users"
	devicesAPI    = "https://graph.microsoft.com/v1.0/devices"
	mobileAppsAPI = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
)

func main() {

	var scopes = []string{msauth.DefaultMSGraphScope}

	ctx := context.Background()
	m := msauth.NewManager()
	tokensource, err := m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, scopes)
	if err != nil {
		log.Fatal(err)
	}

	// this Interface client type can be used with http.Get()
	httpClient := oauth2.NewClient(ctx, tokensource)

	// This is an example making use of the raw API
	// Really helpful is https://developer.microsoft.com/en-us/graph/graph-explorer
	response, err := httpClient.Get(usersAPI)
	if err != nil {
		log.Fatal(err)
	} else {
		io.Copy(os.Stdout, response.Body)
	}
	defer response.Body.Close()

}
