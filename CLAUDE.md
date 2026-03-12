# Bunnyshell CLI (bns) - Developer Guide

This document provides a comprehensive overview of the Bunnyshell CLI codebase to help developers and AI assistants understand the project structure, architecture, and development patterns.

## Project Overview

**Project Name:** Bunnyshell CLI (`bns`)
**Language:** Go 1.23
**Lines of Code:** ~17,000 (13,940 in `pkg/` + 3,093 in `cmd/`)
**Purpose:** Command-line tool for managing Bunnyshell environments, components, and remote development workflows.

## Development Environment

### Docker-based Development (Recommended)

The project uses a Docker-based development environment for consistent builds and testing.

**Container Setup:**
- Location: `.dev/` directory
- Container name: `bunnyshell-cli`
- Base image: `golang:1.23` with goreleaser pre-installed
- Working directory: `/usr/src/app` (mounted from project root)

**Starting the container:**
```bash
cd .dev
docker-compose up -d
```

**Accessing the container:**
```bash
docker exec -it bunnyshell-cli /bin/bash
```

**Building inside the container:**
```bash
# Inside the container
make build-local

# Successful build produces:
# - dist/bns_linux_amd64_v1/bns (Linux binary)
# - dist/bns_darwin_arm64/bns (macOS ARM binary)
# - dist/bns_darwin_amd64_v1/bns (macOS Intel binary)
```

**Testing the build:**
- From container: `./dist/bns_linux_amd64_v1/bns --help`
- From host (macOS): `./dist/bns_darwin_arm64/bns --help` or `./dist/bns_darwin_amd64_v1/bns --help`

**Notes:**
- Docker image building will fail in the dev container (no Docker-in-Docker) - this is expected and OK for development
- The build is considered successful if Linux and/or Darwin binaries are produced
- Host machine Go installation is NOT recommended - use the container for all Go commands

### Local Development with SDK Changes

If you need to test changes to the Bunnyshell SDK (`bunnyshell.com/dev`):

1. Add to `go.mod`:
   ```go
   replace bunnyshell.com/dev v0.7.0 => ../bunnyshellosi-dev/
   ```
2. Ensure the path works both in container and on host for IDE support
3. The path is already mounted in docker-compose.yaml

## Repository Structure

```
/cli
├── main.go                    # Entry point with panic recovery
├── go.mod / go.sum           # Go module dependencies
├── .dev/                      # Docker development environment
│   ├── docker-compose.yaml   # Container orchestration
│   ├── Dockerfile.dev        # Development container image
│   └── Readme.md             # Development quick start
├── cmd/                       # Command implementations (~3,093 LOC)
│   ├── root.go               # Root command setup with Cobra
│   ├── environment/          # Environment management commands
│   ├── component/            # Component operations
│   ├── pipeline/             # Pipeline monitoring
│   ├── project/              # Project management
│   ├── template/             # Template helpers
│   ├── variable/             # Variable management
│   ├── secret/               # Secret management
│   ├── configure/            # CLI configuration
│   ├── remote_development/   # Remote dev features
│   ├── git/                  # Git operations helper
│   ├── port-forward/         # Port forwarding
│   └── k8sIntegration/       # K8s integration
├── pkg/                       # Core packages (~13,940 LOC)
│   ├── api/                  # API client layers
│   ├── config/               # Configuration management
│   ├── formatter/            # Output formatting (stylish, JSON, YAML)
│   ├── interactive/          # Interactive prompts
│   ├── progress/             # Progress tracking for deployments
│   ├── remote_development/   # Remote dev implementation
│   ├── port_forward/         # Port forwarding logic
│   ├── k8s/                  # Kubernetes integration
│   ├── util/                 # Utility functions
│   ├── lib/                  # Helper libraries
│   ├── build/                # Build metadata
│   ├── wizard/               # Interactive wizards
│   └── net/                  # Network utilities
├── .github/workflows/        # CI/CD workflows
├── Dockerfile                # Container build
├── .goreleaser.yaml          # Release configuration
├── Makefile                  # Build targets
└── README.md / LICENSE       # Documentation
```

## Command Structure

The CLI follows a hierarchical command structure using Cobra:

```
bns [global flags] [command] [subcommand] [flags]
```

### Top-Level Command Groups

1. **Resource Commands** (API-backed operations)
   - `environments` (env) - Environment CRUD and deployment
   - `components` (comp) - Component management
   - `pipeline` - Pipeline monitoring
   - `projects` - Project management
   - `variables` - Environment variables
   - `secrets` - Secret management
   - `templates` - Template utilities
   - `k8s-clusters` - Kubernetes integrations
   - And more...

