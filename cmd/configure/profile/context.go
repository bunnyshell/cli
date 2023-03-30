package profile

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/wizard"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()
	profile := &settings.Profile

	command := &cobra.Command{
		Use: "context",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if errors.Is(config.MainManager.Error, config.ErrConfigLoad) {
				return config.MainManager.Error
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := askToFillContext(profile); err != nil {
				return err
			}

			config.MainManager.SetProfile(*profile)

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Updated profile context",
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.ProfileName.GetRequiredFlag("profile"))

	mainCmd.AddCommand(command)
}

func askToFillContextOrSkip(profile *config.Profile) error {
	addContext, err := interactive.ConfirmWithHelp(
		"Do you want to set a context for this profile?",
		"Context is used to determine which organization, project, environment, and component to use when running commands.\n"+
			"See more at https://documentation.bunnyshell.com/docs/bunnyshell-cli-authentication#create-a-profile",
	)
	if err != nil {
		return err
	}

	if !addContext {
		return nil
	}

	return askToFillContext(profile)
}

func askToFillContext(profile *config.Profile) error {
	var err error

	wiz := wizard.New(profile)

	if err = askOrganization(wiz, profile); err != nil {
		return err
	}

	if err = askProject(wiz, profile); err != nil {
		return err
	}

	if err = askEnvironment(wiz, profile); err != nil {
		return err
	}

	if err = askComponent(wiz, profile); err != nil {
		return err
	}

	return nil
}

func askOrganization(wiz *wizard.Wizard, profile *config.Profile) error {
	return addOrRemove("organization", profile.Context.Organization, func() error {
		item, err := wiz.SelectOrganization()
		if err != nil {
			return err
		}

		profile.Context.Organization = item.GetId()

		return nil
	}, func() {
		profile.Context.Organization = ""
	})
}

func askProject(wiz *wizard.Wizard, profile *config.Profile) error {
	return addOrRemove("project", profile.Context.Project, func() error {
		item, err := wiz.SelectProject()
		if err != nil {
			return err
		}

		profile.Context.Project = item.GetId()

		return nil
	}, func() {
		profile.Context.Project = ""
	})
}

func askEnvironment(wiz *wizard.Wizard, profile *config.Profile) error {
	return addOrRemove("environment", profile.Context.Environment, func() error {
		item, err := wiz.SelectEnvironment()
		if err != nil {
			return err
		}

		profile.Context.Environment = item.GetId()

		return nil
	}, func() {
		profile.Context.Environment = ""
	})
}

func askComponent(wiz *wizard.Wizard, profile *config.Profile) error {
	return addOrRemove("component", profile.Context.ServiceComponent, func() error {
		item, err := wiz.SelectComponent()
		if err != nil {
			return err
		}

		profile.Context.ServiceComponent = item.GetId()

		return nil
	}, func() {
		profile.Context.ServiceComponent = ""
	})
}

func addOrRemove(name string, value string, add func() error, remove func()) error {
	if value != "" {
		removeContext, err := interactive.Confirm(fmt.Sprintf("Remove context %s (%s) ?", name, value))
		if err != nil {
			return err
		}

		if !removeContext {
			return nil
		}

		remove()
	}

	addContext, err := interactive.Confirm(fmt.Sprintf("Set context %s ?", name))
	if err != nil {
		return err
	}

	if !addContext {
		return nil
	}

	return add()
}
