package importer

import (
	"testing"
)

// Realistic workflow YAML for benchmarking
var benchmarkWorkflowYAML = []byte(`name: CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'src/**'
      - '*.go'
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

env:
  GO_VERSION: '1.23'
  CI: true

permissions:
  contents: read
  pull-requests: write

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go build ./...

  test:
    name: Test
    runs-on: ubuntu-latest
    needs: build
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test -v ./...

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest
`)

// Realistic job YAML for benchmarking
var benchmarkJobYAML = []byte(`name: Test Workflow

on: push

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    timeout-minutes: 30
    env:
      GOPATH: /home/runner/go
      GOCACHE: /home/runner/.cache/go-build
    strategy:
      matrix:
        go: ['1.22', '1.23']
        os: [ubuntu-latest, macos-latest, windows-latest]
      fail-fast: false
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          submodules: recursive
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
          cache: true
      - name: Run Tests
        run: go test -v -race -coverprofile=coverage.txt ./...
        env:
          GOPROXY: https://proxy.golang.org
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: coverage.txt
`)

// Complex workflow with many features for benchmarking
var benchmarkComplexWorkflowYAML = []byte(`name: Complex CI

on:
  push:
    branches: [main, develop, 'release/*']
    tags: ['v*']
    paths-ignore:
      - 'docs/**'
      - '*.md'
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 0 * * *'
  workflow_dispatch:
    inputs:
      environment:
        type: string
        required: true
        description: Deployment environment

env:
  GO_VERSION: '1.23'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

permissions:
  contents: read
  packages: write
  id-token: write

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

defaults:
  run:
    shell: bash
    working-directory: ./src

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4

  test:
    name: Test
    runs-on: ${{ matrix.os }}
    needs: lint
    strategy:
      matrix:
        go: ['1.22', '1.23']
        os: [ubuntu-latest, macos-latest, windows-latest]
        exclude:
          - os: windows-latest
            go: '1.22'
        include:
          - os: ubuntu-latest
            go: '1.23'
            coverage: true
      fail-fast: false
      max-parallel: 4
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test -v ./...
      - if: ${{ matrix.coverage }}
        run: go test -coverprofile=coverage.txt ./...

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, test]
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - id: version
        run: echo "version=$(git describe --tags --always)" >> $GITHUB_OUTPUT
      - run: go build -ldflags="-s -w" -o ./bin/app ./cmd/app
      - uses: actions/upload-artifact@v4
        with:
          name: app
          path: ./bin/app

  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    needs: build
    if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ needs.build.outputs.version }}

  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: [build, docker]
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://app.example.com
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: app
      - run: ./deploy.sh
        env:
          DEPLOY_TOKEN: ${{ secrets.DEPLOY_TOKEN }}
`)

// BenchmarkImportWorkflow benchmarks YAML to Go code import.
func BenchmarkImportWorkflow(b *testing.B) {
	parser := NewParser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		workflow, err := parser.Parse(benchmarkWorkflowYAML)
		if err != nil {
			b.Fatal(err)
		}
		if workflow.Name == "" {
			b.Fatal("expected workflow name")
		}
	}
}

// BenchmarkImportJob benchmarks job import.
func BenchmarkImportJob(b *testing.B) {
	parser := NewParser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		workflow, err := parser.Parse(benchmarkJobYAML)
		if err != nil {
			b.Fatal(err)
		}
		if len(workflow.Jobs) == 0 {
			b.Fatal("expected jobs")
		}
	}
}

// BenchmarkImportComplexWorkflow benchmarks complex workflow import.
func BenchmarkImportComplexWorkflow(b *testing.B) {
	parser := NewParser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		workflow, err := parser.Parse(benchmarkComplexWorkflowYAML)
		if err != nil {
			b.Fatal(err)
		}
		if len(workflow.Jobs) == 0 {
			b.Fatal("expected jobs")
		}
	}
}

// BenchmarkParseWorkflow is an alias for BenchmarkImportWorkflow for consistency.
func BenchmarkParseWorkflow(b *testing.B) {
	parser := NewParser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		workflow, err := parser.Parse(benchmarkWorkflowYAML)
		if err != nil {
			b.Fatal(err)
		}
		if workflow.Name == "" {
			b.Fatal("expected workflow name")
		}
	}
}

// BenchmarkBuildReferenceGraph benchmarks building a reference graph.
func BenchmarkBuildReferenceGraph(b *testing.B) {
	parser := NewParser()
	workflow, err := parser.Parse(benchmarkComplexWorkflowYAML)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		graph := BuildReferenceGraph(workflow)
		if len(graph.JobDependencies) == 0 {
			b.Fatal("expected job dependencies")
		}
	}
}

// BenchmarkCodeGenerator benchmarks Go code generation.
func BenchmarkCodeGenerator(b *testing.B) {
	parser := NewParser()
	workflow, err := parser.Parse(benchmarkWorkflowYAML)
	if err != nil {
		b.Fatal(err)
	}

	gen := NewCodeGenerator()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		code, err := gen.Generate(workflow, "ci")
		if err != nil {
			b.Fatal(err)
		}
		if len(code.Files) == 0 {
			b.Fatal("expected generated files")
		}
	}
}

// BenchmarkCodeGeneratorSingleFile benchmarks single-file code generation.
func BenchmarkCodeGeneratorSingleFile(b *testing.B) {
	parser := NewParser()
	workflow, err := parser.Parse(benchmarkWorkflowYAML)
	if err != nil {
		b.Fatal(err)
	}

	gen := NewCodeGenerator()
	gen.SingleFile = true

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		code, err := gen.Generate(workflow, "ci")
		if err != nil {
			b.Fatal(err)
		}
		if len(code.Files) == 0 {
			b.Fatal("expected generated files")
		}
	}
}

// BenchmarkCodeGeneratorComplexWorkflow benchmarks code generation for complex workflows.
func BenchmarkCodeGeneratorComplexWorkflow(b *testing.B) {
	parser := NewParser()
	workflow, err := parser.Parse(benchmarkComplexWorkflowYAML)
	if err != nil {
		b.Fatal(err)
	}

	gen := NewCodeGenerator()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		code, err := gen.Generate(workflow, "complex")
		if err != nil {
			b.Fatal(err)
		}
		if len(code.Files) == 0 {
			b.Fatal("expected generated files")
		}
	}
}

// BenchmarkParseMinimalWorkflow benchmarks parsing a minimal workflow.
func BenchmarkParseMinimalWorkflow(b *testing.B) {
	minimalYAML := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo hello
`)

	parser := NewParser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		workflow, err := parser.Parse(minimalYAML)
		if err != nil {
			b.Fatal(err)
		}
		if workflow.Name == "" {
			b.Fatal("expected workflow name")
		}
	}
}