2. **Utility Commands**
   - `git` - Git operations helper
   - `port-forward` - Port forwarding management
   - `remote-development` - Remote development workspace management

3. **CLI Commands**
   - `configure` - Profile and settings configuration
   - `completion` - Shell autocompletion
   - `version` - Version information

### Global Flags

- `--configFile` - Config file path (default: `$HOME/.bunnyshell/config.yaml`)
- `-d, --debug` - Debug network requests
- `--no-progress` - Disable progress spinners
- `--non-interactive` - Non-interactive mode
- `-o, --output` - Output format (stylish, json, yaml)
- `--profile` - Profile selection
- `-v, --verbose` - Verbosity level

## Key Packages

### `pkg/api/` - API Client Layer

Wraps the Bunnyshell SDK with domain-specific logic. Organized by resource type:
- `environment/` - Environment operations (deploy, create, delete, etc.)
- `component/` - Component management
- `pipeline/` - Pipeline queries and monitoring
- `project/` - Project operations
- `variable/` - Variable CRUD
- `template/` - Template handling
- `event/` - Event fetching
- `k8s/` - Kubernetes integrations
- `common/` - Shared request/response types

### `pkg/config/` - Configuration Management

**Key Files:**
- `manager.go` - Central config lifecycle management
- `options.go` - CLI flag/config option handling
- `manager.cobra.go` - Cobra command integration
- `manager.loader.go` - Config file parsing (YAML)

**Core Types:**
- `Profile` - Named credentials (token, host, scheme)
- `Context` - Default resource selection (org, project, env, component)
- `Config` - Full configuration structure

**Features:** Profile switching, context persistence, flag binding

### `pkg/formatter/` - Output Formatting

Multiple output formats with specialized renderers:
- **Stylish** - Human-readable colored terminal output (default)
- **JSON** - Machine-readable output
- **YAML** - Structured data format

Specialized formatters for different resources (environments, pipelines, templates, etc.)

### `pkg/interactive/` - Interactive Prompts

Built on `github.com/AlecAivazis/survey/v2`:
- Input prompts, selections, confirmations
- Password inputs with validation
- Path suggestions
- `AskMissingRequiredFlags()` - Auto-prompt for missing CLI flags

### `pkg/progress/` - Deployment Progress Tracking

- **Event Monitoring** - Watches deployment events in real-time
- **Pipeline Tracking** - Monitors pipeline execution with status updates
- **Spinner UI** - Visual feedback during long operations
- **Retry Logic** - Built-in retry mechanisms

### `pkg/remote_development/` - Remote Development

Complete remote development workspace management:
- **Workspace** - Manages local dev environments
- **Config** - Remote dev configuration
- **Action** - Start/stop remote dev
- Integrates with K8s port forwarding

### `pkg/port_forward/` - Port Forwarding

Kubernetes port forwarding implementation:
- **Port Mapping** - Local-to-remote port mapping (format: `local:remote`)
- **Workspace** - Persistent forwarding state
- **Manager** - Signal handling and lifecycle

### `pkg/k8s/` - Kubernetes Integration

- Direct K8s client integration (k8s.io/client-go)
- Pod management
- Service discovery
- Cluster operations
- **kubectl/exec** - Wrapper around kubectl exec functionality for executing commands in containers
- **wizard/k8s** - Interactive pod and container selection wizards

### `pkg/util/` - Utilities

Common helpers:
- `cobra.go` - Cobra command helpers
- `network.go` - Port availability checks
- `workspace.go` - `.bunnyshell/` directory management
- `spinner.go` - Consistent spinner creation
- `os.go` - File/OS operations

### `pkg/build/` - Build Metadata

```go
var (
  Name = "bns"
  EnvPrefix = "bunnyshell"
  Version = "dev" // Set via ldflags during build
  Commit, Date = "none", "unknown"
)
```

## Technology Stack

### Core Framework
- **CLI Framework:** Spf13 Cobra (v1.8.1)
- **Config Management:** Spf13 Viper (v1.19.0)
- **Flag Parsing:** Spf13 pflag (v1.0.5)

### External Services
- **Bunnyshell SDK** (bunnyshell.com/sdk v0.20.4) - Official API client
- **Bunnyshell Dev** (bunnyshell.com/dev v0.7.2) - Development utilities

