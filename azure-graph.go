package main

/*
  This uses an unofficial Azure SDK auth package. Get it using:

		go get github.com/yaegashi/msgraph.go/msauth

		export AZURE_TENANT_ID="xxx"
		export AZURE_CLIENT_ID="xxx"
		export AZURE_CLIENT_SECRET="xxx"

*/

import (
	"context"
	"fmt"
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

func main() {

	fmt.Println("Azure PoC code for friendly intel")

	var scopes = []string{msauth.DefaultMSGraphScope}

	ctx := context.Background()
	m := msauth.NewManager()
	tokensource, err := m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, scopes)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2.NewClient(ctx, tokensource)
	graphClient := msgraph.NewClient(httpClient)

	// range over all users in the tenant
	//
	r := graphClient.Users().Request()
	users, _ := r.Get(ctx)
	for _, user := range users {
		fmt.Printf("[*] Found user %s\n", *user.UserPrincipalName)
	}

}
