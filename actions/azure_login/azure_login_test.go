package azure_login

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAzureLogin_Action(t *testing.T) {
	a := AzureLogin{}
	if got := a.Action(); got != "azure/login@v2" {
		t.Errorf("Action() = %q, want %q", got, "azure/login@v2")
	}
}

func TestAzureLogin_Inputs_Creds(t *testing.T) {
	a := AzureLogin{
		Creds: "${{ secrets.AZURE_CREDENTIALS }}",
	}

	inputs := a.Inputs()

	if inputs["creds"] != "${{ secrets.AZURE_CREDENTIALS }}" {
		t.Errorf("inputs[creds] = %v, want secret reference", inputs["creds"])
	}
}

func TestAzureLogin_Inputs_Empty(t *testing.T) {
	a := AzureLogin{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AzureLogin.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAzureLogin_Inputs_OIDC(t *testing.T) {
	a := AzureLogin{
		ClientID:       "${{ secrets.AZURE_CLIENT_ID }}",
		TenantID:       "${{ secrets.AZURE_TENANT_ID }}",
		SubscriptionID: "${{ secrets.AZURE_SUBSCRIPTION_ID }}",
	}

	inputs := a.Inputs()

	if inputs["client-id"] != "${{ secrets.AZURE_CLIENT_ID }}" {
		t.Errorf("inputs[client-id] = %v, want secret reference", inputs["client-id"])
	}
	if inputs["tenant-id"] != "${{ secrets.AZURE_TENANT_ID }}" {
		t.Errorf("inputs[tenant-id] = %v, want secret reference", inputs["tenant-id"])
	}
	if inputs["subscription-id"] != "${{ secrets.AZURE_SUBSCRIPTION_ID }}" {
		t.Errorf("inputs[subscription-id] = %v, want secret reference", inputs["subscription-id"])
	}
}

func TestAzureLogin_Inputs_Options(t *testing.T) {
	a := AzureLogin{
		ClientID:          "my-client-id",
		TenantID:          "my-tenant-id",
		Environment:       "azureusgovernment",
		EnableAzPSSession: true,
	}

	inputs := a.Inputs()

	if inputs["environment"] != "azureusgovernment" {
		t.Errorf("inputs[environment] = %v, want azureusgovernment", inputs["environment"])
	}
	if inputs["enable-AzPSSession"] != true {
		t.Errorf("inputs[enable-AzPSSession] = %v, want true", inputs["enable-AzPSSession"])
	}
}

func TestAzureLogin_Inputs_AllowNoSubscriptions(t *testing.T) {
	a := AzureLogin{
		ClientID:             "my-client-id",
		TenantID:             "my-tenant-id",
		AllowNoSubscriptions: true,
	}

	inputs := a.Inputs()

	if inputs["allow-no-subscriptions"] != true {
		t.Errorf("inputs[allow-no-subscriptions] = %v, want true", inputs["allow-no-subscriptions"])
	}
}

func TestAzureLogin_Inputs_Audience(t *testing.T) {
	a := AzureLogin{
		ClientID: "my-client-id",
		TenantID: "my-tenant-id",
		Audience: "api://AzureADTokenExchange",
	}

	inputs := a.Inputs()

	if inputs["audience"] != "api://AzureADTokenExchange" {
		t.Errorf("inputs[audience] = %v, want api://AzureADTokenExchange", inputs["audience"])
	}
}

func TestAzureLogin_Inputs_AuthType(t *testing.T) {
	a := AzureLogin{
		ClientID: "my-client-id",
		TenantID: "my-tenant-id",
		AuthType: "SERVICE_PRINCIPAL",
	}

	inputs := a.Inputs()

	if inputs["auth-type"] != "SERVICE_PRINCIPAL" {
		t.Errorf("inputs[auth-type] = %v, want SERVICE_PRINCIPAL", inputs["auth-type"])
	}
}

func TestAzureLogin_Inputs_BooleanFalse(t *testing.T) {
	a := AzureLogin{
		ClientID:             "my-client-id",
		TenantID:             "my-tenant-id",
		EnableAzPSSession:    false,
		AllowNoSubscriptions: false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["enable-AzPSSession"]; ok {
		t.Error("enable-AzPSSession=false should not be in inputs")
	}
	if _, ok := inputs["allow-no-subscriptions"]; ok {
		t.Error("allow-no-subscriptions=false should not be in inputs")
	}
}

func TestAzureLogin_ImplementsStepAction(t *testing.T) {
	a := AzureLogin{}
	var _ workflow.StepAction = a
}
