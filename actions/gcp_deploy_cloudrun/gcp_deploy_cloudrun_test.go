package gcp_deploy_cloudrun

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGCPDeployCloudRun_Action(t *testing.T) {
	a := GCPDeployCloudRun{}
	if got := a.Action(); got != "google-github-actions/deploy-cloudrun@v2" {
		t.Errorf("Action() = %q, want %q", got, "google-github-actions/deploy-cloudrun@v2")
	}
}

func TestGCPDeployCloudRun_Inputs(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		Region:  "us-central1",
		Image:   "gcr.io/my-project/my-image:latest",
	}

	inputs := a.Inputs()

	if inputs["service"] != "my-service" {
		t.Errorf("inputs[service] = %v, want my-service", inputs["service"])
	}
	if inputs["region"] != "us-central1" {
		t.Errorf("inputs[region] = %v, want us-central1", inputs["region"])
	}
	if inputs["image"] != "gcr.io/my-project/my-image:latest" {
		t.Errorf("inputs[image] = %v, want image reference", inputs["image"])
	}
}

func TestGCPDeployCloudRun_Inputs_Empty(t *testing.T) {
	a := GCPDeployCloudRun{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty GCPDeployCloudRun.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGCPDeployCloudRun_Inputs_EnvVars(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		EnvVars: "KEY1=value1,KEY2=value2",
		Secrets: "SECRET1=projects/123/secrets/my-secret:latest",
	}

	inputs := a.Inputs()

	if inputs["env_vars"] != "KEY1=value1,KEY2=value2" {
		t.Errorf("inputs[env_vars] = %v, want env vars string", inputs["env_vars"])
	}
	if inputs["secrets"] != "SECRET1=projects/123/secrets/my-secret:latest" {
		t.Errorf("inputs[secrets] = %v, want secrets string", inputs["secrets"])
	}
}

func TestGCPDeployCloudRun_Inputs_Options(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:   "my-service",
		NoTraffic: true,
		Tag:       "blue",
		Timeout:   "300s",
	}

	inputs := a.Inputs()

	if inputs["no_traffic"] != true {
		t.Errorf("inputs[no_traffic] = %v, want true", inputs["no_traffic"])
	}
	if inputs["tag"] != "blue" {
		t.Errorf("inputs[tag] = %v, want blue", inputs["tag"])
	}
	if inputs["timeout"] != "300s" {
		t.Errorf("inputs[timeout] = %v, want 300s", inputs["timeout"])
	}
}

func TestGCPDeployCloudRun_ImplementsStepAction(t *testing.T) {
	a := GCPDeployCloudRun{}
	var _ workflow.StepAction = a
}
