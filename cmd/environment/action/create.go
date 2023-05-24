package action

import (
	"errors"
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var (
	errCreateSourceNotProvided   = errors.New("template id, content or git repository must be provided")
	errK8SIntegrationNotProvided = errors.New("kubernetes integration must be provided when deploying")
)

type CreateSource struct {
	TemplateID string

	Git string

	YamlPath string

	GitRepo   string
	GitBranch string
	GitPath   string
}

func (createSource *CreateSource) UpdateCommandFlags(command *cobra.Command) {
	flags := command.Flags()

	flags.StringVar(&createSource.Git, "from-git", createSource.Git, "Use a template git repository during creation")
	flags.StringVar(&createSource.TemplateID, "from-template", createSource.TemplateID, "Use a TemplateID during creation")
	flags.StringVar(&createSource.YamlPath, "from-path", createSource.YamlPath, "Use a local bunnyshell.yaml during creation")
	flags.StringVar(&createSource.GitRepo, "from-git-repo", createSource.GitRepo, "Git repository for the environment")
	flags.StringVar(&createSource.GitBranch, "from-git-branch", createSource.GitBranch, "Git branch for the environment")
	flags.StringVar(&createSource.GitPath, "from-git-path", createSource.GitPath, "Git path for the environment")

	command.MarkFlagsMutuallyExclusive("from-git", "from-template", "from-path", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-branch", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-path", "from-git-repo")

	_ = command.MarkFlagFilename("from-path", "yaml", "yml")
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := environment.NewCreateOptions()
	createSource := CreateSource{}

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if createSource.Git == "" && createSource.TemplateID == "" && createSource.YamlPath == "" && createSource.GitRepo == "" {
				return errCreateSourceNotProvided
			}

			if createOptions.WithDeploy && createOptions.GetKubernetesIntegration() == "" {
				if !settings.IsStylish() {
					return errK8SIntegrationNotProvided
				}
			}

			return validateActionOptions(&createOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.Project = settings.Profile.Context.Project

			if err := parseCreateOptions(createSource, createOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := environment.Create(createOptions)
			if err != nil {
				var apiError api.Error

				if errors.As(err, &apiError) {
					return handleCreateErrors(cmd, apiError, createOptions)
				}

				return lib.FormatCommandError(cmd, err)
			}

			if !createOptions.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &createOptions.DeployOptions
			deployOptions.ID = model.GetId()

			return handleDeploy(cmd, deployOptions, "created", createOptions.GetKubernetesIntegration())
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.AddFlagWithExtraHelp(
		"project",
		"Project for the environment",
		"Projects contain environments along with build settings and project variables",
		util.FlagRequired,
	))

	createOptions.UpdateFlagSet(flags)

	createSource.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}

func handleCreateErrors(cmd *cobra.Command, apiError api.Error, createOptions *environment.CreateOptions) error {
	genesisName := getCreateGenesisName(createOptions)

	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		cmd.Printf("Problem with %s: %s\n", genesisName, violation.GetMessage())
	}

	return lib.ErrGeneric
}

func getCreateGenesisName(createOptions *environment.CreateOptions) string {
	if createOptions.Genesis.FromGitSpec != nil {
		return "--from-git"
	}

	if createOptions.Genesis.FromTemplate != nil {
		return "--from-template"
	}

	if createOptions.Genesis.FromString != nil {
		return "--from-path"
	}

	return "arguments"
}

func parseCreateOptions(createSource CreateSource, createOptions *environment.CreateOptions) error {
	createOptions.Genesis = &sdk.EnvironmentCreateActionGenesis{}

	if createSource.Git != "" {
		fromGitSpec := sdk.NewFromGitSpec()
		fromGitSpec.Spec = &createSource.Git

		createOptions.Genesis.FromGitSpec = fromGitSpec

		return nil
	}

	if createSource.TemplateID != "" {
		fromTemplate := sdk.NewFromTemplate()
		fromTemplate.Template = &createSource.TemplateID

		createOptions.Genesis.FromTemplate = fromTemplate

		return nil
	}

	if createSource.YamlPath != "" {
		fromString := sdk.NewFromString()

		bytes, err := readFile(createSource.YamlPath)
		if err != nil {
			return err
		}

		content := string(bytes)
		fromString.Yaml = &content

		createOptions.Genesis.FromString = fromString

		return nil
	}

	if createSource.GitRepo != "" {
		fromGit := sdk.NewFromGit()
		fromGit.Url = &createSource.GitRepo
		fromGit.Branch = &createSource.GitBranch
		fromGit.YamlPath = &createSource.GitPath

		createOptions.Genesis.FromGit = fromGit

		return nil
	}

	return errCreateSourceNotProvided
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
