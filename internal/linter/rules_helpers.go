package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

// actionInfo describes an action wrapper mapping.
type actionInfo struct {
	pkg        string
	typ        string
	importPath string
}

// knownActions maps action references to their wrapper types.
var knownActions = map[string]actionInfo{
	"actions/checkout":          {pkg: "checkout", typ: "Checkout", importPath: "github.com/lex00/wetwire-github-go/actions/checkout"},
	"actions/setup-go":          {pkg: "setup_go", typ: "SetupGo", importPath: "github.com/lex00/wetwire-github-go/actions/setup_go"},
	"actions/setup-node":        {pkg: "setup_node", typ: "SetupNode", importPath: "github.com/lex00/wetwire-github-go/actions/setup_node"},
	"actions/setup-python":      {pkg: "setup_python", typ: "SetupPython", importPath: "github.com/lex00/wetwire-github-go/actions/setup_python"},
	"actions/cache":             {pkg: "cache", typ: "Cache", importPath: "github.com/lex00/wetwire-github-go/actions/cache"},
	"actions/upload-artifact":   {pkg: "upload_artifact", typ: "UploadArtifact", importPath: "github.com/lex00/wetwire-github-go/actions/upload_artifact"},
	"actions/download-artifact": {pkg: "download_artifact", typ: "DownloadArtifact", importPath: "github.com/lex00/wetwire-github-go/actions/download_artifact"},
}

// recommendedInputs maps action types to their recommended inputs.
var recommendedInputs = map[string][]string{
	"setup_go.SetupGo":         {"GoVersion", "GoVersionFile"},
	"setup_java.SetupJava":     {"JavaVersion", "Distribution"},
	"setup_node.SetupNode":     {"NodeVersion", "NodeVersionFile"},
	"setup_python.SetupPython": {"PythonVersion"},
	"setup_dotnet.SetupDotnet": {"DotnetVersion"},
	"setup_ruby.SetupRuby":     {"RubyVersion"},
	"setup_rust.SetupRust":     {"Toolchain"},
}

