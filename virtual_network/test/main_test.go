package main

// Import key modules.
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var (
	globalBackendConf = make(map[string]interface{})
	globalEnvVars     = make(map[string]string)
	uniquePostfix     = strings.ToLower(random.UniqueId())
	prefix            = "vnet"
	separator         = "-"
)

// Define key global variables.
var (
	subscriptionId      = "2a04288a-8136-4880-b526-c6070e59f004"
	resource_group_name = "example-resources"
)

const (
	apiVersion              = "2019-06-01"
	resourceProvisionStatus = "Succeeded"
)

func setTerraformVariables() (map[string]string, error) {

	// Getting enVars from environment variables
	ARM_CLIENT_ID := os.Getenv("AZURE_CLIENT_ID")
	ARM_CLIENT_SECRET := os.Getenv("AZURE_CLIENT_SECRET")
	ARM_TENANT_ID := os.Getenv("AZURE_TENANT_ID")
	ARM_SUBSCRIPTION_ID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Creating globalEnVars for terraform call through Terratest
	if ARM_CLIENT_ID != "" {
		globalEnvVars["ARM_CLIENT_ID"] = ARM_CLIENT_ID
		globalEnvVars["ARM_CLIENT_SECRET"] = ARM_CLIENT_SECRET
		globalEnvVars["ARM_SUBSCRIPTION_ID"] = ARM_SUBSCRIPTION_ID
		globalEnvVars["ARM_TENANT_ID"] = ARM_TENANT_ID
	}

	return globalEnvVars, nil
}

func TestTerraform_azure_virtualNetwork(t *testing.T) {
	t.Parallel()

	setTerraformVariables()

	expectedLocation := "uksouth"
	expectedAddressSpace := "10.0.0.0/8"

	//uniquePostfix := random.UniqueId() // "mce" - switch for terratest or manual terraform deployment

	// Use Terratest to deploy the infrastructure
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{

		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../provision",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{

			"resource_group_name": resource_group_name,
			"location":            expectedLocation,
			"prefix":              prefix,
			"postfix":             uniquePostfix,
			"address_space":       expectedAddressSpace,
			//tags = var.tags
		},

		// globalvariables for user account
		EnvVars: globalEnvVars,

		// Backend values to set when initialziing Terraform
		BackendConfig: globalBackendConf,

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,

		// Reconfigure is required if module deployment and go test pipelines are running in one stage
		Reconfigure: true,
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	//out, err := terraform.OutputJsonE(t, terraformOptions, "resource_name")
	expectedResourceName := terraform.Output(t, terraformOptions, "resource_name")
	expectedResourceId := terraform.Output(t, terraformOptions, "resource_id")
	expectedVnetAddressSpace := terraform.Output(t, terraformOptions, "vnet_address_space")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	fmt.Printf("ResourceName :: %s\n", expectedResourceName)
	fmt.Printf("resourceId :: %s\n", expectedResourceId)
	fmt.Printf("VnetAddressSpace :: %s\n", expectedVnetAddressSpace)
	fmt.Printf("resourceGroupName :: %s\n", resourceGroupName)

	resp, err := getResourceFromRESTAPI(expectedResourceId)

	if err != nil {
		log.Fatalf("failed to obtain a terraform var output as json: %v", err)
	}

	actualResponse := Deserialize(resp)

	//resourceProvisioningState := *resp.provision
	t.Run("vnet_resource_name_matched", func(t *testing.T) {
		// Check the Storage Account exists
		assert.Equal(t, expectedResourceName, actualResponse.Name, "vnet name matched ")
	})

	//resourceProvisioningState := *resp.provision
	t.Run("vnet_resource_provisioning_state_is_succeeded", func(t *testing.T) {
		// Check the Storage Account exists
		assert.Equal(t, resourceProvisionStatus, actualResponse.Properties.ProvisioningState)
	})

	t.Run("vnet_resource_address_space_is_matching", func(t *testing.T) {
		assert.Equal(t, expectedAddressSpace, actualResponse.Properties.AddressSpace.AddressPrefixes[0])
	})

	t.Run("vnet_resource_location_is_matching", func(t *testing.T) {
		assert.Equal(t, expectedLocation, actualResponse.Location)
	})

}

func getResourceFromRESTAPI(out string) (armresources.ResourcesGetByIDResponse, error) {

	//expected variable
	//expectedVnetName := strings.ToLower(fmt.Sprintf("%s%s%s", prefix, separator, uniquePostfix))

	log.Printf("json output: %s\n", out)

	ctx := context.Background()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	resourceId := fmt.Sprintf("%v", out) //result["resource_name"]
	//resourceId := fmt.Sprintf(resourceIdFormat, subscriptionId, resource_group_name, expectedResource_name)
	// Azure SDK Azure Resource Management clients accept the credential as a parameter
	client := armresources.NewResourcesClient(subscriptionId, cred, nil)

	resp, err := client.GetByID(ctx, resourceId, apiVersion, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}

	return resp, err
}

func Deserialize(resp armresources.ResourcesGetByIDResponse) ReponseBase {

	// // Unmarshal JSON to Result struct.
	// var result Result
	// json.Unmarshal(bytes, &result)
	fmt.Println("================================== JSON OUTPUT =================================================")
	b, _ := json.Marshal(resp)
	// Convert bytes to string.
	sOutput := string(b)
	fmt.Println(sOutput)

	fmt.Println("================================== OBJECT OUTPUT =================================================")

	// Get bytes.
	bytes := []byte(sOutput)

	// Unmarshal JSON to Result struct.
	var result ReponseBase
	json.Unmarshal(bytes, &result)

	//fmt.Println("" result.provisioningState)
	fmt.Printf("Result ProvisioningState:: %s\n", result.Properties.ProvisioningState)
	fmt.Printf("Result AddressSpace:: %s\n", result.Properties.AddressSpace.AddressPrefixes)
	fmt.Printf("Result ResourceId:: %s\n", result.ResourceId)
	fmt.Printf("Result ResourceName:: %s\n", result.Name)
	fmt.Printf("Result Location:: %s\n", result.Location)
	fmt.Printf("Result Location:: %s\n", reflect.TypeOf(result.Location))

	return result
}
