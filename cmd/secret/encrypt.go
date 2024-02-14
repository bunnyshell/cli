package secret

import (
	"errors"
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/secret"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var (
	errMissingValue        = errors.New("the plain value must be provided")
	errMultipleValueInputs = errors.New("the value must be provided either by argument or by stdin, not both")
)

var encryptCommandExample = heredoc.Docf(`
	%[1]s%[2]s secret encrypt --organization dMVwZO5jGN "my plain secret value"
	%[1]scat plain.txt | %[2]s secret encrypt --organization dMVwZO5jGN
`, "\t", build.Name)

func init() {
	settings := config.GetSettings()

	encryptOptions := secret.NewEncryptOptions()

	command := &cobra.Command{
		Use: "encrypt <value>",

		Short:   "Encrypts a secret for the given organization",
		Example: encryptCommandExample,

		Args:              cobra.MaximumNArgs(1),
		ArgAliases:        []string{"value"},
		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			if len(args) == 0 && !hasStdin {
				return errMissingValue
			}

			if len(args) == 1 && hasStdin {
				return errMultipleValueInputs
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if encryptOptions.Organization == "" {
				encryptOptions.Organization = settings.Profile.Context.Organization
			}

			if len(args) == 1 {
				encryptOptions.PlainText = args[0]
			} else {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				encryptOptions.PlainText = string(buf)
			}

			encryptedSecret, err := secret.Encrypt(encryptOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, encryptedSecret)
		},
	}

	encryptOptions.UpdateSelfFlags(command.Flags())

	mainCmd.AddCommand(command)
}