### Kubernetes Integration
- k8s.io/client-go (v0.30.2)
- k8s.io/cli-runtime (v0.30.2)
- k8s.io/kubectl (v0.30.2)

### User Interaction
- **Survey** (AlecAivazis/survey v2.3.7) - Interactive prompts
- **Color** (fatih/color v1.17.0) - Terminal colors
- **Spinner** (briandowns/spinner v1.23.1) - Progress spinners

### Utilities
- **Retry Logic** (avast/retry-go v4.6.0)
- **YAML** (gopkg.in/yaml.v3 v3.0.1)
- **Enum Flags** (thediveo/enumflag v2.0.5)

### Build & Deployment
- GoReleaser (v1.x) - Multi-platform binary releases
- Docker - Container builds (Alpine base)

## Architecture Patterns

### Command Pattern (Cobra)

```
cmd/
├── [resource]/
│   ├── root.go        # Main command with API setup
│   ├── list.go        # List subcommand
│   ├── show.go        # Detail view subcommand
│   └── action/        # Nested actions (create, delete, deploy, etc.)
```

Each resource follows a consistent structure:
- Resource main command registers API client
- Subcommands (list, show, actions)
- Actions organized in `action/` subdirectory

### Layered Architecture

```
cmd/ (Command handlers)
  ↓
pkg/api/ (Business logic)
  ↓
bunnyshell.com/sdk (API client)
  ↓
Bunnyshell Backend API
```

### Configuration Injection

- `config.MainManager` - Singleton manager
- `CommandWithAPI()` - Auto-injects API client
- `CommandWithGlobalOptions()` - Adds standard flags
- Flag defaults bound to config file values

### Error Handling

- Domain-specific error types (ErrInvalidValue, ErrUnknownProfile)
- Panic recovery at main() level
- Command error formatting with `lib.FormatCommandError()`

### Option Pattern

```go
type DeployOptions struct {
  ID string
  WithoutPipeline bool
  Interval time.Duration
  // ...
}
```

### Fluent/Builder Patterns

- `NewInput()` with method chaining for prompts
- `NewOptions()` factories for API calls

## Main Features

1. **Environment Management**
   - Create, read, update, delete environments
   - Deploy with K8s integration
   - Start/stop/pause/resume environments
   - View endpoints
   - Clone environments

2. **Component Operations**
   - List/show components
   - Component-specific variables
   - Component debugging
   - SSH access to containers

**Note:** Command execution (`exec`) is now a top-level utility command. See "Utility Commands" section.

3. **Deployment Pipeline**
   - Monitor deployment pipelines
   - Track events in real-time
   - Queue deployments
   - Configure pipeline intervals

4. **Remote Development**
   - Start remote dev workspace
   - Stop remote dev
   - Full local-remote synchronization

5. **Utility Commands**
   - **exec** - Execute arbitrary commands in component containers
   - **git** - Git operations
   - **port-forward** - Port forwarding management
   - **remote-development** - Remote development workspace
   - **debug** - Component debugging
   - **ssh** - SSH into running containers

6. **Port Forwarding**
   - Local-to-remote port mapping
   - K8s pod port forwarding
   - Persistent workspace management

7. **Configuration & Profiles**
   - Multi-profile support
   - Context persistence (org, project, env, component)
   - Profile switching and default setting
   - Interactive configuration setup

8. **Resource Management**
   - Variables (project and environment level)
   - Secrets management
   - Templates
   - Registries
   - K8s cluster integrations

9. **Git Integration**
   - Git repository preparation
   - Git operation utilities

10. **Output Formatting**
   - Stylish (terminal-optimized)
   - JSON (programmatic)
   - YAML (structured)

## Configuration Files

**Default Config Template:** `config.sample.yaml`
```yaml
defaultprofile: sample
outputformat: json
profiles:
  sample:
    token: ''
```

**User Config:** `$HOME/.bunnyshell/config.yaml`

## Build System

### Makefile

```makefile
build-local:
  goreleaser release --snapshot --rm-dist
```

### GoReleaser Configuration (`.goreleaser.yaml`)

- **Multi-platform builds:** Linux, Windows, Darwin (macOS)
- **Architecture support:** amd64, arm64, 386
- **Build flags:** CGO_ENABLED=0, -trimpath
- **Version injection:** Via ldflags (Version, Commit, Date)
- **Docker builds:** Alpine-based images for amd64 and arm64
- **Homebrew distribution:** Via homebrew-tap
- **Artifacts:** Compressed archives (tar.gz, zip), checksums

### CI/CD Pipelines

