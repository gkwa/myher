# myher - Go Module Helper

A command-line tool to help manage Go module dependencies.

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
Generate commands to downgrade dependencies to their second-latest versions:
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