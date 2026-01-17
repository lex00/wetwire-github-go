package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// WAG011 detects potential unreachable jobs.
type WAG011 struct{}

func (r *WAG011) ID() string          { return "WAG011" }
func (r *WAG011) Description() string { return "Detect potential unreachable jobs" }

func (r *WAG011) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	jobDeps := make(map[string][]string)
	jobPositions := make(map[string]token.Position)

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if len(valueSpec.Values) <= i {
					continue
				}

				lit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				typeName := getTypeName(lit.Type)
				if typeName != "workflow.Job" && typeName != "Job" {
					continue
				}

				jobPositions[name.Name] = fset.Position(name.Pos())

				for _, elt := range lit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := kv.Key.(*ast.Ident)
					if !ok || key.Name != "Needs" {
						continue
					}

					deps := extractDependencyNames(kv.Value)
					jobDeps[name.Name] = deps
				}
			}
		}
		return true
	})

	definedJobs := make(map[string]bool)
	for name := range jobPositions {
		definedJobs[name] = true
	}

	for jobName, deps := range jobDeps {
		for _, dep := range deps {
			if !definedJobs[dep] {
				pos := jobPositions[jobName]
				issues = append(issues, LintIssue{
					File:     path,
					Line:     pos.Line,
					Column:   pos.Column,
					Severity: "error",
					Message:  fmt.Sprintf("Job %q depends on undefined job %q", jobName, dep),
					Rule:     r.ID(),
					Fixable:  false,
				})
			}
		}
	}

	return issues
}

// WAG019 detects circular dependencies in job dependency graphs.
type WAG019 struct{}

func (r *WAG019) ID() string { return "WAG019" }
func (r *WAG019) Description() string {
	return "Detect circular dependencies in job needs"
}

func (r *WAG019) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	jobDeps := make(map[string][]string)
	jobPositions := make(map[string]token.Position)

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if len(valueSpec.Values) <= i {
					continue
				}

				lit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				typeName := getTypeName(lit.Type)
				if typeName != "workflow.Job" && typeName != "Job" {
					continue
				}

				jobPositions[name.Name] = fset.Position(name.Pos())

				for _, elt := range lit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := kv.Key.(*ast.Ident)
					if !ok || key.Name != "Needs" {
						continue
					}

					deps := extractDependencyNames(kv.Value)
					jobDeps[name.Name] = deps
				}
			}
		}
		return true
	})

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	reportedCycles := make(map[string]bool)

	var detectCycle func(job string, currentPath []string) []string
	detectCycle = func(job string, currentPath []string) []string {
		visited[job] = true
		recStack[job] = true
		currentPath = append(currentPath, job)

		for _, dep := range jobDeps[job] {
			if recStack[dep] {
				cycleStart := -1
				for i, j := range currentPath {
					if j == dep {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					return currentPath[cycleStart:]
				}
				return []string{dep, job}
			}

			if !visited[dep] {
				if cycle := detectCycle(dep, currentPath); cycle != nil {
					return cycle
				}
			}
		}

		recStack[job] = false
		return nil
	}

	for job := range jobDeps {
		if !visited[job] {
			if cycle := detectCycle(job, nil); cycle != nil {
				cycleKey := normalizeCycle(cycle)
				if !reportedCycles[cycleKey] {
					reportedCycles[cycleKey] = true

					pos := jobPositions[cycle[0]]
					issues = append(issues, LintIssue{
						File:     path,
						Line:     pos.Line,
						Column:   pos.Column,
						Severity: "error",
						Message:  fmt.Sprintf("Circular dependency detected: %s", strings.Join(cycle, " -> ")+" -> "+cycle[0]),
						Rule:     r.ID(),
						Fixable:  false,
					})
				}
			}
		}
	}

	return issues
}
