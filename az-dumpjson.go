package main

/*
  This uses an unofficial Azure SDK auth package. Get it using:

		go get github.com/yaegashi/msgraph.go/msauth

		export AZURE_TENANT_ID="xxx"
		export AZURE_CLIENT_ID="xxx"
		export AZURE_CLIENT_SECRET="xxx"

	It dumps all JSON info from the API to stdout. Normally this would
	be ingested by a log shipper.

	Note: There is no official Azure Graph API SDK for Golang. Second
    	best is a officially supported Python SDK:
		https://developer.microsoft.com/en-us/graph/get-started/python

*/

import (
	"context"
	"fmt"
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
	// list all users available in the Azure AD
	// https://docs.microsoft.com/en-us/graph/api/resources/users?view=graph-rest-1.0
	usersAPI string = "https://graph.microsoft.com/beta/users"

	// list all groups in the Azure AD
	// https://docs.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-1.0
	groupsAPI string = "https://graph.microsoft.com/beta/groups"

	// after resolving all group ID's in the Azure AD, we can dump its members
	// on a per ID basis
	// https://docs.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-1.0&tabs=http
	groupMembersAPI string = "https://graph.microsoft.com/v1.0/groups/<id>/members?$count=true"

	// list all devices in the Azure AD, these can be managed by Intune
	// https://docs.microsoft.com/en-us/graph/api/resources/device?view=graph-rest-1.0
	devicesAPI string = "https://graph.microsoft.com/beta/devices"

	// Conditional Access policies
	// List all the policies applied in the Azure IAM
	// See: https://docs.microsoft.com/en-us/graph/api/conditionalaccessroot-list-policies?view=graph-rest-1.0&tabs=http
	// Alerting/Audit logging on Conditional Access:
	//    => https://www.chorus.co/resources/news/setting-up-conditional-access-alerts
	//    => https://docs.microsoft.com/en-us/azure/sentinel/connect-azure-active-directory
	//    => https://techcommunity.microsoft.com/t5/azure-active-directory-identity/conditional-access-reporting/m-p/1379594
	//    => https://docs.microsoft.com/en-us/azure/security-center/security-center-alerts-overview
	conditionalAccessAPI string = "https://graph.microsoft.com/v1.0/identity/conditionalAccess/policies"

	// Azure AD audit logging
	// This allows logging on Azure AD activity, including user signin events and activity
	// See https://docs.microsoft.com/en-us/graph/api/resources/azure-ad-auditlog-overview?view=graph-rest-1.0
	// This also allows for logging of Microsoft 365 products: OneDrive activity, Outlook, etc
	auditLogsAPI  string = "https://graph.microsoft.com/v1.0/auditLogs/directoryaudits"
	signinLogsAPI string = "https://graph.microsoft.com/v1.0/auditLogs/signIns"

	// Access the Microsoft Azure Identity Protection features
	// See also: https://docs.microsoft.com/en-us/graph/api/resources/identityprotectionroot?view=graph-rest-1.0
	// This is licensed. You need a Premium 1 or Premium 2 license
	riskyUsersAPI     string = "https://graph.microsoft.com/v1.0/identityProtection/riskyUsers"
	riskDetectionsAPI string = "https://graph.microsoft.com/v1.0/identityProtection/riskDetections"

	// Access the Microsoft Security graph API
	// Basically a unified API for multiple Azure Security products
	// See also: https://docs.microsoft.com/en-us/graph/api/resources/security-api-overview?view=graph-rest-1.0
	// Alerts can be filtered on a per product basis: https://docs.microsoft.com/en-us/graph/api/alert-list?view=graph-rest-1.0&tabs=http
	securityAlertsAPI string = "https://graph.microsoft.com/v1.0/security/alerts"

	// Microsoft Intune API reference guide:
	// https://docs.microsoft.com/en-us/graph/api/resources/intune-graph-overview?view=graph-rest-beta
	// => Now called Microsoft Endpoint Manager
	// Get devices registered to a user
	userDeviceAPI string = "https://graph.microsoft.com/beta/users/{user}/ownedDevices"

	// Get all Mobile Apps
	mobileAppsAPI   string = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
	applicationsAPI string = "https://graph.microsoft.com/v1.0/applications"
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

	{
		// This is an example making use of the raw API
		// Really helpful is https://developer.microsoft.com/en-us/graph/graph-explorer
		response, err := httpClient.Get(usersAPI)
		fmt.Println("\n\n=====>", usersAPI)
		if err != nil {
			log.Fatal(err)
		} else {
			io.Copy(os.Stdout, response.Body)
		}
		defer response.Body.Close()
	}
	{
		response, err := httpClient.Get(groupsAPI)
		fmt.Println("\n\n=====>", groupsAPI)
		if err != nil {
			log.Fatal(err)
		} else {
			io.Copy(os.Stdout, response.Body)
		}
		defer response.Body.Close()
	}
}
