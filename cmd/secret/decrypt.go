package secret

import (
	"errors"
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/secret"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var (
	errMissingEncryptedExpression        = errors.New("the encrypted expression must be provided")
	errMultipleEncryptedExpressionInputs = errors.New("the encrypted expression must be provided either by argument or by stdin, not both")
)

var decryptCommandExample = heredoc.Docf(`
	%[1]s%[2]s secret decrypt --organization dMVwZO5jGN "ENCRYPTED[F67baQQZ6XXHNxcmRz]"
	%[1]scat encrypted.txt | %[2]s secret decrypt --organization dMVwZO5jGN
`, "\t", build.Name)

func init() {
	settings := config.GetSettings()

	decryptOptions := secret.NewDecryptOptions()

	command := &cobra.Command{
		Use: "decrypt <expression>",

		Short:   "Decrypts a secret expression of the given organization",
		Example: decryptCommandExample,

		Args:              cobra.MaximumNArgs(1),
		ArgAliases:        []string{"expression"},
		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := isStdinPresent()
			if err != nil {
				return err
			}

			if len(args) == 0 && !hasStdin {
				return errMissingEncryptedExpression
			}

			if len(args) == 1 && hasStdin {
				return errMultipleEncryptedExpressionInputs
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if decryptOptions.Organization == "" {
				decryptOptions.Organization = settings.Profile.Context.Organization
			}

			if len(args) == 1 {
				decryptOptions.Expression = args[0]
			} else {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				decryptOptions.Expression = string(buf)
			}

			decryptedSecret, err := secret.Decrypt(decryptOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, decryptedSecret)
		},
	}

	decryptOptions.UpdateSelfFlags(command.Flags())

	mainCmd.AddCommand(command)
}
