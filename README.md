# myher - Dependency Update Testing Tool

A command-line tool designed to test automated dependency update workflows using GitHub Actions, Renovate, and Dependabot.

## Purpose

`myher` is a testing tool that helps verify your automated dependency update pipeline works correctly. It does this by:

1. Intentionally downgrading dependencies in your go.mod file
2. Creating test scenarios for Renovate/Dependabot to detect
3. Triggering your automated update workflow to ensure it:
   - Detects outdated dependencies
   - Creates update branches
   - Runs tests
   - Merges successful updates to master

This allows you to verify your CI/CD pipeline handles dependency updates correctly before implementing it in production projects.

## Test Workflow

1. `myher` downgrades a dependency in go.mod
2. You commit and push the change
3. Renovate/Dependabot detects the outdated dependency
4. GitHub Actions workflow triggers:
   - Updates the dependency
   - Runs tests
   - Merges to master if tests pass

This creates a complete test cycle of your dependency update automation.

## Installation

```bash
go install github.com/gkwa/myher@latest
```

## Quick Start

### Parse go.mod Dependencies
Shows current module dependencies:
```bash
myher parse
```

### Generate Downgrade Commands
Generate commands to downgrade dependencies (creates test scenarios):
```bash
# Basic usage
myher downgrade

# With concurrent version checks
myher downgrade -c 10

# With alternating commented commands
myher downgrade --enable-alternating-comments
```

### Version Info
```bash
myher version
```

## Cheatsheet

```bash
# Basic Commands
myher parse                              # Show current dependencies
myher downgrade                          # Generate downgrade commands
myher version                            # Show version info

# Downgrade Options
myher downgrade -c 3                     # Run 3 concurrent version checks
myher downgrade --concurrent 10          # Run 10 concurrent version checks
myher downgrade --enable-alternating-comments  # Comment out every other command

# Global Options
myher --verbose                          # Increase output verbosity
myher --log-format json                  # Output logs in JSON format
myher --config /path/to/config.yaml      # Use custom config file
```

## Testing Your Update Workflow

1. Set up your repository with:
   - GitHub Actions workflow for testing and merging
   - Renovate configuration
   - Dependabot configuration (optional)

2. Use `myher` to create test scenarios:
   ```bash
   # Downgrade a dependency
   myher downgrade
   
   # Commit and push
   go mod tidy
   git commit -am "test: downgrade dependency for testing"
   git push
   ```

3. Watch your automation:
   - Renovate should detect the outdated dependency
   - GitHub Actions should run your tests
   - If tests pass, changes should merge to master

This allows you to verify your entire dependency update pipeline works as expected.
