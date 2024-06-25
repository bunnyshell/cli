package environment

import (
	"errors"
	"fmt"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config/enum"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var (
	errUnknownEnvironmentType   = errors.New("unknown environment type")
	errUndefinedEnvironmentType = errors.New("undefined environment type")
)

type EditSettingsOptions struct {
	common.ItemOptions

	*sdk.EnvironmentEditSettings

	EditSettingsData
}

type EditSettingsData struct {
	Name string

	RemoteDevelopmentAllowed enum.Bool
	AutoUpdate               enum.Bool

	Labels       map[string]string
	LabelReplace bool

	K8SIntegration string

	EphemeralK8SIntegration   string
	CreateEphemeralOnPrCreate enum.Bool
	DestroyEphemeralOnPrClose enum.Bool
	AutoDeployEphemeral       enum.Bool
	TerminationProtection     enum.Bool

	DisableEphemeralBranchWhitelist bool
	EphemeralBranchWhitelistRegex   string
}

func NewEditSettingsOptions(environment string) *EditSettingsOptions {
	options := &EditSettingsOptions{
		ItemOptions: *common.NewItemOptions(environment),

		EditSettingsData: EditSettingsData{
			RemoteDevelopmentAllowed: enum.BoolNone,
			AutoUpdate:               enum.BoolNone,
		},

		EnvironmentEditSettings: sdk.NewEnvironmentEditSettingsWithDefaults(),
	}

	return options
}

func (eso *EditSettingsOptions) UpdateEditSettingsForType(environmentType string) error {
	switch environmentType {
	case "primary":
		primaryEdit := sdk.PrimaryAsEnvironmentEditSettingsEdit(sdk.NewPrimaryWithDefaults())
		eso.Edit = &primaryEdit
	case "ephemeral":
		ephemeralEdit := sdk.EphemeralAsEnvironmentEditSettingsEdit(sdk.NewEphemeralWithDefaults())
		eso.Edit = &ephemeralEdit
	default:
		return fmt.Errorf("%w: %s", errUnknownEnvironmentType, environmentType)
	}

	return nil
}

func (eso *EditSettingsOptions) UpdateCommandFlags(command *cobra.Command) {
	flags := command.Flags()

	data := &eso.EditSettingsData

	flags.StringVar(&data.Name, "name", data.Name, "Update environment name")

	autoUpdateFlag := enum.BoolFlag(
		&data.AutoUpdate,
		"auto-update",
		"Automatically update the environment when components git refs change",
	)
	flags.AddFlag(autoUpdateFlag)
	autoUpdateFlag.NoOptDefVal = "true"

	rdevFlag := enum.BoolFlag(
		&data.RemoteDevelopmentAllowed,
		"remote-development",
		"Allow remote development for the environment",
	)
	flags.AddFlag(rdevFlag)
	rdevFlag.NoOptDefVal = "true"

	flags.StringToStringVar(&data.Labels, "label", data.Labels, "Set labels for the environment (key=value)")
	flags.BoolVar(&data.LabelReplace, "label-replace", data.LabelReplace, "Set label strategy to replace (default: merge)")

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	ephCreateFlag := enum.BoolFlag(
		&data.CreateEphemeralOnPrCreate,
		"create-ephemeral-on-pr",
		"Create ephemeral environments when pull requests are created (for 'primary' environments)",
	)
	flags.AddFlag(ephCreateFlag)
	ephCreateFlag.NoOptDefVal = "true"

	ephDestroyFlag := enum.BoolFlag(
		&data.DestroyEphemeralOnPrClose,
		"destroy-ephemeral-on-pr-close",
		"Destroys the created ephemerals when the pull request is closed (or merged) (for 'primary' environments)",
	)
	flags.AddFlag(ephDestroyFlag)
	ephDestroyFlag.NoOptDefVal = "true"

	ephAutoDeployFlag := enum.BoolFlag(
		&data.AutoDeployEphemeral,
		"auto-deploy-ephemerals",
		"Auto deploy the created ephemerals (for 'primary' environments)",
	)
	flags.AddFlag(ephAutoDeployFlag)
	ephAutoDeployFlag.NoOptDefVal = "true"

	terminationProtectionFlag := enum.BoolFlag(
		&data.TerminationProtection,
		"termination-protection",
		"Prevent environment from being accidentally terminated (for 'primary' environments)",
	)
	flags.AddFlag(terminationProtectionFlag)
	terminationProtectionFlag.NoOptDefVal = "true"

	flags.BoolVar(&data.DisableEphemeralBranchWhitelist, "disable-ephemeral-branch-whitelist", data.DisableEphemeralBranchWhitelist, "Disable ephemerl branch whitelist")
	flags.StringVar(&data.EphemeralBranchWhitelistRegex, "ephemeral-branch-whitelist", data.EphemeralBranchWhitelistRegex, "Set ephemeral branch whitelist regex")
	command.MarkFlagsMutuallyExclusive("disable-ephemeral-branch-whitelist", "ephemeral-branch-whitelist")

	flags.StringVar(&data.EphemeralK8SIntegration, "ephemerals-k8s", data.EphemeralK8SIntegration, "The Kubernetes integration to be used for the ephemeral environments triggered by this environment (for 'primary' environments)")
}