// deprecatedVersions maps action patterns to their deprecated versions and recommended versions.
var deprecatedVersions = map[string]struct {
	deprecated  []string
	recommended string
}{
	"actions/checkout":          {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/setup-go":          {[]string{"v1", "v2", "v3", "v4"}, "v5"},
	"actions/setup-node":        {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/setup-python":      {[]string{"v1", "v2", "v3", "v4"}, "v5"},
	"actions/setup-java":        {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/cache":             {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/upload-artifact":   {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/download-artifact": {[]string{"v1", "v2", "v3"}, "v4"},
	"actions/setup-dotnet":      {[]string{"v1", "v2", "v3"}, "v4"},
}

// setupActionsNeedingCache maps setup action types to their display names.
var setupActionsNeedingCache = map[string]string{
	"setup_go.SetupGo":         "setup-go",
	"setup_node.SetupNode":     "setup-node",
	"setup_python.SetupPython": "setup-python",
	"SetupGo":                  "setup-go",
	"SetupNode":                "setup-node",
	"SetupPython":              "setup-python",
}

// secretPattern defines a secret detection pattern with a description.
type secretPattern struct {
	pattern     *regexp.Regexp
	description string
}

// secretPatterns contains all known secret patterns to detect.
var secretPatterns = []secretPattern{
	// AWS credentials
	{regexp.MustCompile(`AKIA[0-9A-Z]{16}`), "AWS access key"},
	{regexp.MustCompile(`(?i)(aws_secret_access_key|aws_secret_key)\s*[:=]\s*['"][A-Za-z0-9/+=]{40}['"]`), "AWS secret key"},

	// GitHub tokens
	{regexp.MustCompile(`ghp_[A-Za-z0-9]{36,}`), "GitHub personal access token"},
	{regexp.MustCompile(`ghs_[A-Za-z0-9]{36,}`), "GitHub server token"},
	{regexp.MustCompile(`ghu_[A-Za-z0-9]{36,}`), "GitHub user token"},
	{regexp.MustCompile(`ghr_[A-Za-z0-9]{36,}`), "GitHub refresh token"},
	{regexp.MustCompile(`gho_[A-Za-z0-9]{36,}`), "GitHub OAuth token"},
	{regexp.MustCompile(`github_pat_[A-Za-z0-9]{22,}`), "GitHub fine-grained PAT"},

	// Private keys
	{regexp.MustCompile(`-----BEGIN RSA PRIVATE KEY-----`), "RSA private key"},
	{regexp.MustCompile(`-----BEGIN PRIVATE KEY-----`), "private key"},
	{regexp.MustCompile(`-----BEGIN EC PRIVATE KEY-----`), "EC private key"},
	{regexp.MustCompile(`-----BEGIN DSA PRIVATE KEY-----`), "DSA private key"},
	{regexp.MustCompile(`-----BEGIN OPENSSH PRIVATE KEY-----`), "OpenSSH private key"},
	{regexp.MustCompile(`-----BEGIN PGP PRIVATE KEY BLOCK-----`), "PGP private key"},

	// Stripe keys
	{regexp.MustCompile(`sk_live_[A-Za-z0-9]{20,}`), "Stripe live secret key"},
	{regexp.MustCompile(`sk_test_[A-Za-z0-9]{20,}`), "Stripe test secret key"},
	{regexp.MustCompile(`rk_live_[A-Za-z0-9]{20,}`), "Stripe live restricted key"},
	{regexp.MustCompile(`rk_test_[A-Za-z0-9]{20,}`), "Stripe test restricted key"},

	// Slack tokens
	{regexp.MustCompile(`xox[baprs]-[0-9]{10,}-[0-9]{10,}-[A-Za-z0-9]{20,}`), "Slack token"},

	// Google API keys
	{regexp.MustCompile(`AIza[A-Za-z0-9_-]{35}`), "Google API key"},

	// Twilio
	{regexp.MustCompile(`SK[a-f0-9]{32}`), "Twilio API key"},
	{regexp.MustCompile(`AC[a-f0-9]{32}`), "Twilio Account SID"},

	// SendGrid
	{regexp.MustCompile(`SG\.[A-Za-z0-9_-]{22}\.[A-Za-z0-9_-]{43}`), "SendGrid API key"},

	// Mailgun
	{regexp.MustCompile(`key-[A-Za-z0-9]{32}`), "Mailgun API key"},

	// NPM tokens
	{regexp.MustCompile(`npm_[A-Za-z0-9]{36,}`), "NPM token"},

	// PyPI tokens
	{regexp.MustCompile(`pypi-[A-Za-z0-9]{50,}`), "PyPI token"},

	// JWT tokens (detect base64-encoded JWT structure)
	{regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`), "JWT token"},

	// Heroku
	{regexp.MustCompile(`[hH]eroku.*[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}`), "Heroku API key"},

	// DigitalOcean
	{regexp.MustCompile(`dop_v1_[A-Za-z0-9]{64}`), "DigitalOcean personal access token"},
	{regexp.MustCompile(`doo_v1_[A-Za-z0-9]{64}`), "DigitalOcean OAuth token"},

	// Azure
	{regexp.MustCompile(`(?i)azure[A-Za-z0-9_-]*['\"][A-Za-z0-9/+=]{40,}['\"]`), "Azure credential"},
}

// parseActionRef extracts the action name from a reference like "actions/checkout@v4".
func parseActionRef(ref string) string {
	if idx := strings.Index(ref, "@"); idx != -1 {
		ref = ref[:idx]
	}
	return ref
}

// getTypeName extracts the type name from an AST expression.
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	}
	return ""
}

// isWorkflowType checks if a type name is a workflow-related type.
func isWorkflowType(typeName string) bool {
	workflowTypes := map[string]bool{
		"workflow.Workflow":    true,
		"workflow.Job":         true,
		"workflow.Step":        true,
		"workflow.Strategy":    true,
		"workflow.Matrix":      true,
		"workflow.Triggers":    true,
		"workflow.Concurrency": true,
		"Workflow":             true,
		"Job":                  true,
		"Step":                 true,
		"Strategy":             true,
		"Matrix":               true,
		"Triggers":             true,
		"Concurrency":          true,
	}
	return workflowTypes[typeName]
}

// extractDependencyNames extracts job names from a Needs field value.
func extractDependencyNames(expr ast.Expr) []string {
	var names []string

	switch v := expr.(type) {
	case *ast.Ident:
		names = append(names, v.Name)
	case *ast.CompositeLit:
		for _, elt := range v.Elts {
			if ident, ok := elt.(*ast.Ident); ok {
				names = append(names, ident.Name)
			}
		}
	}

	return names
}

// addImportIfNeeded adds an import statement if it doesn't already exist.
func addImportIfNeeded(src []byte, importPath, alias string) []byte {
	srcStr := string(src)

	// Check if import already exists
	if strings.Contains(srcStr, fmt.Sprintf(`"%s"`, importPath)) {
		return src
	}

	// Find the import block
	importPattern := regexp.MustCompile(`import \(([^)]*)\)`)
	if match := importPattern.FindStringIndex(srcStr); match != nil {
		// Find the closing paren
		closeIdx := match[1] - 1
		newImport := fmt.Sprintf("\n\t\"%s\"", importPath)
		result := srcStr[:closeIdx] + newImport + srcStr[closeIdx:]
		return []byte(result)
	}

	// Try single import statement
	singleImportPattern := regexp.MustCompile(`import "([^"]*)"`)
	if match := singleImportPattern.FindStringIndex(srcStr); match != nil {
		// Convert to import block
		existingImport := singleImportPattern.FindString(srcStr)
		existingPath := strings.TrimPrefix(strings.TrimSuffix(existingImport, `"`), `import "`)
		newBlock := fmt.Sprintf("import (\n\t\"%s\"\n\t\"%s\"\n)", existingPath, importPath)
		result := srcStr[:match[0]] + newBlock + srcStr[match[1]:]
		return []byte(result)
	}

	return src
}

// checkTriggersForPRTarget checks if triggers contain pull_request_target.
func checkTriggersForPRTarget(expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return false
	case *ast.CompositeLit:
		typeName := getTypeName(v.Type)
		if typeName != "workflow.Triggers" && typeName != "Triggers" {
			return false
		}

		for _, elt := range v.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if key.Name == "PullRequestTarget" {
				if _, ok := kv.Value.(*ast.Ident); ok {
					return true
				}
				if unary, ok := kv.Value.(*ast.UnaryExpr); ok && unary.Op == token.AND {
					return true
				}
			}
		}
	}
	return false
}

// hasPullRequestTargetWithTracking checks if a workflow has pull_request_target trigger.
func hasPullRequestTargetWithTracking(workflowLit *ast.CompositeLit, triggersWithPRTarget map[string]bool) bool {
	for _, elt := range workflowLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "On" {
			continue
		}

		switch v := kv.Value.(type) {
		case *ast.Ident:
			return triggersWithPRTarget[v.Name]
		case *ast.CompositeLit:
			return checkTriggersForPRTarget(v)
		}
	}
	return false
}

// hasCheckoutAction checks if a workflow contains checkout action.
func hasCheckoutAction(workflowLit *ast.CompositeLit) bool {
	hasCheckout := false

	ast.Inspect(workflowLit, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName == "checkout.Checkout" || typeName == "Checkout" {
			hasCheckout = true
			return false
		}

		// Also check for raw Uses field with checkout
		if typeName == "workflow.Step" || typeName == "Step" {
			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || key.Name != "Uses" {
					continue
				}
				if bl, ok := kv.Value.(*ast.BasicLit); ok {
					usesVal := strings.Trim(bl.Value, `"'`)
					actionName := parseActionRef(usesVal)
					if actionName == "actions/checkout" {
						hasCheckout = true
						return false
					}
				}
			}
		}

		return true
	})

	return hasCheckout
}

// normalizeCycle creates a canonical representation of a cycle for deduplication.
func normalizeCycle(cycle []string) string {
	if len(cycle) == 0 {
		return ""
	}

	minIdx := 0
	for i, job := range cycle {
		if job < cycle[minIdx] {
			minIdx = i
		}
	}

	normalized := make([]string, len(cycle))
	for i := 0; i < len(cycle); i++ {
		normalized[i] = cycle[(minIdx+i)%len(cycle)]
	}

	return strings.Join(normalized, "->")
}
