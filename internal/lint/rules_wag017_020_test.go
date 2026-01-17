package lint

import (
	"strings"
	"testing"
)

// WAG017 Tests - Suggest workflow permissions scope

func TestWAG017_Check_MissingPermissions(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	On:   workflow.Triggers{},
}
`)

	l := NewLinter(&WAG017{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG017 should have flagged missing Permissions")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG017" {
			found = true
			if issue.Severity != SeverityInfo {
				t.Error("WAG017 issues should be severity 'info'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG017 issue not found")
	}
}

func TestWAG017_Check_HasPermissions(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name:        "CI",
	On:          workflow.Triggers{},
	Permissions: workflow.Permissions{Contents: "read"},
}
`)

	l := NewLinter(&WAG017{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG017 should not flag when Permissions is set")
	}
}

// WAG019 - Circular Dependency Detection Tests

func TestWAG019_Check_SimpleCycle(t *testing.T) {
	// A -> B -> A (simple 2-job cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected circular dependency between JobA and JobB")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
			if issue.Severity != SeverityError {
				t.Errorf("WAG019 issues should be severity 'error', got %s", issue.Severity)
			}
			if !strings.Contains(issue.Message, "JobA") || !strings.Contains(issue.Message, "JobB") {
				t.Errorf("WAG019 message should contain job names, got: %s", issue.Message)
			}
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_ThreeJobCycle(t *testing.T) {
	// A -> B -> C -> A (3-job cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobC},
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected circular dependency in 3-job cycle")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_SelfReference(t *testing.T) {
	// A -> A (self-referencing cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected self-referencing circular dependency")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
			if !strings.Contains(issue.Message, "JobA") {
				t.Errorf("WAG019 message should contain job name, got: %s", issue.Message)
			}
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_NoCycle(t *testing.T) {
	// Linear chain: A -> B -> C (no cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Unexpected issue: %s", issue.Message)
		}
		t.Error("WAG019 should not flag linear dependency chain")
	}
}

func TestWAG019_Check_DiamondDependency(t *testing.T) {
	// Diamond shape: A -> B, A -> C, B -> D, C -> D (no cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobD = workflow.Job{
	Name:   "job-d",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB, JobC},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Unexpected issue: %s", issue.Message)
		}
		t.Error("WAG019 should not flag diamond dependency pattern (no cycle)")
	}
}

func TestWAG019_Check_SingleDependencyFormat(t *testing.T) {
	// Test with single dependency (not slice): A -> B -> A
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  JobB,
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  JobA,
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should detect cycle with single dependency format")
	}
}

func TestWAG019_Check_NoJobs(t *testing.T) {
	// No jobs at all
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG019 should not flag when there are no jobs")
	}
}

// WAG020 - Secret Pattern Detection Tests

func TestWAG020_Check_AWSAccessKey(t *testing.T) {
	content := []byte(`package main

var awsKey = "AKIAIOSFODNN7EXAMPLE"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect AWS access key pattern")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG020" {
			found = true
			if issue.Severity != SeverityError {
				t.Errorf("WAG020 issues should be severity 'error', got %s", issue.Severity)
			}
			if !strings.Contains(issue.Message, "AWS") {
				t.Errorf("WAG020 message should mention AWS, got: %s", issue.Message)
			}
		}
	}
	if !found {
		t.Error("Expected WAG020 issue not found")
	}
}

func TestWAG020_Check_PrivateKey(t *testing.T) {
	content := []byte(`package main

var privateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIE..."
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect private key header")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG020" && strings.Contains(issue.Message, "private key") {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG020 private key issue not found")
	}
}

func TestWAG020_Check_StripeKey(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	stripePrefix := "sk_" + "live_"
	content := []byte(`package main

var stripeKey = "` + stripePrefix + `51H7xxxxxxxxxxxxxxxxxxxxxxxF"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect Stripe secret key")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG020" && strings.Contains(issue.Message, "Stripe") {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG020 Stripe key issue not found")
	}
}

func TestWAG020_Check_SlackToken(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	slackPrefix := "xox" + "b-"
	content := []byte(`package main

var slackToken = "` + slackPrefix + `1234567890123-1234567890123-abcdefghijklmnopqrstuvwx"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect Slack token")
	}
}

func TestWAG020_Check_GoogleAPIKey(t *testing.T) {
	content := []byte(`package main

var googleKey = "AIzaSyD1234567890abcdefghijklmnopqrstuvwx"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect Google API key")
	}
}

func TestWAG020_Check_GitHubOAuthToken(t *testing.T) {
	content := []byte(`package main

var ghOAuth = "gho_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect GitHub OAuth token")
	}
}

func TestWAG020_Check_SendGridKey(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	sgPrefix := "SG" + "."
	content := []byte(`package main

var sendgridKey = "` + sgPrefix + `abcdefghijklmnopqrstuv.abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect SendGrid API key")
	}
}

func TestWAG020_Check_TwilioKey(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	twilioPrefix := "S" + "K"
	content := []byte(`package main

var twilioKey = "` + twilioPrefix + `1234567890abcdef1234567890abcdef"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect Twilio API key")
	}
}

func TestWAG020_Check_NPMToken(t *testing.T) {
	content := []byte(`package main

var npmToken = "npm_abcdefghijklmnopqrstuvwxyz1234567890"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect NPM token")
	}
}

func TestWAG020_Check_PyPIToken(t *testing.T) {
	content := []byte(`package main

var pypiToken = "pypi-AgEIcHlwaS5vcmcxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect PyPI token")
	}
}

func TestWAG020_Check_JWTToken(t *testing.T) {
	content := []byte(`package main

var jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect JWT token")
	}
}

func TestWAG020_Check_NoSecrets(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}

var normalString = "hello world"
var envVar = "${{ secrets.MY_SECRET }}"
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Unexpected issue: %s", issue.Message)
		}
		t.Error("WAG020 should not flag normal strings or secrets references")
	}
}

func TestWAG020_Check_InEnvValue(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	stripeTestPrefix := "sk_" + "test_"
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	Env: map[string]any{
		"STRIPE_KEY": "` + stripeTestPrefix + `51H7xxxxxxxxxxxxxxxxxxxF",
	},
}
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect secrets in Env map values")
	}
}

func TestWAG020_Check_InRunCommand(t *testing.T) {
	// Using string concatenation to avoid GitHub secret scanning
	stripeLivePrefix := "sk_" + "live_"
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	Run: "curl -H 'Authorization: Bearer ` + stripeLivePrefix + `1234567890abcdefghijklmnop' https://api.stripe.com",
}
`)

	l := NewLinter(&WAG020{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG020 should detect secrets in Run commands")
	}
}
