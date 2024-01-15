package secret

import (
	"bunnyshell.com/cli/pkg/api/secret"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var encryptDefinitionCommandExample = heredoc.Docf(`
	%[1]s%[2]s secret encrypt-definition --organization dMVwZO5jGN --file plain.txt
	%[1]scat plain.txt | %[2]s secret encrypt-definition --organization dMVwZO5jGN
`, "\t", build.Name)

func init() {
	transcriptConfigurationOptions := secret.NewTranscriptConfigurationOptions()

	command := &cobra.Command{
		Use: "encrypt-definition",

		Short:   "Encrypts an environment definition for the given organization",
		Example: encryptDefinitionCommandExample,

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateDefinitionCommand(transcriptConfigurationOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			encryptedDefinition, err := executeTranscriptConfiguration(transcriptConfigurationOptions, secret.TranscriptModeObfuscated)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Println(encryptedDefinition)

			return nil
		},
	}

	transcriptConfigurationOptions.UpdateSelfFlags(command.Flags())

	mainCmd.AddCommand(command)
}
