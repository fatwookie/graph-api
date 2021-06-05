package main

/*
  This uses the Azure SDK auth package. Get it using:
		go get -u github.com/Azure/go-autorest/autorest/azure/auth

		export AZURE_TENANT_ID="2afb9002-f677-491c-b76a-b7c2f8f026ee"
		export AZURE_CLIENT_ID="b940b2a4-38cd-40c7-a360-52def813e105"
		export AZURE_CLIENT_SECRET="_wd_kUm1ZMkQRc6~_Yet_zGi8Yq_.lH1_Y"

		This uses the (deprecated) Azure AD Graph API, thus it needs
		API permissions on the Azure AD graph API
*/

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func main() {

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println(err)
	}

	client := graphrbac.NewUsersClient(os.Getenv("AZURE_TENANT_ID"))
	client.Authorizer = authorizer

	if _, err := client.List(context.Background(), "", ""); err != nil {
		fmt.Println("list users", err)
	}
}
