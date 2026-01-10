// Package setup_java provides a typed wrapper for actions/setup-java.
package setup_java

// SetupJava wraps the actions/setup-java@v4 action.
// Set up a specific version of the Java JDK and add it to PATH.
type SetupJava struct {
	// The Java version to set up. Required.
	// Examples: "11", "17", "21", "17.0.x", "21-ea"
	JavaVersion string `yaml:"java-version,omitempty"`

	// Java distribution. Required.
	// Supported: "temurin", "zulu", "adopt", "liberica", "microsoft", "corretto", "semeru", "oracle", "dragonwell"
	Distribution string `yaml:"distribution,omitempty"`

	// Path to the .java-version file
	JavaVersionFile string `yaml:"java-version-file,omitempty"`

	// The package type (jdk, jre, jdk+fx, jre+fx). Default: jdk
	JavaPackage string `yaml:"java-package,omitempty"`

	// The architecture of the package (defaults to the action runner's architecture)
	Architecture string `yaml:"architecture,omitempty"`

	// Path to where the compressed JDK is located
	JdkFile string `yaml:"jdk-file,omitempty"`

	// Set this option to true if you want the action to check for the latest available version
	CheckLatest bool `yaml:"check-latest,omitempty"`

	// ID of the distributionManagement repository in pom.xml
	ServerID string `yaml:"server-id,omitempty"`

	// Environment variable name for the username for authentication to the Apache Maven repository
	ServerUsername string `yaml:"server-username,omitempty"`

	// Environment variable name for password for authentication to the Apache Maven repository
	ServerPassword string `yaml:"server-password,omitempty"`

	// Path to where the settings.xml file will be written
	SettingsPath string `yaml:"settings-path,omitempty"`

	// Overwrite the settings.xml file if it exists
	OverwriteSettings bool `yaml:"overwrite-settings,omitempty"`

	// GPG private key to import
	GPGPrivateKey string `yaml:"gpg-private-key,omitempty"`

	// Environment variable name for the GPG private key passphrase
	GPGPassphrase string `yaml:"gpg-passphrase,omitempty"`

	// Used to specify whether caching is needed. Set to true if you'd like to enable caching
	Cache string `yaml:"cache,omitempty"`

	// Used to specify the path to a dependency file: pom.xml, build.gradle, etc.
	CacheDependencyPath string `yaml:"cache-dependency-path,omitempty"`

	// Path to where the compressed JDK is located
	Token string `yaml:"token,omitempty"`

	// Name of a target release of Maven Toolchains
	MvnToolchainID string `yaml:"mvn-toolchain-id,omitempty"`

	// Name of Maven Toolchain Vendor
	MvnToolchainVendor string `yaml:"mvn-toolchain-vendor,omitempty"`
}

// Action returns the action reference.
func (a SetupJava) Action() string {
	return "actions/setup-java@v4"
}

// Inputs returns the action inputs as a map.
func (a SetupJava) Inputs() map[string]any {
	with := make(map[string]any)

	if a.JavaVersion != "" {
		with["java-version"] = a.JavaVersion
	}
	if a.Distribution != "" {
		with["distribution"] = a.Distribution
	}
	if a.JavaVersionFile != "" {
		with["java-version-file"] = a.JavaVersionFile
	}
	if a.JavaPackage != "" {
		with["java-package"] = a.JavaPackage
	}
	if a.Architecture != "" {
		with["architecture"] = a.Architecture
	}
	if a.JdkFile != "" {
		with["jdk-file"] = a.JdkFile
	}
	if a.CheckLatest {
		with["check-latest"] = a.CheckLatest
	}
	if a.ServerID != "" {
		with["server-id"] = a.ServerID
	}
	if a.ServerUsername != "" {
		with["server-username"] = a.ServerUsername
	}
	if a.ServerPassword != "" {
		with["server-password"] = a.ServerPassword
	}
	if a.SettingsPath != "" {
		with["settings-path"] = a.SettingsPath
	}
	if a.OverwriteSettings {
		with["overwrite-settings"] = a.OverwriteSettings
	}
	if a.GPGPrivateKey != "" {
		with["gpg-private-key"] = a.GPGPrivateKey
	}
	if a.GPGPassphrase != "" {
		with["gpg-passphrase"] = a.GPGPassphrase
	}
	if a.Cache != "" {
		with["cache"] = a.Cache
	}
	if a.CacheDependencyPath != "" {
		with["cache-dependency-path"] = a.CacheDependencyPath
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.MvnToolchainID != "" {
		with["mvn-toolchain-id"] = a.MvnToolchainID
	}
	if a.MvnToolchainVendor != "" {
		with["mvn-toolchain-vendor"] = a.MvnToolchainVendor
	}

	return with
}
