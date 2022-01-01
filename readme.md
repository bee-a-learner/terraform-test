https://azure.github.io/azure-sdk/golang_introduction.html

set variables
https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#environment-variables

az ad sp create-for-rbac --name "golang-spn" --role contributor --scopes /subscriptions/2a04288a-8136-4880-b526-c6070e59f004 --sdk-auth

{
  "clientId": "<clientId>",
  "clientSecret": "<clientSecret>",
  "subscriptionId": "<subscriptionId>",
  "tenantId": "<active_directory_tenant_id",
  "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
  "resourceManagerEndpointUrl": "https://management.azure.com/",
  "activeDirectoryGraphResourceId": "https://graph.windows.net/",
  "sqlManagementEndpointUrl": "https://management.core.windows.net:8443/",
  "galleryEndpointUrl": "https://gallery.azure.com/",
  "managementEndpointUrl": "https://management.core.windows.net/"
}

$env:AZURE_TENANT_ID="<active_directory_tenant_id"
$env:AZURE_CLIENT_ID="<service_principal_appid>"
$env:AZURE_CLIENT_SECRET="<service_principal_password>"
$env:AZURE_SUBSCRIPTION_ID= "2a04288a-8136-4880-b526-c6070e59f004"


run the the current file
```
1. run the init command
go mod init main -- where main is the package name in your go file

2. run tidy
go mod tidy

3. run go file local in the same folder
go run .

```
for test file
```
go mod init terratestmodules 
or
go mod init terratestmodules && go mod tidy
4. run test file
go test -v -timeout 60m
```