package azure_docker_login

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAzureDockerLogin_Action(t *testing.T) {
	a := AzureDockerLogin{}
	if got := a.Action(); got != "azure/docker-login@v2" {
		t.Errorf("Action() = %q, want %q", got, "azure/docker-login@v2")
	}
}

func TestAzureDockerLogin_Inputs(t *testing.T) {
	a := AzureDockerLogin{
		LoginServer: "myregistry.azurecr.io",
		Username:    "${{ secrets.ACR_USERNAME }}",
		Password:    "${{ secrets.ACR_PASSWORD }}",
	}

	inputs := a.Inputs()

	if inputs["login-server"] != "myregistry.azurecr.io" {
		t.Errorf("inputs[login-server] = %v, want myregistry.azurecr.io", inputs["login-server"])
	}
	if inputs["username"] != "${{ secrets.ACR_USERNAME }}" {
		t.Errorf("inputs[username] = %v, want secret reference", inputs["username"])
	}
	if inputs["password"] != "${{ secrets.ACR_PASSWORD }}" {
		t.Errorf("inputs[password] = %v, want secret reference", inputs["password"])
	}
}

func TestAzureDockerLogin_Inputs_Empty(t *testing.T) {
	a := AzureDockerLogin{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AzureDockerLogin.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAzureDockerLogin_ImplementsStepAction(t *testing.T) {
	a := AzureDockerLogin{}
	var _ workflow.StepAction = a
}
