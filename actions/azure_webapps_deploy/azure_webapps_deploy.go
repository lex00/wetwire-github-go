// Package azure_webapps_deploy provides a typed wrapper for azure/webapps-deploy.
package azure_webapps_deploy

// AzureWebappsDeploy wraps the azure/webapps-deploy@v3 action.
// Deploy to Azure Web Apps or Azure Web App for Containers.
type AzureWebappsDeploy struct {
	// Name of the Azure Web App (required).
	AppName string `yaml:"app-name,omitempty"`

	// Publish profile for authentication (alternative to azure/login).
	PublishProfile string `yaml:"publish-profile,omitempty"`

	// Deployment slot name.
	SlotName string `yaml:"slot-name,omitempty"`

	// Path to package or folder for Web App deployment.
	Package string `yaml:"package,omitempty"`

	// Container image(s) for Web App Containers.
	Images string `yaml:"images,omitempty"`

	// Docker-Compose file path for multi-container deployment.
	ConfigurationFile string `yaml:"configuration-file,omitempty"`

	// Startup command for the app.
	StartupCommand string `yaml:"startup-command,omitempty"`

	// Resource group name of the web app.
	ResourceGroupName string `yaml:"resource-group-name,omitempty"`

	// Deployment type: JAR, WAR, EAR, ZIP, Static.
	Type string `yaml:"type,omitempty"`

	// Target path in the web app.
	TargetPath string `yaml:"target-path,omitempty"`

	// Delete existing files before deploying.
	Clean bool `yaml:"clean,omitempty"`

	// Restart the app service after deployment.
	Restart bool `yaml:"restart,omitempty"`
}

// Action returns the action reference.
func (a AzureWebappsDeploy) Action() string {
	return "azure/webapps-deploy@v3"
}

// Inputs returns the action inputs as a map.
func (a AzureWebappsDeploy) Inputs() map[string]any {
	with := make(map[string]any)

	if a.AppName != "" {
		with["app-name"] = a.AppName
	}
	if a.PublishProfile != "" {
		with["publish-profile"] = a.PublishProfile
	}
	if a.SlotName != "" {
		with["slot-name"] = a.SlotName
	}
	if a.Package != "" {
		with["package"] = a.Package
	}
	if a.Images != "" {
		with["images"] = a.Images
	}
	if a.ConfigurationFile != "" {
		with["configuration-file"] = a.ConfigurationFile
	}
	if a.StartupCommand != "" {
		with["startup-command"] = a.StartupCommand
	}
	if a.ResourceGroupName != "" {
		with["resource-group-name"] = a.ResourceGroupName
	}
	if a.Type != "" {
		with["type"] = a.Type
	}
	if a.TargetPath != "" {
		with["target-path"] = a.TargetPath
	}
	if a.Clean {
		with["clean"] = a.Clean
	}
	if a.Restart {
		with["restart"] = a.Restart
	}

	return with
}
