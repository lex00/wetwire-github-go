package importer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IRCodeowners represents a parsed CODEOWNERS file.
type IRCodeowners struct {
	Rules []IRCodeownersRule
}

// IRCodeownersRule represents a single ownership rule.
type IRCodeownersRule struct {
	Pattern string
	Owners  []string
	Comment string
}

// ParseCodeownersFile parses a CODEOWNERS file from disk.
func ParseCodeownersFile(path string) (*IRCodeowners, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return ParseCodeownersContent(string(content))
}

// ParseCodeownersContent parses CODEOWNERS content from a string.
func ParseCodeownersContent(content string) (*IRCodeowners, error) {
	ir := &IRCodeowners{
		Rules: []IRCodeownersRule{},
	}

	if content == "" {
		return ir, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var pendingComment string

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		}

		// Handle full-line comments
		if strings.HasPrefix(trimmedLine, "#") {
			// Store as pending comment for next rule
			commentText := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "#"))
			pendingComment = commentText
			continue
		}

		// Parse rule line
		rule, err := parseCodeownersLine(trimmedLine)
		if err != nil {
			return nil, fmt.Errorf("parsing line %q: %w", trimmedLine, err)
		}

		// Attach pending comment if present
		if pendingComment != "" && rule.Comment == "" {
			rule.Comment = pendingComment
			pendingComment = ""
		}

		ir.Rules = append(ir.Rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning content: %w", err)
	}

	return ir, nil
}

// parseCodeownersLine parses a single CODEOWNERS rule line.
func parseCodeownersLine(line string) (IRCodeownersRule, error) {
	rule := IRCodeownersRule{}

	// Check for inline comment
	inlineCommentIdx := strings.Index(line, " #")
	if inlineCommentIdx > 0 {
		rule.Comment = strings.TrimSpace(line[inlineCommentIdx+2:])
		line = strings.TrimSpace(line[:inlineCommentIdx])
	}

	// Split on whitespace
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return rule, fmt.Errorf("invalid rule: expected pattern and at least one owner")
	}

	rule.Pattern = parts[0]
	rule.Owners = parts[1:]

	return rule, nil
}

// CodeownersCodeGenerator generates Go code from parsed CODEOWNERS.
type CodeownersCodeGenerator struct {
	PackageName string
}

// CodeownersGeneratedCode contains the generated Go code.
type CodeownersGeneratedCode struct {
	Files map[string]string // filename -> content
	Rules int
}

// Generate generates Go code from parsed CODEOWNERS.
func (g *CodeownersCodeGenerator) Generate(ir *IRCodeowners) (*CodeownersGeneratedCode, error) {
	result := &CodeownersGeneratedCode{
		Files: make(map[string]string),
		Rules: len(ir.Rules),
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/codeowners\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("var Codeowners = codeowners.Owners{\n")
	if len(ir.Rules) == 0 {
		sb.WriteString("\tRules: []codeowners.Rule{},\n")
	} else {
		sb.WriteString("\tRules: []codeowners.Rule{\n")
		for _, rule := range ir.Rules {
			sb.WriteString(g.generateRule(rule))
		}
		sb.WriteString("\t},\n")
	}
	sb.WriteString("}\n")

	result.Files["codeowners.go"] = sb.String()
	return result, nil
}

func (g *CodeownersCodeGenerator) generateRule(rule IRCodeownersRule) string {
	var sb strings.Builder
	sb.WriteString("\t\t{\n")
	sb.WriteString(fmt.Sprintf("\t\t\tPattern: %q,\n", rule.Pattern))

	// Generate Owners slice
	sb.WriteString("\t\t\tOwners: []string{")
	for i, owner := range rule.Owners {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q", owner))
	}
	sb.WriteString("},\n")

	// Add comment if present
	if rule.Comment != "" {
		sb.WriteString(fmt.Sprintf("\t\t\tComment: %q,\n", rule.Comment))
	}

	sb.WriteString("\t\t},\n")
	return sb.String()
}
