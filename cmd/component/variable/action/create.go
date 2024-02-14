package action

import (
	"fmt"
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/component_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	settings := config.GetSettings()

	createOptions := component_variable.NewCreateOptions()
	var componentName string

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			if !flags.Changed("value") && !hasStdin {
				return errMissingValue
			}

			if flags.Changed("value") && hasStdin {
				return errMultipleValueInputs
			}

			if !flags.Changed("component-name") {
				cmd.MarkFlagRequired("component")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.ServiceComponent = settings.Profile.Context.ServiceComponent

			if componentName != "" {
				component, err := findComponentByName(componentName, &settings.Profile)
				if err != nil {
					return err
				}
				createOptions.ServiceComponent = component.GetId()
			}

			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			if hasStdin {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				createOptions.Value = string(buf)
			}

			model, err := component_variable.Create(createOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	updateComponentIdentifierFlags(command, &componentName)

	createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func findComponentByName(componentName string, profile *config.Profile) (*sdk.ComponentCollection, error) {
	matchedComponents, err := findComponentsByName(componentName, profile)
	if err != nil {
		return nil, fmt.Errorf("failed while searching for the component: %w", err)
	}

	if matchedComponents.GetTotalItems() == 0 {
		return nil, fmt.Errorf("component '%s' not found", componentName)
	}

	if matchedComponents.GetTotalItems() > 1 {
		return nil, fmt.Errorf("multiple components with the name '%s' were found", componentName)
	}

	return &matchedComponents.GetEmbedded().Item[0], nil
}
