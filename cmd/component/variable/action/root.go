package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/component_variable"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var (
	errMissingValue        = errors.New("the plain value must be provided")
	errMultipleValueInputs = errors.New("the value must be provided either by argument or by stdin, not both")
	errComponentRequired   = errors.New("either the 'component' or the 'component-name' arguments are required")
)

var mainCmd = &cobra.Command{}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func GetIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available component variables with "%s variables list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Component Variable Id", help)

	return idOption
}

func updateComponentIdentifierFlags(cmd *cobra.Command, componentName *string) {
	options := config.GetOptions()

	flags := cmd.Flags()

	flags.AddFlag(options.ServiceComponent.AddFlagWithExtraHelp(
		"component",
		"Component for the variable",
		"Components contain multiple variables",
	))

	flags.StringVar(componentName, "component-name", *componentName, "Component Name")
	cmd.MarkFlagsMutuallyExclusive("component-name", "component")

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
}

func findComponentsByName(componentName string, profile *config.Profile) (*sdk.PaginatedComponentCollection, error) {
	listOptions := component.NewListOptions()

	listOptions.Organization = profile.Context.Organization
	listOptions.Project = profile.Context.Project
	listOptions.Environment = profile.Context.Environment
	listOptions.Name = componentName

	components, error := component.List(listOptions)
	if error != nil {
		return nil, error
	}

	return components, nil
}

func findComponentVariablesByName(componentVariableName string, componentName string, profile *config.Profile) (*sdk.PaginatedServiceComponentVariableCollection, error) {
	listOptions := component_variable.NewListOptions()

	listOptions.Organization = profile.Context.Organization
	listOptions.Project = profile.Context.Project
	listOptions.Environment = profile.Context.Environment
	listOptions.Component = profile.Context.ServiceComponent
	listOptions.Name = componentVariableName

	if componentName != "" {
		listOptions.ComponentName = componentName
	}

	component_variables, error := component_variable.List(listOptions)
	if error != nil {
		return nil, error
	}

	return component_variables, nil
}

func findComponentVariableByName(componentVariableName string, componentName string, profile *config.Profile) (*sdk.ServiceComponentVariableCollection, error) {
	matchedComponentVariables, err := findComponentVariablesByName(componentVariableName, componentName, profile)
	if err != nil {
		return nil, fmt.Errorf("failed while searching for the component variable: %w", err)
	}

	if matchedComponentVariables.GetTotalItems() == 0 {
		return nil, fmt.Errorf("component variable '%s' not found", componentVariableName)
	}

	if matchedComponentVariables.GetTotalItems() > 1 {
		return nil, fmt.Errorf("multiple component variables '%s' were found", componentVariableName)
	}

	return &matchedComponentVariables.GetEmbedded().Item[0], nil
}
