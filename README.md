# myher - Go Module Dependency Downgrade Tool

myher helps you validate your dependency update automation by deliberately downgrading dependencies to test update detection.

## Install

```bash
go install github.com/gkwa/myher@latest
```

## Usage

### Parse Dependencies
```bash
myher parse

# Module: github.com/gkwa/myher
# Go Version: 1.21
# 
# Direct Dependencies:
# github.com/pkg/errors v0.9.1
# github.com/stretchr/testify v1.8.4
```

### Generate Downgrade Commands
```bash
myher downgrade

# Output:
# go get github.com/pkg/errors@v0.9.0  # Downgrades to previous version
# go get github.com/stretchr/testify@v1.8.3
```

Options:
```bash
myher downgrade -c 10                         # Run 10 concurrent version checks
myher downgrade --enable-alternating-comments # Comment every other command
```

## Workflow

1. Generate downgrade commands:
```bash
myher downgrade > downgrade.sh
```

2. Apply downgrades:
```bash 
chmod +x downgrade.sh
./downgrade.sh
go mod tidy
```

3. Commit & push to trigger your update automation:
```bash
git commit -am "test: downgrade dependencies"
git push
```

4. Verify your automation:
- Detects outdated dependencies
- Creates update PRs
- Runs tests
- Merges successful updates

## Downgrade Strategy

For each dependency, myher:
1. Gets the current version from go.mod
2. Fetches all available versions 
3. Finds the version directly before current in the version list
4. Generates command to downgrade to that version

Example:
- Current: v1.8.4
- Available: [v1.7.0, v1.8.0, v1.8.3, v1.8.4] 
- Downgrades to: v1.8.3

## Command Reference

```bash
# Show current dependencies
myher parse

# Generate downgrade commands
myher downgrade

# Options
-c, --concurrent N             # N concurrent version checks
--enable-alternating-comments # Comment every other command 

# Global flags
--verbose      # Increase output detail
--log-format   # text/json logging
```