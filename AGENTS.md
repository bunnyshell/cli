# AI Agent Development Guide

This document provides instructions for AI assistants working on the Bunnyshell CLI codebase. It is maintained in an AI-agnostic format for use with any AI development assistant.

## Quick Start

**Project:** Bunnyshell CLI (`bns`)
**Language:** Go 1.23
**Development Environment:** Docker-based (required)

## Development Environment Setup

### Using Docker (Required for Go Commands)

All Go-related commands MUST be executed inside the Docker development container.

**1. Check if container is running:**
```bash
docker ps --filter "name=bunnyshell-cli"
```

**2. Start container if needed:**
```bash
cd .dev && docker-compose up -d
```

**3. Execute commands in container:**
```bash
docker exec -it bunnyshell-cli <command>
```

**Common commands:**
```bash
# Build the project
docker exec -it bunnyshell-cli make build-local

# Run go mod tidy
docker exec -it bunnyshell-cli go mod tidy

# Run tests
docker exec -it bunnyshell-cli go test ./...

# Access container shell
docker exec -it bunnyshell-cli /bin/bash
```

### Build Success Criteria

A build is considered **successful** when:
- ✅ Linux binary is produced: `dist/bns_linux_amd64_v1/bns`
- ✅ Darwin binary is produced: `dist/bns_darwin_arm64/bns` or `dist/bns_darwin_amd64_v1/bns`

A build may show errors for:
- ❌ Docker image building (no Docker-in-Docker in dev container) - **THIS IS EXPECTED AND OK**

### Testing Builds

**From container:**
```bash
./dist/bns_linux_amd64_v1/bns --help
```

**From host (macOS):**
```bash
./dist/bns_darwin_arm64/bns --help
# or
./dist/bns_darwin_amd64_v1/bns --help
```

### Important Notes

- **DO NOT** rely on host machine Go installation
- **DO NOT** run `go` commands directly on the host
- **ALWAYS** use the Docker container for Go commands
- The host may not have Go installed or may have a different version

## Project Structure

```
/
├── .dev/                      # Docker development environment
│   ├── docker-compose.yaml   # Container setup
│   ├── Dockerfile.dev        # golang:1.23 with goreleaser
│   └── Readme.md             # Quick reference
├── cmd/                       # Command implementations
│   └── [resource]/           # Command groups (environments, components, etc.)
│       ├── root.go           # Main command
│       ├── list.go           # List subcommand
│       ├── show.go           # Show subcommand
│       └── action/           # Action subcommands (create, delete, etc.)
├── pkg/                       # Core packages
│   ├── api/                  # API client wrappers
│   ├── config/               # Configuration management
│   ├── formatter/            # Output formatters (stylish, JSON, YAML)
│   ├── interactive/          # Interactive prompts
│   └── ...                   # Other core packages
├── main.go                    # Application entry point
├── go.mod / go.sum           # Go dependencies
├── Makefile                   # Build targets
├── .goreleaser.yaml          # Release configuration
└── AGENTS.md                  # This file
```

## Adding a New Command

Follow this pattern when adding new commands:

### 1. Create API Layer

Location: `pkg/api/[resource]/`

```go
// pkg/api/[resource]/list.go
package resource

import (
    "bunnyshell.com/cli/pkg/api"
    "bunnyshell.com/cli/pkg/api/common"
    "bunnyshell.com/cli/pkg/lib"
    "bunnyshell.com/sdk"
)

type ListOptions struct {
    common.ListOptions
    // Add your filters here
}

func NewListOptions() *ListOptions {
    return &ListOptions{
        ListOptions: *common.NewListOptions(),
    }
}

func List(options *ListOptions) (*sdk.PaginatedResourceCollection, error) {
    model, resp, err := ListRaw(options)
    if err != nil {
        return nil, api.ParseError(resp, err)
    }
    return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedResourceCollection, *http.Response, error) {
    profile := options.GetProfile()
    ctx, cancel := lib.GetContextFromProfile(profile)
    defer cancel()

    request := lib.GetAPIFromProfile(profile).ResourceAPI.ResourceList(ctx)
    return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiResourceListRequest, options *ListOptions) sdk.ApiResourceListRequest {
    if options == nil {
        return request
    }

    if options.Page > 1 {
        request = request.Page(options.Page)
    }

    // Add your filters here

    return request
}
```

### 2. Create Command Implementation

Location: `cmd/[resource]/`

