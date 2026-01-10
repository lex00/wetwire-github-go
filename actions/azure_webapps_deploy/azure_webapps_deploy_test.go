package azure_webapps_deploy

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAzureWebappsDeploy_Action(t *testing.T) {
	a := AzureWebappsDeploy{}
	if got := a.Action(); got != "azure/webapps-deploy@v3" {
		t.Errorf("Action() = %q, want %q", got, "azure/webapps-deploy@v3")
	}
}

func TestAzureWebappsDeploy_Inputs(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName: "my-web-app",
		Package: "./dist",
	}

	inputs := a.Inputs()

	if inputs["app-name"] != "my-web-app" {
		t.Errorf("inputs[app-name] = %v, want my-web-app", inputs["app-name"])
	}
	if inputs["package"] != "./dist" {
		t.Errorf("inputs[package] = %v, want ./dist", inputs["package"])
	}
}

func TestAzureWebappsDeploy_Inputs_Empty(t *testing.T) {
	a := AzureWebappsDeploy{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AzureWebappsDeploy.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAzureWebappsDeploy_Inputs_Container(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName: "my-container-app",
		Images:  "myregistry.azurecr.io/myapp:latest",
	}

	inputs := a.Inputs()

	if inputs["images"] != "myregistry.azurecr.io/myapp:latest" {
		t.Errorf("inputs[images] = %v, want image reference", inputs["images"])
	}
}

func TestAzureWebappsDeploy_Inputs_Slot(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:  "my-web-app",
		SlotName: "staging",
		Package:  "./dist",
	}

	inputs := a.Inputs()

	if inputs["slot-name"] != "staging" {
		t.Errorf("inputs[slot-name] = %v, want staging", inputs["slot-name"])
	}
}

func TestAzureWebappsDeploy_Inputs_Options(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName: "my-web-app",
		Package: "./dist",
		Clean:   true,
		Restart: true,
		Type:    "ZIP",
	}

	inputs := a.Inputs()

	if inputs["clean"] != true {
		t.Errorf("inputs[clean] = %v, want true", inputs["clean"])
	}
	if inputs["restart"] != true {
		t.Errorf("inputs[restart] = %v, want true", inputs["restart"])
	}
	if inputs["type"] != "ZIP" {
		t.Errorf("inputs[type] = %v, want ZIP", inputs["type"])
	}
}

func TestAzureWebappsDeploy_ImplementsStepAction(t *testing.T) {
	a := AzureWebappsDeploy{}
	var _ workflow.StepAction = a
}
