package main

// validProviders is the list of supported LLM providers.
var validProviders = []string{"anthropic", "kiro"}

// isValidProvider checks if a provider name is valid.
func isValidProvider(provider string) bool {
	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}
	return false
}