- **audit.yaml** - Code quality checks
- **release.yaml** - Automated releases

### Docker Image

- Base: Alpine Linux
- Includes: jq, bash, curl, sed
- Sets up config directory at `/root/.bunnyshell`
- Entrypoint: `bns` command

## Development Workflow

### Example: Deploy Environment

```
$ bns environments deploy my-env
  ↓
cmd/environment/action/deploy.go
  ↓
Validates K8s integration via pkg/api/environment
  ↓
Calls SDK Deploy API
  ↓
Monitors pipeline via pkg/progress/event.go
  ↓
Displays spinner with real-time status updates
  ↓
Shows endpoints via pkg/api/component/endpoint
```

### Example: Execute Command in Component

```
$ bns exec comp-123 -- ls -la
  ↓
cmd/exec/root.go
  ↓
Parses component ID from positional arg
  ↓
Fetches component details via pkg/api/component
  ↓
Retrieves kubeconfig via pkg/api/environment
  ↓
Interactive pod selection via pkg/wizard/k8s
  ↓
Interactive container selection via pkg/wizard/k8s
  ↓
Creates exec command via pkg/k8s/kubectl/exec
  ↓
Executes command in container using k8s.io/kubectl/pkg/cmd/exec
```

### Config Binding Flow

```go
// Flags automatically:
// 1. Read from config file
// 2. Overridden by CLI flags
// 3. Prompted interactively if missing
config.MainManager.CommandWithGlobalOptions(cmd)
```

### Interactive Flows

```go
// Gracefully asks for missing required values
interactive.AskMissingRequiredFlags(cmd)
```

## Adding New Commands

To add a new command or resource:

1. Create directory under `cmd/[resource]/`
2. Create `root.go` with main command definition
3. Add subcommands (list, show, etc.)
4. Create action handlers in `action/` subdirectory
5. Implement API layer in `pkg/api/[resource]/`
6. Add formatters in `pkg/formatter/` if needed
7. Register command in `cmd/root.go`

## Testing

The codebase includes:
- Unit tests throughout packages
- CI/CD automation via GitHub Actions
- Code quality checks (audit workflow)

## Utility Commands

The CLI provides several top-level utility commands for working with components:

### `bns exec`
**Location:** `cmd/exec/root.go`

Execute arbitrary commands in component containers, similar to `kubectl exec` or `docker exec`.

**Command Structure:**
- Component ID as **positional argument**: `bns exec <component-id>`
- Falls back to context if ID not provided
- Top-level utility command (not under `components`)

**Features:**
- Interactive or explicit pod/container selection
- Support for TTY and stdin allocation via `--tty` and `--stdin` flags
- Defaults to interactive shell (`/bin/sh`) when no command specified
- Auto-enables `--tty` and `--stdin` for interactive shells
- Supports `namespace/pod-name` format for pod specification
- Uses existing Kubernetes exec infrastructure (`pkg/k8s/kubectl/exec`)

**Usage:**
```bash
# Interactive shell (auto-enables --tty and --stdin)
bns exec comp-123

# Run specific command
bns exec comp-123 -- ls -la /app

# Explicit pod and container
bns exec comp-123 --tty --stdin --pod my-pod -c api -- /bin/bash

# Use component from context (no ID needed)
bns configure set-context --component comp-123
bns exec --tty --stdin

# Pipe local script to remote container
bns exec comp-123 --stdin -- python3 < local-script.py
```

**Key Components:**
- `cmd/exec/root.go` - Main exec command implementation
- `pkg/k8s/kubectl/exec/exec.go` - Wraps kubectl exec functionality
- `pkg/wizard/k8s/pod.go` - Interactive pod selection
- `pkg/wizard/k8s/container.go` - Interactive container selection
- `pkg/api/environment/action_kubeconfig.go` - Kubeconfig retrieval

**Design Notes:**
- No shorthand flags for `--tty` and `--stdin` to avoid conflicts with global `-t` (timeout) flag
- Follows docker/kubectl conventions for TTY and stdin handling

## Component Actions

The CLI provides component-level action commands under the `bns components` namespace:

### `bns components ssh`
**Location:** `cmd/component/action/ssh.go`

Provides interactive shell access with environment banner (MOTD).

**Differences from `bns exec`:**
- Always allocates TTY and enables stdin
- Shows environment information banner
- Fixed shell command with MOTD display
- Designed specifically for interactive sessions
- Requires `--id` flag for component ID

### `bns components port-forward`
**Location:** `cmd/component/action/port_forward.go`

