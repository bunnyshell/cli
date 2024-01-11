package action

import (
	"errors"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

type EditSource struct {
	TemplateID string

	Git string

	YamlPath string

	GitRepo   string
	GitBranch string
	GitPath   string
}

func (es *EditSource) UpdateCommandFlags(command *cobra.Command) {
	flags := command.Flags()

	flags.StringVar(&es.Git, "from-git", es.Git, "Use a template git repository during update")
	flags.StringVar(&es.TemplateID, "from-template", es.TemplateID, "Use a template ID during update")
	flags.StringVar(&es.YamlPath, "from-path", es.YamlPath, "Use a local environment yaml during update")
	flags.StringVar(&es.GitRepo, "from-git-repo", es.GitRepo, "Git repository for the environment template")
	flags.StringVar(&es.GitBranch, "from-git-branch", es.GitBranch, "Git branch for the environment template")
	flags.StringVar(&es.GitPath, "from-git-path", es.GitPath, "Git path for the environment template")

	command.MarkFlagsMutuallyExclusive("from-git", "from-template", "from-path", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-branch", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-path", "from-git-repo")

	_ = command.MarkFlagFilename("from-path", "yaml", "yml")
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	editConfigurationOptions := environment.NewEditConfigurationOptions("")
	editSource := EditSource{}

	command := &cobra.Command{
		Use: "update-configuration",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editConfigurationOptions.ID = settings.Profile.Context.Environment

			if err := parseEditConfigurationOptions(editSource, editConfigurationOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := environment.EditConfiguration(editConfigurationOptions)
			if err != nil {
				var apiError api.Error

				if errors.As(err, &apiError) {
					return handleEditErrors(cmd, apiError, editConfigurationOptions)
				}

				return lib.FormatCommandError(cmd, err)
			}

			if !editConfigurationOptions.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &editConfigurationOptions.DeployOptions
			deployOptions.ID = model.GetId()

			return HandleDeploy(cmd, deployOptions, "updated", editConfigurationOptions.K8SIntegration, settings.IsStylish())
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	editConfigurationOptions.UpdateFlagSet(flags)
	editSource.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}

func handleEditErrors(cmd *cobra.Command, apiError api.Error, editOptions *environment.EditConfigurationOptions) error {
	genesisName := getEditGenesisName(editOptions)

	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		cmd.Printf("Problem with %s: %s\n", genesisName, violation.GetMessage())
	}

	return lib.ErrGeneric
}

func getEditGenesisName(editOptions *environment.EditConfigurationOptions) string {
	configuration := editOptions.Configuration

	if configuration.FromGitSpec != nil {
		return "--from-git"
	}

	if configuration.FromTemplate != nil {
		return "--from-template"
	}

	if configuration.FromString != nil {
		return "--from-path"
	}

	return "arguments"
}

//nolint:dupl
func parseEditConfigurationOptions(
	editSource EditSource,
	editConfigurationOptions *environment.EditConfigurationOptions,
) error {
	editConfigurationOptions.Configuration = &sdk.EnvironmentEditConfigurationConfiguration{}

	if editSource.Git != "" {
		fromGitSpec := sdk.NewFromGitSpec()
		fromGitSpec.Spec = &editSource.Git

		editConfigurationOptions.Configuration.FromGitSpec = fromGitSpec

		return nil
	}

	if editSource.TemplateID != "" {
		fromTemplate := sdk.NewFromTemplate()
		fromTemplate.Template = &editSource.TemplateID

		editConfigurationOptions.Configuration.FromTemplate = fromTemplate

		return nil
	}

	if editSource.YamlPath != "" {
		fromString := sdk.NewFromString()

		bytes, err := readFile(editSource.YamlPath)
		if err != nil {
			return err
		}

		content := string(bytes)
		fromString.Yaml = &content

		editConfigurationOptions.Configuration.FromString = fromString

		return nil
	}

	if editSource.GitRepo != "" {
		fromGit := sdk.NewFromGit()
		fromGit.Url = &editSource.GitRepo
		fromGit.Branch = &editSource.GitBranch
		fromGit.YamlPath = &editSource.GitPath

		editConfigurationOptions.Configuration.FromGit = fromGit

		return nil
	}

	return errCreateSourceNotProvided
}
