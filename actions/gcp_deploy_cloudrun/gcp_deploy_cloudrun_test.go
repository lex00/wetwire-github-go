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

func TestGCPDeployCloudRun_Inputs_Job(t *testing.T) {
	a := GCPDeployCloudRun{
		Job:   "my-job",
		Image: "gcr.io/my-project/my-image:latest",
	}

	inputs := a.Inputs()

	if inputs["job"] != "my-job" {
		t.Errorf("inputs[job] = %v, want my-job", inputs["job"])
	}
	if inputs["image"] != "gcr.io/my-project/my-image:latest" {
		t.Errorf("inputs[image] = %v, want image reference", inputs["image"])
	}
}

func TestGCPDeployCloudRun_Inputs_Source(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		Source:  "./src",
	}

	inputs := a.Inputs()

	if inputs["source"] != "./src" {
		t.Errorf("inputs[source] = %v, want ./src", inputs["source"])
	}
}

func TestGCPDeployCloudRun_Inputs_Metadata(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:  "my-service",
		Metadata: "./service.yaml",
	}

	inputs := a.Inputs()

	if inputs["metadata"] != "./service.yaml" {
		t.Errorf("inputs[metadata] = %v, want ./service.yaml", inputs["metadata"])
	}
}

func TestGCPDeployCloudRun_Inputs_EnvVarsUpdateStrategy(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:               "my-service",
		EnvVars:               "KEY=value",
		EnvVarsUpdateStrategy: "merge",
	}

	inputs := a.Inputs()

	if inputs["env_vars_update_strategy"] != "merge" {
		t.Errorf("inputs[env_vars_update_strategy] = %v, want merge", inputs["env_vars_update_strategy"])
	}
}

func TestGCPDeployCloudRun_Inputs_SecretsUpdateStrategy(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:               "my-service",
		Secrets:               "SECRET=my-secret",
		SecretsUpdateStrategy: "overwrite",
	}

	inputs := a.Inputs()

	if inputs["secrets_update_strategy"] != "overwrite" {
		t.Errorf("inputs[secrets_update_strategy] = %v, want overwrite", inputs["secrets_update_strategy"])
	}
}

func TestGCPDeployCloudRun_Inputs_Labels(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		Labels:  "env=prod,team=backend",
	}

	inputs := a.Inputs()

	if inputs["labels"] != "env=prod,team=backend" {
		t.Errorf("inputs[labels] = %v, want env=prod,team=backend", inputs["labels"])
	}
}

func TestGCPDeployCloudRun_Inputs_Flags(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		Flags:   "--memory=512Mi --cpu=1",
	}

	inputs := a.Inputs()

	if inputs["flags"] != "--memory=512Mi --cpu=1" {
		t.Errorf("inputs[flags] = %v, want --memory=512Mi --cpu=1", inputs["flags"])
	}
}

func TestGCPDeployCloudRun_Inputs_ProjectID(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:   "my-service",
		ProjectID: "my-gcp-project",
	}

	inputs := a.Inputs()

	if inputs["project_id"] != "my-gcp-project" {
		t.Errorf("inputs[project_id] = %v, want my-gcp-project", inputs["project_id"])
	}
}

func TestGCPDeployCloudRun_Inputs_Suffix(t *testing.T) {
	a := GCPDeployCloudRun{
		Service: "my-service",
		Suffix:  "v1",
	}

	inputs := a.Inputs()

	if inputs["suffix"] != "v1" {
		t.Errorf("inputs[suffix] = %v, want v1", inputs["suffix"])
	}
}

func TestGCPDeployCloudRun_Inputs_SkipDefaultLabels(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:           "my-service",
		SkipDefaultLabels: true,
	}

	inputs := a.Inputs()

	if inputs["skip_default_labels"] != true {
		t.Errorf("inputs[skip_default_labels] = %v, want true", inputs["skip_default_labels"])
	}
}

func TestGCPDeployCloudRun_Inputs_AllFields(t *testing.T) {
	a := GCPDeployCloudRun{
		Service:               "my-service",
		Job:                   "my-job",
		Image:                 "gcr.io/my-project/my-image:latest",
		Source:                "./src",
		Metadata:              "./service.yaml",
		EnvVars:               "KEY1=value1,KEY2=value2",
		EnvVarsUpdateStrategy: "merge",
		Secrets:               "SECRET1=projects/123/secrets/my-secret:latest",
		SecretsUpdateStrategy: "overwrite",
		Labels:                "env=prod,team=backend",
		Tag:                   "blue",
		Timeout:               "300s",
		Flags:                 "--memory=512Mi --cpu=1",
		NoTraffic:             true,
		ProjectID:             "my-gcp-project",
		Region:                "us-central1",
		Suffix:                "v1",
		SkipDefaultLabels:     true,
	}

	inputs := a.Inputs()

	// Verify all fields are present
	expectedFields := []string{
		"service", "job", "image", "source", "metadata",
		"env_vars", "env_vars_update_strategy",
		"secrets", "secrets_update_strategy",
		"labels", "tag", "timeout", "flags",
		"no_traffic", "project_id", "region", "suffix",
		"skip_default_labels",
	}

	for _, field := range expectedFields {
		if _, ok := inputs[field]; !ok {
			t.Errorf("inputs[%s] not present, want it to be set", field)
		}
	}

	if len(inputs) != len(expectedFields) {
		t.Errorf("inputs has %d fields, want %d", len(inputs), len(expectedFields))
	}
}
