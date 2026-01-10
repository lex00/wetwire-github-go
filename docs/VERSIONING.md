# Versioning

This document describes the versioning strategy for wetwire-github-go and its action wrappers.

## Semantic Versioning

wetwire-github-go follows [Semantic Versioning](https://semver.org/):

- **Major (X.0.0)**: Breaking changes to the API or behavior
- **Minor (0.X.0)**: New features, backward-compatible
- **Patch (0.0.X)**: Bug fixes, backward-compatible

## Action Wrapper Versions

Action wrappers target specific versions of GitHub Actions:

| Wrapper | Action Reference | Version |
|---------|-----------------|---------|
| `checkout.Checkout` | `actions/checkout@v4` | v4 |
| `setup_go.SetupGo` | `actions/setup-go@v5` | v5 |
| `setup_node.SetupNode` | `actions/setup-node@v4` | v4 |
| `cache.Cache` | `actions/cache@v4` | v4 |
| `docker_build_push.DockerBuildPush` | `docker/build-push-action@v6` | v6 |

### When Action Versions Update

When upstream actions release new major versions:

1. **New wrapper version**: We update the wrapper to target the new version
2. **Deprecation notice**: Previous version may be flagged by WAG012 linter rule
3. **Changelog entry**: Breaking changes documented in CHANGELOG.md

### Multiple Action Versions

If you need a specific action version different from the wrapper default:

```go
// Use the raw Step type with custom action reference
var CustomStep = workflow.Step{
    Uses: "actions/checkout@v3",  // Override to use v3 instead of v4
    With: map[string]any{
        "fetch-depth": 0,
    },
}
```

## Breaking Changes Policy

### What Constitutes a Breaking Change

- Removing a struct field from an action wrapper
- Changing the return type of `Action()` or `Inputs()` methods
- Changing the behavior of `Job.Steps` slice handling
- Removing or renaming exported types

### What Is Not a Breaking Change

- Adding new struct fields to action wrappers
- Adding new action wrappers
- Adding new lint rules (they can be ignored)
- Updating action version references (old code still works)

## Deprecation Timeline

When deprecating features:

1. **Announce**: Document in CHANGELOG.md with deprecation warning
2. **Warn**: Add lint rule to warn about deprecated usage
3. **Remove**: Remove in next major version (minimum 3 months)

### Currently Deprecated

- `.ToStep()` method on action wrappers (use wrappers directly in `Job.Steps`)
- `checkout@v2`, `checkout@v3` (use v4)
- `setup-go@v3`, `setup-go@v4` (use v5)

## Version Compatibility Matrix

| wetwire-github-go | Go Version | actions/checkout | actions/setup-go |
|-------------------|------------|------------------|------------------|
| 0.1.x | 1.22+ | v4 | v5 |
| 0.2.x (planned) | 1.23+ | v4 | v5 |

## Checking Your Version

```bash
wetwire-github version
```

## Upgrading

1. Update the module:
   ```bash
   go get -u github.com/lex00/wetwire-github-go@latest
   ```

2. Run the linter to check for deprecated patterns:
   ```bash
   wetwire-github lint .
   ```

3. Use `--fix` to automatically fix some issues:
   ```bash
   wetwire-github lint . --fix
   ```

4. Review CHANGELOG.md for breaking changes
