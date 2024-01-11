package action

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/api/component/git"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	githelper "bunnyshell.com/cli/pkg/helper/git"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/slice"
)

type EditComponentsSource struct {
	// filters
	Component string
	GitSource string

	// updates
	GitTarget string

	// deployment
	K8SIntegration string
}

var commandExample = fmt.Sprintf(`This command updates the Git details for components in an environment.

You can update the Git details for a specific component in an environment by using the --component-name flag:
%[1]s%[2]s env update-components --component-name my-component --git-target https://github.com/my-fork/my-repo@my-main

You can update all components matching a specific repository:
%[1]s%[2]s env update-components --git-source https://github.com/original/repo --git-target git@github.com/my-fork/my-repo

You can update all components matching a specific branch:
%[1]s%[2]s env update-components --git-source @main --git-target @feature-branch`, "\t", build.Name)

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

			matchedComponentNames := componentToNamesList(matched)
			cmd.Printf(`Updating components "%s"%s`, strings.Join(matchedComponentNames, ", "), "\n\n")

			model, err := environment.EditComponents(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("Successfully updated Git details...%s", "\n\n")

			if !editOptions.WithDeploy {
				return showUpdateGitInfo(cmd, model.GetId(), matchedComponentNames)
			}

			deployOptions := &editOptions.DeployOptions
			deployOptions.ID = model.GetId()

			if err = HandleDeploy(cmd, deployOptions, "updated", editSource.K8SIntegration); err != nil {
				return err
			}

			return showUpdateGitInfo(cmd, model.GetId(), matchedComponentNames)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.AddFlag("environment", "Update components for environment", util.FlagRequired))

	editOptions.UpdateFlagSet(flags)

	flags.StringVar(&editSource.GitTarget, "git-target", editSource.GitTarget, "Target git spec (e.g. https://github.com/fork/templates@main)")

	targetFlag := flags.Lookup("git-target")
	util.MarkFlagRequiredWithHelp(targetFlag, "Update components git repository and branch. Example: https://github.com/my-fork/my-repo@my-branch")

	flags.StringVar(&editSource.GitSource, "git-source", editSource.GitSource, "Filter by git spec (e.g. https://github.com/bunnyshell/templates@main)")
	flags.StringVar(&editSource.Component, "component-name", editSource.Component, "Filter by component name")
	command.MarkFlagsMutuallyExclusive("git-source", "component-name")

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

	if editOptions.Component != "" {
		aggOptions.Name = editOptions.Component
	}

	return git.Aggregate(aggOptions)
}

func showUpdateGitInfo(cmd *cobra.Command, environment string, componentNames []string) error {
	listOptions := git.NewAggregateOptions()
	listOptions.Environment = environment

	agg, err := git.Aggregate(listOptions)
	if err != nil {
		return err
	}

	filtered := []sdk.ComponentGitCollection{}
	for _, element := range agg {
		if slice.ContainsString(componentNames, element.GetName(), nil) {
			filtered = append(filtered, element)
		}
	}

	return lib.FormatCommandData(cmd, filtered)
}

func componentToNamesList(matched []sdk.ComponentGitCollection) []string {
	components := []string{}

	for _, item := range matched {
		components = append(components, item.GetName())
	}

	return components
}

func fillWithGitSpec(editSource *EditComponentsSource, editOptions *environment.EditComponentOptions) error {
	if editSource.GitTarget != "" {
		target, branch, err := githelper.ParseGitSec(editSource.GitTarget)
		if err != nil {
			return fmt.Errorf("invalid git spec for %s: %w", editSource.GitTarget, err)
		}

		editOptions.TargetRepository = target
		editOptions.TargetBranch = branch
	}

	if editSource.Component != "" {
		editOptions.Component = editSource.Component

		return nil
	}

	if editSource.GitSource == "" {
		return nil
	}

	target, branch, err := githelper.ParseGitSec(editSource.GitSource)
	if err != nil {
		return fmt.Errorf("invalid git spec for %s: %w", editSource.GitSource, err)
	}

	editOptions.SourceRepository = target
	editOptions.SourceBranch = branch

	return nil
}
