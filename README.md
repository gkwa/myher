# myher - Go Module Dependency Downgrade Tool

A command-line tool that helps test dependency update automation by generating commands to downgrade Go module dependencies to their second-latest versions.

## How It Works

`myher` performs these steps:

1. Parses your `go.mod` file to extract import paths
2. For each import, runs `go list -m versions {import path}` to get all available versions
3. Identifies the latest version and the version immediately before it
4. Generates `go get` commands to downgrade to the second-latest version

This creates test scenarios for dependency update automation tools (Renovate, Dependabot) to detect and update.

## Purpose

This tool helps test your dependency update automation pipeline:

1. Use `myher` to generate downgrade commands
2. Apply the downgrades to create outdated dependencies
3. Push changes to trigger your update automation
4. Verify that your CI/CD pipeline:
   - Detects the outdated dependencies
   - Creates update branches
   - Runs tests
   - Merges successful updates

## Installation

```bash
go install github.com/gkwa/myher@latest
```

## Usage

### Parse go.mod
Shows the current module dependencies:
```bash
myher parse
```

Example output:
```
github.com/pkg/errors v0.9.1
github.com/stretchr/testify v1.8.4
```

### Generate Downgrade Commands
Generates commands to downgrade dependencies to their second-latest versions:
```bash
myher downgrade
```

Example output:
```
go get github.com/pkg/errors@v0.9.0
go get github.com/stretchr/testify@v1.8.3
```

Options:
```bash
# Run multiple version checks concurrently
myher downgrade -c 10

# Alternate between commented and uncommented commands
myher downgrade --enable-alternating-comments
```

### Version Info
```bash
myher version
```

## Testing Your Update Workflow

1. Generate downgrade commands:
```bash
myher downgrade > downgrade.sh
```

2. Apply the downgrades:
```bash
chmod +x downgrade.sh
./downgrade.sh
go mod tidy
```

3. Commit and push to trigger automation:
```bash
git commit -am "test: downgrade dependencies for testing"
git push
```

4. Watch your automation workflow:
   - Renovate/Dependabot should detect outdated dependencies
   - GitHub Actions should run tests
   - Changes should merge if tests pass

## Command Options

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

## Integration Notes

To complete the testing pipeline, you'll need:

- GitHub Actions workflow for running tests
- Renovate or Dependabot configuration
- Branch protection rules (optional)

The tool works best as part of an automated dependency update workflow but can also be used standalone to generate downgrade commands.
