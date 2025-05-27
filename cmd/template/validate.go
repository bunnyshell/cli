package template

import (
	"errors"
	"io"
	"os"
	"strings"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

type ValidateGit struct {
	URL       string
	Branch    string
	Directory string
}

type ValidateFile struct {
	BunnyshellYaml string
	TemplateYaml   string
}

type ValidateDir struct {
	Directory string
}

type ValidateSource struct {
	ValidateGit
	ValidateFile
	ValidateDir
}

var errCreateSourceNotProvided = errors.New("git repository or yaml files must be provided")

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	validateOptions := template.NewValidateOptions()
	validateSource := ValidateSource{}

	command := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"val"},
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if validateSource.BunnyshellYaml == "" && validateSource.URL == "" && validateSource.ValidateDir.Directory == "" {
				return errCreateSourceNotProvided
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptions.Organization = settings.Profile.Context.Organization

			if err := parseValidateOptions(validateOptions, validateSource); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if _, err := template.Validate(validateOptions); err != nil {
				var apiError api.Error

				if errors.As(err, &apiError) {
					return handleValidateErrors(cmd, apiError)
				}

				return lib.FormatCommandError(cmd, err)
			}

			cmd.Println("Template is valid.")

			return nil
		},
	}

	flags := command.Flags()

	orgFlag := options.Organization.AddFlag("organization", "Validate template in this organization")
	flags.AddFlag(orgFlag)
	_ = command.MarkFlagRequired(orgFlag.Name)

	validateOptions.UpdateFlagSet(flags)

	validateGit := &validateSource.ValidateGit
	flags.StringVar(&validateGit.URL, "git-url", validateGit.URL, "Git repository URL")
	flags.StringVar(&validateGit.Branch, "git-branch", validateGit.Branch, "Git repository branch")
	flags.StringVar(&validateGit.Directory, "git-directory", validateGit.Directory, "Git repository directory")
	command.MarkFlagsRequiredTogether("git-url", "git-branch", "git-directory")

	validateFile := &validateSource.ValidateFile
	flags.StringVar(&validateFile.BunnyshellYaml, "bunnyshell-yaml", validateFile.BunnyshellYaml, "Bunnyshell yaml file")
	flags.StringVar(&validateFile.TemplateYaml, "template-yaml", validateFile.TemplateYaml, "Template yaml file")
	command.MarkFlagsRequiredTogether("bunnyshell-yaml", "template-yaml")

	validateDir := &validateSource.ValidateDir
	flags.StringVar(&validateDir.Directory, "directory", validateDir.Directory, "Directory")

	command.MarkFlagsMutuallyExclusive("git-url", "bunnyshell-yaml", "directory")

	mainCmd.AddCommand(command)
}

func handleValidateErrors(cmd *cobra.Command, apiError api.Error) error {
	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		if violation.GetPropertyPath() == "source.dirPath" {
			cmd.Printf("%s\n", violation.GetMessage())
		}

		if strings.HasPrefix(violation.GetPropertyPath(), "bunnyshellYaml") {
			cmd.Printf("%s: %s\n", "Error validating bunnyshell.yaml:", violation.GetMessage())
		}

		if strings.HasPrefix(violation.GetPropertyPath(), "templateYaml") {
			cmd.Printf("%s: %s\n", "Error validating template.yaml:", violation.GetMessage())
		}
	}

	return lib.ErrGeneric
}

func parseValidateOptions(validateOptions *template.ValidateOptions, validateSource ValidateSource) error {
	source, err := getSource(validateSource, *validateOptions)
	if err != nil {
		return err
	}

	validateOptions.TemplateValidateAction.SetSource(*source)

	return nil
}

func getSource(validateSource ValidateSource, validateOptions template.ValidateOptions) (*sdk.TemplateValidateActionSource, error) {
	if validateSource.ValidateGit.URL != "" {
		source := fromGitSource(validateSource.ValidateGit, validateOptions.Organization)

		if validateOptions.WithComponents {
			source.SetValidateComponents(true)
		}

		if !validateOptions.AllowExtraFields {
			source.SetValidateAllowExtraFields(false)
		}

		action := sdk.ValidateSourceGitAsTemplateValidateActionSource(source)

		return &action, nil
	}

	if validateSource.ValidateDir.Directory != "" {
		// proxy to ValidateFile logic
		validateSource.ValidateFile.BunnyshellYaml = validateSource.ValidateDir.Directory + "/bunnyshell.yaml"
		validateSource.ValidateFile.TemplateYaml = validateSource.ValidateDir.Directory + "/template.yaml"
	}

	if validateSource.ValidateFile.BunnyshellYaml != "" {
		source, err := fromFile(validateSource.ValidateFile)
		if err != nil {
			return nil, err
		}

		if validateOptions.WithComponents {
			source.SetValidateComponents(true)
			source.SetValidateForOrganizationId(validateOptions.Organization)
		}

		if !validateOptions.AllowExtraFields {
			source.SetValidateAllowExtraFields(false)
		}

		action := sdk.ValidateSourceStringAsTemplateValidateActionSource(source)

		return &action, nil
	}

	return nil, errCreateSourceNotProvided
}

func fromGitSource(git ValidateGit, organization string) *sdk.ValidateSourceGit {
	return sdk.NewValidateSourceGit(git.URL, git.Branch, git.Directory, organization)
}

func fromFile(git ValidateFile) (*sdk.ValidateSourceString, error) {
	bunnyshellYaml, err := readFile(git.BunnyshellYaml)
	if err != nil {
		return nil, err
	}

	templateYaml, err := readFile(git.TemplateYaml)
	if err != nil {
		return nil, err
	}

	return sdk.NewValidateSourceString(string(bunnyshellYaml), string(templateYaml)), nil
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
