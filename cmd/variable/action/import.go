package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/api/variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/enum"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BulkImport struct {
	Vars    map[string]string `json:"vars" yaml:"vars"`
	Secrets map[string]string `json:"secrets" yaml:"secrets"`
}

func init() {
	varFile := ""
	secretFile := ""
	options := config.GetOptions()
	data := BulkImport{
		Vars:    make(map[string]string),
		Secrets: make(map[string]string),
	}

	command := &cobra.Command{
		Use: "import",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			//todo add this to preRunE
			if varFile == "" && secretFile == "" {
				return errors.New("must provide a either a var or secret file")
			}

			if varFile != "" {
				if err := readFile(varFile, &data.Vars); err != nil {
					return err
				}
			}

			if secretFile != "" {
				if err := readFile(secretFile, &data.Secrets); err != nil {
					return err
				}
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := config.GetSettings()

			if len(data.Vars) > 0 {
				for key, value := range data.Vars {
					createOptions := variable.NewCreateOptions()
					createOptions.Environment = settings.Profile.Context.Environment
					createOptions.Name = key
					createOptions.Value = value

					model, err := variable.Create(createOptions)
					if err != nil {
						// Log the error or notify that the particular variable couldn't be created
						fmt.Printf("Error creating variable %s: %v\n", key, err)
					} else {
						// Successfully created model, output the results immediately
						lib.FormatCommandData(cmd, model)
					}
				}
			}

			if len(data.Secrets) > 0 {
				for key, value := range data.Secrets {
					createOptions := variable.NewCreateOptions()
					createOptions.Environment = settings.Profile.Context.Environment
					createOptions.Name = key
					createOptions.Value = value
					createOptions.IsSecret = enum.BoolTrue
					model, err := variable.Create(createOptions)
					if err != nil {
						// Log the error or notify that the particular variable couldn't be created
						fmt.Printf("Error creating variable %s: %v\n", key, err)
					} else {
						// Successfully created model, output the results immediately
						lib.FormatCommandData(cmd, model)
					}
				}
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.StringVar(&varFile, "vars-file", varFile, "File to import variables from")
	flags.StringVar(&secretFile, "secrets-file", secretFile, "File to import secrets from")

	flags.AddFlag(options.Environment.AddFlagWithExtraHelp(
		"environment",
		"Environment for the variable",
		"Environments contain multiple variables",
		util.FlagRequired,
	))

	// createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func readFile(fileName string, data *map[string]string) error {
	viper := viper.New()
	viper.SetConfigFile(fileName)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		// @review update go:1.20 errors.join
		// return fmt.Errorf("%w: %s", ErrConfigLoad, err.Error())
		return err
	}

	if err := viper.Unmarshal(&data); err != nil {
		// @review update go:1.20 errors.join
		// return fmt.Errorf("%w: %s", ErrConfigLoad, err.Error())
		return err
	}

	return nil
}
