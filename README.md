```sh
Bunnyshell CLI

Usage:
  bunnyshell-cli [command]

Available Commands:
  completion    Generate the autocompletion script for the specified shell
  components    Bunnyshell Components
  configure     Configure CLI settings
  environments  Bunnyshell Environments
  events        Bunnyshell Events
  help          Help about any command
  organizations Bunnyshell Organizations
  projects      Bunnyshell Projects
  version       Version Information

Flags:
  -c, --configFile string   Config file
  -d, --debug               Show network debug
      --feedback            Add feedback final output
  -h, --help                help for bunnyshell-cli
      --no-progress         Disable progress spinners
  -o, --output string       Output format: stylish | json | yaml (default "stylish")
  -p, --profile string      Force profile usage from config file
  -t, --timeout duration    Network timeout on requests (default 30s)
  -v, --verbose count       Number for the log level verbosity

Use "bunnyshell-cli [command] --help" for more information about a command.
```

- [Installing](#installing)
  - [Homebrew](#homebrew)
  - [Downloading a Release from GitHub](#download-github-release)
  - [Docker Hub](#docker-hub)
- [Authentication](#authentication)
  - [Profiles](#profiles)
- [Shell Autocomplete](#shell-autocomplete)

## Installing

### Homebrew
```sh
brew install bunnyshellosi/tap/cli
```

### Download Github Release

Download the appropriate archive for your architecture on the [releases page](https://github.com/bunnyshellosi/cli/releases)

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
# creates the config file
bunnyshell-cli configure init
# sets up the profile
bunnyshell-cli configure profiles add
```

## Shell Autocomplete
Using `bunnyshell-cli completion SHELL` you can generate autocomplete for your current shell.

### ZSH
```sh
echo 'source <(bunnyshell-cli completion zsh)' >> ~/.zshrc
echo 'compdef _bunnyshell-cli bunnyshell-cli' >> ~/.zshrc
```

### Bash
```sh
echo 'source <(bunnyshell-cli completion bash)' >> ~/.bashrc
```