func EditSettings(options *EditSettingsOptions) (*sdk.EnvironmentItem, error) {
	if options.Edit == nil {
		return nil, fmt.Errorf("%w", errUndefinedEnvironmentType)
	}

	model, resp, err := EditSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditSettingsRaw(options *EditSettingsOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentEditSettings(ctx, options.ID)

	return applyEditSettingsOptions(request, options).Execute()
}

func applyEditSettingsOptions(
	request sdk.ApiEnvironmentEditSettingsRequest,
	options *EditSettingsOptions,
) sdk.ApiEnvironmentEditSettingsRequest {
	if options.EditSettingsData.Name != "" {
		options.EnvironmentEditSettings.SetName(options.EditSettingsData.Name)
	}

	if options.EditSettingsData.K8SIntegration != "" {
		options.EnvironmentEditSettings.SetKubernetesIntegration(options.EditSettingsData.K8SIntegration)
	}

	if options.EditSettingsData.RemoteDevelopmentAllowed != enum.BoolNone {
		options.EnvironmentEditSettings.SetRemoteDevelopmentAllowed(options.EditSettingsData.RemoteDevelopmentAllowed == enum.BoolTrue)
	}

	if options.EditSettingsData.AutoUpdate != enum.BoolNone {
		options.EnvironmentEditSettings.SetAutoUpdate(options.EditSettingsData.AutoUpdate == enum.BoolTrue)
	}

	if options.EditSettingsData.Labels != nil {
		labelsEdit := *sdk.NewEdit()
		if options.EditSettingsData.LabelReplace {
			labelsEdit.SetStrategy("replace")
		}

		labelsEdit.SetValues(options.EditSettingsData.Labels)

		options.EnvironmentEditSettings.SetLabels(labelsEdit)
	} else if options.EditSettingsData.LabelReplace {
		labelsEdit := *sdk.NewEdit()
		labelsEdit.SetStrategy("replace")

		options.EnvironmentEditSettings.SetLabels(labelsEdit)
	}

	isPrimaryType := options.EnvironmentEditSettings.Edit.Primary != nil
	if isPrimaryType {
		applyPrimaryEditSettingsOptions(options)
	}

	request = request.EnvironmentEditSettings(*options.EnvironmentEditSettings)

	return request
}

func applyPrimaryEditSettingsOptions(options *EditSettingsOptions) {
	if options.EphemeralK8SIntegration != "" {
		options.EnvironmentEditSettings.Edit.Primary.SetEphemeralKubernetesIntegration(options.EphemeralK8SIntegration)
	}

	if options.CreateEphemeralOnPrCreate != enum.BoolNone {
		options.EnvironmentEditSettings.Edit.Primary.SetCreateEphemeralOnPrCreate(options.CreateEphemeralOnPrCreate == enum.BoolTrue)
	}

	if options.AutoDeployEphemeral != enum.BoolNone {
		options.EnvironmentEditSettings.Edit.Primary.SetAutoDeployEphemeral(options.AutoDeployEphemeral == enum.BoolTrue)
	}

	if options.TerminationProtection != enum.BoolNone {
		options.EnvironmentEditSettings.Edit.Primary.SetTerminationProtection(options.TerminationProtection == enum.BoolTrue)
	}

	if options.DestroyEphemeralOnPrClose != enum.BoolNone {
		options.EnvironmentEditSettings.Edit.Primary.SetDestroyEphemeralOnPrClose(options.DestroyEphemeralOnPrClose == enum.BoolTrue)
	}

	if options.DisableEphemeralBranchWhitelist {
		primaryOptions := sdk.NewPrimaryOptionsEdit()
		primaryOptions.SetWhitelistEnabled(false)

		options.EnvironmentEditSettings.Edit.Primary.SetPrimaryOptions(*primaryOptions)
	}

	if options.EphemeralBranchWhitelistRegex != "" {
		primaryOptions := sdk.NewPrimaryOptionsEdit()
		primaryOptions.SetWhitelistEnabled(true)
		primaryOptions.SetBranchWhitelist(options.EphemeralBranchWhitelistRegex)

		options.EnvironmentEditSettings.Edit.Primary.SetPrimaryOptions(*primaryOptions)
	}
}
