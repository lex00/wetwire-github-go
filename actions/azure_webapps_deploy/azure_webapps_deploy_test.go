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

func TestAzureWebappsDeploy_Inputs_PublishProfile(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:        "my-web-app",
		PublishProfile: "${{ secrets.AZURE_WEBAPP_PUBLISH_PROFILE }}",
	}

	inputs := a.Inputs()

	if inputs["publish-profile"] != "${{ secrets.AZURE_WEBAPP_PUBLISH_PROFILE }}" {
		t.Errorf("inputs[publish-profile] = %v, want secret reference", inputs["publish-profile"])
	}
}

func TestAzureWebappsDeploy_Inputs_ConfigurationFile(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:           "my-container-app",
		ConfigurationFile: "docker-compose.yml",
	}

	inputs := a.Inputs()

	if inputs["configuration-file"] != "docker-compose.yml" {
		t.Errorf("inputs[configuration-file] = %v, want docker-compose.yml", inputs["configuration-file"])
	}
}

func TestAzureWebappsDeploy_Inputs_StartupCommand(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:        "my-web-app",
		StartupCommand: "dotnet myapp.dll",
	}

	inputs := a.Inputs()

	if inputs["startup-command"] != "dotnet myapp.dll" {
		t.Errorf("inputs[startup-command] = %v, want dotnet myapp.dll", inputs["startup-command"])
	}
}

func TestAzureWebappsDeploy_Inputs_ResourceGroupName(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:           "my-web-app",
		ResourceGroupName: "my-resource-group",
	}

	inputs := a.Inputs()

	if inputs["resource-group-name"] != "my-resource-group" {
		t.Errorf("inputs[resource-group-name] = %v, want my-resource-group", inputs["resource-group-name"])
	}
}

func TestAzureWebappsDeploy_Inputs_TargetPath(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName:    "my-web-app",
		Package:    "./dist",
		TargetPath: "/home/site/wwwroot",
	}

	inputs := a.Inputs()

	if inputs["target-path"] != "/home/site/wwwroot" {
		t.Errorf("inputs[target-path] = %v, want /home/site/wwwroot", inputs["target-path"])
	}
}

func TestAzureWebappsDeploy_Inputs_BooleanFalse(t *testing.T) {
	a := AzureWebappsDeploy{
		AppName: "my-web-app",
		Package: "./dist",
		Clean:   false,
		Restart: false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["clean"]; ok {
		t.Error("clean=false should not be in inputs")
	}
	if _, ok := inputs["restart"]; ok {
		t.Error("restart=false should not be in inputs")
	}
}

func TestAzureWebappsDeploy_ImplementsStepAction(t *testing.T) {
	a := AzureWebappsDeploy{}
	var _ workflow.StepAction = a
}
