package action

import (
	"fmt"
	"net/url"
	"strings"

	"bunnyshell.com/cli/pkg/api/component/git"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

type EditComponentsSource struct {
	// filters
	Component string
	Source    string

	// updates
	Target string

	// deployment
	K8SIntegration string
}

var commandExample = fmt.Sprintf(`This command updates the Git details for components in an environment.

You can update the Git details for a specific component in an environment by using the --component-name flag:
%[1]s%[2]s env update-components --component-name my-component --target https://github.com/my-fork/my-repo@my-main

You can update all components matching a specific repository:
%[1]s%[2]s env update-components --source https://github.com/original/repo --target git@github.com/my-fork/my-repo

You can update all components matching a specific branch:
%[1]s%[2]s env update-components --source @main --target @feature-branch`, "\t", build.Name)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	editOptions := environment.NewEditComponentOptions()
	editSource := &EditComponentsSource{}

	command := &cobra.Command{
		Use: "update-components",

		Example: commandExample,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editOptions.ID = settings.Profile.Context.Environment

			if err := fillWithGitSpec(editSource, editOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			matched, err := findMatchedComponents(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if len(matched) == 0 {
				cmd.Println("No components matched the filter")

				return nil
			}

			cmd.Printf(`Updating components "%s"%s`, componentToString(matched), "\n\n")

			model, err := environment.EditComponents(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("Successfully updated Git details...%s", "\n\n")

			if !editOptions.WithDeploy {
				return showGitInfo(cmd, model.GetId())
			}

			deployOptions := &editOptions.DeployOptions
			deployOptions.ID = model.GetId()

			if err = handleDeploy(cmd, deployOptions, "updated", editSource.K8SIntegration); err != nil {
				return err
			}

			return showGitInfo(cmd, model.GetId())
		},
	}

	flags := command.Flags()

	envFlag := options.Environment.AddFlag("environment", "Update components for environment")
	flags.AddFlag(envFlag)
	_ = command.MarkFlagRequired(envFlag.Name)

	editOptions.UpdateFlagSet(flags)

	flags.StringVar(&editSource.Target, "target", editSource.Target, "Target git spec (e.g. https://github.com/fork/templates@main)")

	_ = command.MarkFlagRequired("target")

	flags.StringVar(&editSource.Source, "source", editSource.Source, "Filter by git spec (e.g. https://github.com/bunnyshell/templates@main)")
	flags.StringVar(&editSource.Component, "component-name", editSource.Component, "Filter by component name")
	command.MarkFlagsMutuallyExclusive("source", "component-name")

	mainCmd.AddCommand(command)
}

func findMatchedComponents(editOptions *environment.EditComponentOptions) ([]sdk.ComponentGitCollection, error) {
	aggOptions := git.NewAggregateOptions()
	aggOptions.Environment = editOptions.ID

	if editOptions.SourceRepository != "" {
		aggOptions.GitRepository = editOptions.SourceRepository
	}

	if editOptions.SourceBranch != "" {
		aggOptions.GitBranch = editOptions.SourceBranch
	}

	return git.Aggregate(aggOptions)
}

func showGitInfo(cmd *cobra.Command, environment string) error {
	listOptions := git.NewAggregateOptions()
	listOptions.Environment = environment

	agg, err := git.Aggregate(listOptions)
	if err != nil {
		return err
	}

	return lib.FormatCommandData(cmd, agg)
}

func componentToString(matched []sdk.ComponentGitCollection) string {
	components := []string{}

	for _, item := range matched {
		components = append(components, item.GetName())
	}

	return strings.Join(components, ", ")
}

func fillWithGitSpec(editSource *EditComponentsSource, editOptions *environment.EditComponentOptions) error {
	if editSource.Target != "" {
		target, branch, err := parseGitSec(editSource.Target)
		if err != nil {
			return fmt.Errorf("invalid git spec for %s: %w", editSource.Target, err)
		}

		editOptions.TargetRepository = target
		editOptions.TargetBranch = branch
	}

	if editSource.Component != "" {
		editOptions.Component = editSource.Component

		return nil
	}

	if editSource.Source == "" {
		return nil
	}

	target, branch, err := parseGitSec(editSource.Source)
	if err != nil {
		return fmt.Errorf("invalid git spec for %s: %w", editSource.Source, err)
	}

	editOptions.SourceRepository = target
	editOptions.SourceBranch = branch

	return nil
}

func parseGitSec(spec string) (string, string, error) {
	if spec[0] == '@' {
		return "", spec[1:], nil
	}

	info, err := url.Parse(spec)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(info.Path, "@") {
		return spec, "", nil
	}

	chunks := strings.SplitN(info.Path, "@", 2)
	info.Path = chunks[0]

	return info.String(), chunks[1], nil
}