```go
// cmd/[resource]/list.go
package resource

import (
    "bunnyshell.com/cli/pkg/api/resource"
    "bunnyshell.com/cli/pkg/lib"
    "github.com/spf13/cobra"
)

func init() {
    listOptions := resource.NewListOptions()

    command := &cobra.Command{
        Use: "list",
        Short: "List resources",
        ValidArgsFunction: cobra.NoFileCompletions,

        RunE: func(cmd *cobra.Command, args []string) error {
            return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
                return resource.List(listOptions)
            })
        },
    }

    flags := command.Flags()
    // Add your flags here
    listOptions.UpdateFlagSet(flags)

    mainCmd.AddCommand(command)
}
```

### 3. Add Formatter (if needed)

Location: `pkg/formatter/`

```go
// pkg/formatter/stylish.resource.go
package formatter

import (
    "fmt"
    "text/tabwriter"
    "bunnyshell.com/sdk"
)

func tabulateResourceCollection(writer *tabwriter.Writer, data *sdk.PaginatedResourceCollection) {
    fmt.Fprintf(writer, "%v\t %v\t %v\n", "ID", "Name", "Status")

    if data.Embedded != nil {
        for _, item := range data.Embedded.Item {
            fmt.Fprintf(writer, "%v\t %v\t %v\n",
                item.GetId(),
                item.GetName(),
                item.GetStatus(),
            )
        }
    }
}
```

Then add the case to `pkg/formatter/stylish.go`:

```go
case *sdk.PaginatedResourceCollection:
    tabulateResourceCollection(writer, dataType)
```

## Common Patterns

### Flag Conflicts to Avoid

These global flags are already registered:
- `-t` = `--timeout`
- `-d` = `--debug`
- `-v` = `--verbose`
- `-o` = `--output`

Do not use these shorthands for command-specific flags.

### Repeatable Flags

Use `StringArrayVar` for repeatable flags:

```go
var statuses []string
flags.StringArrayVar(&statuses, "status", statuses, "Filter by status (repeatable)")
```

### Required Flags

```go
flags.AddFlag(option.GetRequiredFlag("id"))
```

### Optional Context-aware Flags

```go
flags.AddFlag(options.Organization.GetFlag("organization"))
```

## Testing Your Changes

### 1. Build the project

```bash
docker exec -it bunnyshell-cli make build-local
```

### 2. Verify build succeeded

Check for:
- `dist/bns_linux_amd64_v1/bns` (Linux)
- `dist/bns_darwin_arm64/bns` (macOS ARM)
- `dist/bns_darwin_amd64_v1/bns` (macOS Intel)

### 3. Test the binary

From host (macOS):
```bash
./dist/bns_darwin_arm64/bns [your-command] --help
```

From container (Linux):
```bash
./dist/bns_linux_amd64_v1/bns [your-command] --help
```

## SDK Dependencies

The project depends on:
- `bunnyshell.com/sdk` - Official Bunnyshell API client
- `bunnyshell.com/dev` - Development utilities

If you need to test with local SDK changes:

1. Add to `go.mod`:
   ```go
   replace bunnyshell.com/dev v0.7.0 => ../bunnyshellosi-dev/
   ```

2. Ensure the path exists in both container and host (for IDE support)

3. The docker-compose.yaml already mounts this path

## Code Quality

### Before committing:

```bash
# Format code
docker exec -it bunnyshell-cli go fmt ./...

# Tidy dependencies
docker exec -it bunnyshell-ci go mod tidy

# Run tests
docker exec -it bunnyshell-cli go test ./...

# Build to verify
docker exec -it bunnyshell-cli make build-local
```

## Architecture Principles

1. **Separation of Concerns:**
   - `cmd/` = CLI interface and command handling
   - `pkg/api/` = Business logic and API interaction
   - `pkg/formatter/` = Output formatting
   - `pkg/lib/` = Shared utilities

2. **Configuration Management:**
   - Support multiple profiles
   - Store context (org, project, env, component)
   - Allow flag overrides
   - Interactive prompts for missing values

3. **User Experience:**
   - Provide interactive mode for missing parameters
   - Support non-interactive mode for automation
   - Multiple output formats (stylish, JSON, YAML)
   - Progress indicators for long operations

4. **Error Handling:**
   - Use domain-specific error types
   - Parse and format API errors
   - Provide helpful error messages

## Documentation

When adding features, update:
- **AGENTS.md** (this file) - For AI-agnostic instructions
- **CLAUDE.md** - For detailed codebase overview
- **README.md** - For user-facing documentation
- `.dev/Readme.md` - For development quick reference

## Getting Help

- Check **CLAUDE.md** for comprehensive codebase documentation
- Look at existing commands for patterns (e.g., `cmd/pipeline/jobs.go`)
- Examine API packages for SDK usage (e.g., `pkg/api/workflow_job/list.go`)
- Review formatters for output patterns (e.g., `pkg/formatter/stylish.workflow_job.go`)