Forward local ports to component containers for development and debugging.

## Common Code Locations

- **Entry point:** `main.go:23`
- **Root command:** `cmd/root.go:29`
- **Config manager:** `pkg/config/manager.go`
- **API clients:** `pkg/api/[resource]/`
- **Interactive prompts:** `pkg/interactive/`
- **Progress tracking:** `pkg/progress/`
- **Output formatting:** `pkg/formatter/`
- **Component actions:** `cmd/component/action/`
- **Utility commands:** `cmd/exec/`, `cmd/git/`, `cmd/utils/`
- **K8s exec wrapper:** `pkg/k8s/kubectl/exec/exec.go`
- **Pod/Container wizards:** `pkg/wizard/k8s/`

## Environment Variables

All CLI-specific environment variables use the `bunnyshell_` prefix (from `build.EnvPrefix`).

## Key Dependencies

```
bunnyshell.com/sdk (v0.20.4)          # HTTP API client
bunnyshell.com/dev (v0.7.2)           # Development utilities
github.com/spf13/cobra (v1.8.1)       # Command framework
github.com/spf13/viper (v1.19.0)      # Configuration
k8s.io/* (v0.30.2)                    # Kubernetes integration
```

## File Statistics

- **Total Go Files:** 206
- **Command LOC:** 3,093
- **Package LOC:** 13,940
- **Total LOC:** ~17,000

## Important Notes

### Flag Conflicts
When adding new flags to commands, be aware of global flags that are already registered:
- `-t` is used by `--timeout` (global flag in `pkg/config/options.go:90`)
- `-d` is used by `--debug` (global flag)
- `-v` is used by `--verbose` (global flag)
- `-o` is used by `--output` (global flag)

Command-specific flags should avoid these shorthands. For example, the `exec` command uses `--tty` and `--stdin` without shorthands to avoid conflict with the global `-t` timeout flag.

### Pod/Container Selection Pattern
Commands that interact with Kubernetes pods follow a consistent pattern:
1. Get component from context or flag
2. Retrieve kubeconfig for the environment
3. Create K8s client
4. Interactive pod selection (if not specified via `--pod` flag)
5. Interactive container selection (if not specified via `--container` flag)
6. Execute the operation

See `cmd/component/action/exec.go` and `cmd/component/action/ssh.go` for reference implementations.

## Notes for AI Assistants

### Development Workflow for AI Assistants

**IMPORTANT: Always use the Docker container for Go commands**

When you need to run Go commands (build, test, mod tidy, etc.):

1. **Check if the container is running:**
   ```bash
   docker ps --filter "name=bunnyshell-cli"
   ```

2. **Start container if not running:**
   ```bash
   cd .dev && docker-compose up -d
   ```

3. **Execute Go commands inside the container:**
   ```bash
   docker exec -it bunnyshell-cli <command>
   # Example: docker exec -it bunnyshell-cli make build-local
   ```

4. **Build success criteria:**
   - Build is considered successful if Linux and/or Darwin binaries are produced
   - Docker image building will fail (no Docker-in-Docker) - this is EXPECTED and OK
   - Look for: `dist/bns_linux_amd64_v1/bns` and/or `dist/bns_darwin_arm64/bns`

5. **DO NOT rely on host machine Go:**
   - Host may not have Go installed
   - Host Go version may differ
   - Container provides consistent environment

### General Development Patterns

- The codebase follows clear separation of concerns: CLI commands, business logic, and utilities
- Commands are structured hierarchically using Cobra
- API layer wraps SDK calls with domain-specific logic
- Interactive mode gracefully handles missing parameters
- Progress tracking provides real-time feedback for long operations
- Configuration supports multiple profiles and contexts
- Build system supports cross-platform releases via GoReleaser
- When adding new utility commands, follow the pattern in `cmd/exec/root.go` or `cmd/git/root.go`
- When adding new component actions, follow the pattern in `cmd/component/action/ssh.go`
- Always check for flag conflicts with global flags before adding shorthand flags
- Positional arguments are preferred for primary identifiers (e.g., `bns exec <component-id>` instead of `--id` flag)

### AI-Agnostic Documentation

This project maintains AI-agnostic documentation for use across different AI assistants:
- **AGENTS.md** - General instructions for all AI agents (to be created/maintained)
- **CLAUDE.md** - This file (Claude-specific but should contain general knowledge)

When updating documentation, consider whether the knowledge should be in AI-agnostic format in AGENTS.md.
