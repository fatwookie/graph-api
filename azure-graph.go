package main

/*
  This uses an unofficial Azure SDK auth package. Get it using:

		go get github.com/yaegashi/msgraph.go/msauth

		export AZURE_TENANT_ID="xxx"
		export AZURE_CLIENT_ID="xxx"
		export AZURE_CLIENT_SECRET="xxx"

		Note: There is no official Azure Graph API SDK for Golang. Second
		best is a officially supported Python SDK:
		https://developer.microsoft.com/en-us/graph/get-started/python

*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	msgraph "github.com/yaegashi/msgraph.go/beta"
	msauth "github.com/yaegashi/msgraph.go/msauth"
	"golang.org/x/oauth2"
)

var (
	tenantID     = os.Getenv("AZURE_TENANT_ID")
	clientID     = os.Getenv("AZURE_CLIENT_ID")
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
)

const (
	usersAPI        = "https://graph.microsoft.com/v1.0/users"
	groupsAPI       = "https://graph.microsoft.com/v1.0/groups"
	groupMembersAPI = "https://graph.microsoft.com/v1.0/groups/<id>/members?$count=true"
)

func main() {

	fmt.Println(
		"Azure PoC code for friendly intel\n",
		"This uses the community msgraph API package\n",
	)

	var scopes = []string{msauth.DefaultMSGraphScope}

	ctx := context.Background()
	m := msauth.NewManager()
	tokensource, err := m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, scopes)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2.NewClient(ctx, tokensource)

	// range over all users in the tenant
	// Note: this uses the deprecated Azure AD Graph API
	fmt.Println("====== msgraph package example ==========")

	graphClient := msgraph.NewClient(httpClient)
	r := graphClient.Users().Request()
	users, _ := r.Get(ctx)
	for _, user := range users {
		fmt.Printf("[*] Found user %s\n", *user.UserPrincipalName)
	}

	// This is an example making use of the raw API
	// Really helpful is https://developer.microsoft.com/en-us/graph/graph-explorer
	fmt.Println("====== raw API examples ==========")
	response, err := httpClient.Get(usersAPI)
	if err != nil {
		log.Fatal(err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

}
