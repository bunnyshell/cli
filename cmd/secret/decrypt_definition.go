package secret

import (
	"bunnyshell.com/cli/pkg/api/secret"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var decryptDefinitionCommandExample = heredoc.Docf(`
	%[1]s%[2]s secret decrypt-definition --organization dMVwZO5jGN --file plain.txt
	%[1]scat plain.txt | %[2]s secret decrypt-definition --organization dMVwZO5jGN
`, "\t", build.Name)

var resolvedExpressions bool

func init() {
	transcriptConfigurationOptions := secret.NewTranscriptConfigurationOptions()

	command := &cobra.Command{
		Use: "decrypt-definition",

		Short:   "Decrypts an environment definition for the given organization",
		Example: decryptDefinitionCommandExample,

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateDefinitionCommand(transcriptConfigurationOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			mode := secret.TranscriptModeExposed
			if resolvedExpressions {
				mode = secret.TranscriptModeResolved
			}

			decryptedDefinition, err := executeTranscriptConfiguration(transcriptConfigurationOptions, mode)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Println(decryptedDefinition)

			return nil
		},
	}

	flags := command.Flags()

	transcriptConfigurationOptions.UpdateSelfFlags(flags)
	flags.BoolVar(&resolvedExpressions, "resolved", false, "Resolve the expressions and include directly the value")

	mainCmd.AddCommand(command)
}
