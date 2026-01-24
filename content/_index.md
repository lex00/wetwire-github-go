---
title: "Wetwire GitHub"
---

[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-github-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-github-go)
[![CI](https://github.com/lex00/wetwire-github-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-github-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-github-go/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-github-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-github-go)](https://goreportcard.com/report/github.com/lex00/wetwire-github-go)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate GitHub Actions workflows from Go structs with AI-assisted design.

## Philosophy

Wetwire uses typed constraints to reduce the model capability required for accurate code generation.

**Core hypothesis:** Typed input + smaller model ≈ Semantic input + larger model

The type system and lint rules act as a force multiplier — cheaper models produce quality output when guided by schema-generated types and iterative lint feedback.

## Documentation

| Document | Description |
|----------|-------------|
| [CLI Reference]({{< relref "/cli" >}}) | Command-line interface |
| [Quick Start]({{< relref "/quick-start" >}}) | Get started in 5 minutes |
| [Examples]({{< relref "/examples" >}}) | Sample workflow projects |
| [FAQ]({{< relref "/faq" >}}) | Frequently asked questions |

## Installation

```bash
go install github.com/lex00/wetwire-github-go@latest
```

## Quick Example

```go
var CI = workflow.Workflow{
    Name: "CI",
    On:   workflow.OnPush{Branches: []string{"main"}},
    Jobs: []workflow.Job{BuildJob},
}
```
