```
Bunnyshell CLI helps you manage environments in Bunnyshell and enable Remote Development.

Usage:
  bns [command]

Commands for Bunnyshell Resources:
  components         Components
  environments       Environments
  events             Events
  k8s-clusters       Kubernetes Cluster Integrations
  organizations      Organizations
  pipeline           Pipeline
  projects           Projects
  variables          Environment Variables

Commands for Utilities:
  git                Git Operations
  port-forward       Port Forward
  remote-development Remote Development

Commands for CLI:
  completion         Generate the autocompletion script for the specified shell
  configure          Configure CLI settings
  help               Help about any command
  version            Version Information

Flags:
      --configFile string   Bunnyshell CLI Config File (default "$HOME/.bunnyshell/config.yaml")
  -d, --debug               Debug network requests
  -h, --help                Help for bns
      --no-progress         Disable progress spinners
      --non-interactive     Disable interactive terminal
  -o, --output string       Output format: stylish | json | yaml (default "stylish")
      --profile string      Use profile from config file
  -v, --verbose count       Increase log verbosity
      --version             version for bns

Use "bns [command] --help" for more information about a command.
```

- [Installing](#installing)
  - [Generic Installer](#generic-installer)
  - [Homebrew](#homebrew)
  - [Downloading a Release from GitHub](#download-github-release)
  - [Docker Hub](#docker-hub)
- [Authentication](#authentication)
  - [Profiles](#profiles)
- [Shell Autocomplete](#shell-autocomplete)

## Installing

### Generic Installer
```sh
installer=$(curl --silent https://raw.githubusercontent.com/bunnyshell/cli/master/installer.sh) \
    && [ "33cdd268adb1e26511b0cc90c9f4bf017bb145041677ca52d23a0f93cd13bd58 *stdin" = "$(echo -n "${installer}" | openssl dgst -sha256 -r)" ] \
    && (
        sh -c "${installer}" || :
    ) \
    || echo "Checksum did not match $(echo -n "${installer}" | openssl dgst -sha256 -r)" \
;
```

| Env Var         | DEFAULT  | Description               |
| --------------- | -------- | ------------------------- |
| DEBUG_INSTALLER | false    | Add set -x for the script |
| PREFER_BREW     | false    | Allow brew installation   |
| INSTALL_PATH    | .        | Install path for binary   |
| SUDO_INSTALL    | false    | Use sudo when moving binary to INSTALL_PATH |

Examples:
```sh
installer=$(curl --silent https://raw.githubusercontent.com/bunnyshell/cli/master/installer.sh) \
    && [ "33cdd268adb1e26511b0cc90c9f4bf017bb145041677ca52d23a0f93cd13bd58 *stdin" = "$(echo -n "${installer}" | openssl dgst -sha256 -r)" ] \
    && (
        SUDO_INSTALL=true INSTALL_PATH=/usr/local/bin sh -c "${installer}" || :
    ) \
    || echo "Checksum did not match $(echo -n "${installer}" | openssl dgst -sha256 -r)" \
;
```

### Homebrew
```sh
brew install bunnyshell/tap/bunnyshell-cli
```

### Download Github Release

Download the appropriate archive for your architecture on the [releases page](https://github.com/bunnyshell/cli/releases)

And make it available in your `$PATH` or move the binary to `/usr/local/bin`

### Docker Hub
All the releases are found on: https://hub.docker.com/r/bunnyshell/cli

```sh
docker run --volume ~/.bunnyshell:/root/.bunnyshell bunnyshell/cli environments list
```

## Authentication
You will need an access token from https://environments.bunnyshell.com/access-token

You can then setup a profile for easy access to your acccount with:
```sh
bns configure profiles add
```

## Shell Autocomplete
Using `bns completion SHELL` you can generate autocomplete for your current shell.

### ZSH
```sh
echo 'source <(bns completion zsh)' >> ~/.zshrc
echo 'compdef _bns bns' >> ~/.zshrc
```

### Bash
```sh
echo 'source <(bns completion bash)' >> ~/.bashrc
```
