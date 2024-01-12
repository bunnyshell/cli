package secret

import (
	"errors"
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/secret"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
)

var (
	errYamlBlankValue          = errors.New("the environment definition provided is blank")
	errYamlMissingValue        = errors.New("the environment definition must be provided")
	errYamlFileNotFound        = errors.New("the provided filepath does not exist")
	errMultipleYamlValueInputs = errors.New("the environment definition must be provided either by argument or by stdin, not both")
)

func executeTranscriptConfiguration(options *secret.TranscriptConfigurationOptions, mode secret.TranscriptMode) (string, error) {
	settings := config.GetSettings()

	options.Mode = string(mode)
	if options.Organization == "" {
		options.Organization = settings.Profile.Context.Organization
	}

	err := loadDefinition(options)
	if err != nil {
		return "", err
	}

	if options.Yaml == "" {
		return "", errYamlBlankValue
	}

	resultedDefinition, err := secret.TranscriptConfiguration(options)
	if err != nil {
		return "", err
	}

	return string(resultedDefinition.Bytes), nil
}

func validateDefinitionCommand(options *secret.TranscriptConfigurationOptions) error {
	hasStdin, err := isStdinPresent()
	if err != nil {
		return err
	}

	if options.DefinitionFilePath == "" && !hasStdin {
		return errYamlMissingValue
	}

	if options.DefinitionFilePath != "" {
		if hasStdin {
			return errMultipleYamlValueInputs
		}

		exists, err := util.FileExists(options.DefinitionFilePath)
		if err != nil {
			return err
		}
		if !exists {
			return errYamlFileNotFound
		}
	}

	return nil
}

func loadDefinition(options *secret.TranscriptConfigurationOptions) error {
	if options.DefinitionFilePath != "" {
		contents, err := loadDefinitionFromFile(options.DefinitionFilePath)
		if err != nil {
			return err
		}

		options.Yaml = contents
	} else {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		options.Yaml = string(buf)
	}

	return nil
}

func loadDefinitionFromFile(filepath string) (string, error) {
	exists, err := util.FileExists(filepath)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("the provided file does not exist")
	}

	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(fileContents), nil
}
