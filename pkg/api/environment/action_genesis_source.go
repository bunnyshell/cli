package environment

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

type GenesisSourceOptions struct {
	TemplateID            string
	TemplateVariablePairs []string

	Git string

	YamlPath string

	GitRepo   string
	GitBranch string
	GitPath   string
}

const variableSplitSize = 2

var (
	errGenesisSourceNotProvided = errors.New("template id, content or git repository must be provided")
	errInvalidVarDefinition     = errors.New("invalid template variable definition")
	errUnknownVar               = errors.New("unknown variable")
	errUnknownEnum              = errors.New("unknown enum value")
)

func NewGenesisSourceOptions() *GenesisSourceOptions {
	return &GenesisSourceOptions{}
}

func (gs *GenesisSourceOptions) updateCommandFlags(command *cobra.Command, genesis string) {
	flags := command.Flags()

	flags.StringVar(&gs.TemplateID, "from-template", gs.TemplateID, "Use a TemplateID during environment "+genesis)
	flags.StringArrayVar(
		&gs.TemplateVariablePairs,
		"template-var",
		gs.TemplateVariablePairs,
		"Template variables to use during environment "+genesis,
	)

	flags.StringVar(&gs.YamlPath, "from-path", gs.YamlPath, "Use a local bunnyshell.yaml during environment "+genesis)

	flags.StringVar(&gs.Git, "from-git", gs.Git, "Use a template git repository during environment "+genesis)

	flags.StringVar(&gs.GitRepo, "from-git-repo", gs.GitRepo, "Git repository for the environment template")
	flags.StringVar(&gs.GitBranch, "from-git-branch", gs.GitBranch, "Git branch for the environment template")
	flags.StringVar(&gs.GitPath, "from-git-path", gs.GitPath, "Git path for the environment template")

	command.MarkFlagsMutuallyExclusive("from-git", "from-template", "from-path", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-branch", "from-git-repo")
	command.MarkFlagsRequiredTogether("from-git-path", "from-git-repo")

	_ = command.MarkFlagFilename("from-path", "yaml", "yml")
}

func (gs *GenesisSourceOptions) validate() error {
	if gs.Git == "" && gs.TemplateID == "" && gs.YamlPath == "" && gs.GitRepo == "" {
		return errGenesisSourceNotProvided
	}

	return nil
}
func (gs *GenesisSourceOptions) handleError(cmd *cobra.Command, apiError api.Error) error {
	genesisName := gs.getGenesisName()

	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		cmd.Printf("Problem with %s: %s\n", genesisName, violation.GetMessage())
	}

	return lib.ErrGeneric
}

func (gs *GenesisSourceOptions) getGenesisName() string {
	if gs.Git != "" {
		return "--from-git"
	}

	if gs.TemplateID != "" {
		return "--from-template"
	}

	if gs.YamlPath != "" {
		return "--from-path"
	}

	return "arguments"
}

func (gs *GenesisSourceOptions) getGenesis() (*sdk.FromGit, *sdk.FromGitSpec, *sdk.FromString, *sdk.FromTemplate, error) {
	if gs.Git != "" {
		return nil, gs.getFromGitSpec(), nil, nil, nil
	}

	if gs.GitRepo != "" {
		return gs.getFromGit(), nil, nil, nil, nil
	}

	if gs.TemplateID != "" {
		fromTemplate, err := gs.getFromTemplate()
		if err != nil {
			return nil, nil, nil, nil, err
		}

		return nil, nil, nil, fromTemplate, nil
	}

	if gs.YamlPath != "" {
		fromString, err := gs.getFromString()
		if err != nil {
			return nil, nil, nil, nil, err
		}

		return nil, nil, fromString, nil, nil
	}

	return nil, nil, nil, nil, errGenesisSourceNotProvided
}

func (gs *GenesisSourceOptions) getFromGit() *sdk.FromGit {
	fromGit := sdk.NewFromGit()
	fromGit.Url = &gs.GitRepo
	fromGit.Branch = &gs.GitBranch
	fromGit.YamlPath = &gs.GitPath

	return fromGit

}

func (gs *GenesisSourceOptions) getFromGitSpec() *sdk.FromGitSpec {
	fromGitSpec := sdk.NewFromGitSpec()
	fromGitSpec.Spec = &gs.Git

	return fromGitSpec
}

func (gs *GenesisSourceOptions) getFromString() (*sdk.FromString, error) {
	fromString := sdk.NewFromString()

	bytes, err := readFile(gs.YamlPath)
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	fromString.Yaml = &content

	return fromString, nil
}

