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

	editOptions := environment.NewEditOptions()
	editSource := EditSource{}

	command := &cobra.Command{
		Use: "update",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editOptions.ID = settings.Profile.Context.Environment

			if err := parseEditOptions(editSource, editOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := environment.Edit(editOptions)
			if err != nil {
				var apiError api.Error

				if errors.As(err, &apiError) {
					return handleEditErrors(cmd, apiError, editOptions)
				}

				return lib.FormatCommandError(cmd, err)
			}

			if !editOptions.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &editOptions.DeployOptions
			deployOptions.ID = model.GetId()

			return handleDeploy(cmd, deployOptions, "updated", editOptions.K8SIntegration)
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	editOptions.UpdateFlagSet(flags)

	editSource.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}

func handleEditErrors(cmd *cobra.Command, apiError api.Error, editOptions *environment.EditOptions) error {
	genesisName := getEditGenesisName(editOptions)

	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		cmd.Printf("Problem with %s: %s\n", genesisName, violation.GetMessage())
	}

	return lib.ErrGeneric
}

func getEditGenesisName(editOptions *environment.EditOptions) string {
	if editOptions.Genesis.FromGitSpec != nil {
		return "--from-git"
	}

	if editOptions.Genesis.FromTemplate != nil {
		return "--from-template"
	}

	if editOptions.Genesis.FromString != nil {
		return "--from-path"
	}

	return "arguments"
}

//nolint:dupl
func parseEditOptions(editSource EditSource, editOptions *environment.EditOptions) error {
	editOptions.Genesis = &sdk.EnvironmentEditActionGenesis{}

	if editSource.Git != "" {
		fromGitSpec := sdk.NewFromGitSpec()
		fromGitSpec.Spec = &editSource.Git

		editOptions.Genesis.FromGitSpec = fromGitSpec

		return nil
	}

	if editSource.TemplateID != "" {
		fromTemplate := sdk.NewFromTemplate()
		fromTemplate.Template = &editSource.TemplateID

		editOptions.Genesis.FromTemplate = fromTemplate

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

		editOptions.Genesis.FromString = fromString

		return nil
	}

	if editSource.GitRepo != "" {
		fromGit := sdk.NewFromGit()
		fromGit.Url = &editSource.GitRepo
		fromGit.Branch = &editSource.GitBranch
		fromGit.YamlPath = &editSource.GitPath

		editOptions.Genesis.FromGit = fromGit

		return nil
	}

	return errCreateSourceNotProvided
}