func (gs *GenesisSourceOptions) getFromTemplate() (*sdk.FromTemplate, error) {
	fromTemplate := sdk.NewFromTemplate()
	fromTemplate.Template = &gs.TemplateID

	if len(gs.TemplateVariablePairs) > 0 {
		templateVariablesSchema, schemaError := getTemplateVariableSchema(gs.TemplateID)
		if schemaError != nil {
			return nil, schemaError
		}

		variables := map[string]sdk.FromTemplateVariablesValue{}
		for _, pair := range gs.TemplateVariablePairs {
			name, value, err := parseDefinition(pair, templateVariablesSchema)
			if err != nil {
				return nil, err
			}

			variables[*name] = *value
		}

		fromTemplate.SetVariables(variables)
	}

	return fromTemplate, nil
}

func getTemplateVariableSchema(templateID string) ([]sdk.TemplateItemVariablesSchemaInner, error) {
	templateItem, err := template.Get(template.NewItemOptions(templateID))
	if err != nil {
		return nil, err
	}

	return templateItem.GetVariablesSchema(), nil
}

func parseDefinition(
	definition string,
	templateVariablesSchema []sdk.TemplateItemVariablesSchemaInner,
) (*string, *sdk.FromTemplateVariablesValue, error) {
	parts := strings.SplitN(definition, "=", variableSplitSize)
	if len(parts) != variableSplitSize {
		return nil, nil, fmt.Errorf("%w: %s", errInvalidVarDefinition, definition)
	}

	name := parts[0]
	value, err := getVariableValue(name, parts[1], templateVariablesSchema)

	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s - %w", errInvalidVarDefinition, definition, err)
	}

	return &name, value, nil
}

func getVariableValue(
	name string,
	stringValue string,
	templateVariablesSchema []sdk.TemplateItemVariablesSchemaInner,
) (*sdk.FromTemplateVariablesValue, error) {
	for _, schema := range templateVariablesSchema {
		switch {
		case schema.BooleanTypeItem != nil && schema.BooleanTypeItem.GetName() == name:
			return stringToBoolVariableValue(stringValue), nil
		case schema.IntegerTypeItem != nil && schema.IntegerTypeItem.GetName() == name:
			return stringToIntVariableValue(stringValue), nil
		case schema.FloatTypeItem != nil && schema.FloatTypeItem.GetName() == name:
			return stringToFloatVariableValue(stringValue), nil
		case schema.StringTypeItem != nil && schema.StringTypeItem.GetName() == name:
			return stringVariableValue(stringValue), nil
		case schema.EnumTypeItem != nil && schema.EnumTypeItem.GetName() == name:
			return stringToEnumValue(stringValue, schema.EnumTypeItem.GetValues())
		}
	}

	return nil, fmt.Errorf("%w %s", errUnknownVar, name)
}

func stringToBoolVariableValue(stringValue string) *sdk.FromTemplateVariablesValue {
	boolValue := cast.ToBool(stringValue)
	variable := sdk.BoolAsFromTemplateVariablesValue(&boolValue)

	return &variable
}

func stringToIntVariableValue(stringValue string) *sdk.FromTemplateVariablesValue {
	intValue := cast.ToInt32(stringValue)
	variable := sdk.Int32AsFromTemplateVariablesValue(&intValue)

	return &variable
}

func stringToFloatVariableValue(stringValue string) *sdk.FromTemplateVariablesValue {
	floatValue := cast.ToFloat32(stringValue)
	variable := sdk.Float32AsFromTemplateVariablesValue(&floatValue)

	return &variable
}

func stringVariableValue(stringValue string) *sdk.FromTemplateVariablesValue {
	variable := sdk.StringAsFromTemplateVariablesValue(&stringValue)

	return &variable
}

func stringToEnumValue(stringValue string, values []sdk.EnumTypeItemValuesInner) (*sdk.FromTemplateVariablesValue, error) {
	boolValue := cast.ToBool(stringValue)
	int32Value := cast.ToInt32(stringValue)
	float32Value := cast.ToFloat32(stringValue)

	for _, item := range values {
		switch {
		case item.BooleanValueItem != nil && item.BooleanValueItem.GetValue() == boolValue:
			variable := sdk.BoolAsFromTemplateVariablesValue(&boolValue)

			return &variable, nil
		case item.IntegerValueItem != nil && item.IntegerValueItem.GetValue() == int32Value:
			variable := sdk.Int32AsFromTemplateVariablesValue(&int32Value)

			return &variable, nil
		case item.FloatValueItem != nil && item.FloatValueItem.GetValue() == float32Value:
			variable := sdk.Float32AsFromTemplateVariablesValue(&float32Value)

			return &variable, nil
		case item.StringValueItem != nil && item.StringValueItem.GetValue() == stringValue:
			variable := sdk.StringAsFromTemplateVariablesValue(&stringValue)

			return &variable, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", errUnknownEnum, stringValue)
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
