package environment

import (
	"errors"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"

	"github.com/spf13/cobra"
)

type CreateOptions struct {
	DeployOptions

	sdk.EnvironmentCreateAction

	genesisSourceOptions GenesisSourceOptions

	WithDeploy bool

	EphemeralBranchWhitelist string
}

func NewCreateOptions() *CreateOptions {
	environmentCreateAction := sdk.NewEnvironmentCreateActionWithDefaults()
	environmentCreateAction.SetKubernetesIntegration("")
	environmentCreateAction.SetEphemeralKubernetesIntegration("")
	environmentCreateAction.SetLabels(map[string]string{})
	environmentCreateAction.SetAutoDeployEphemeral(false)
	environmentCreateAction.SetTerminationProtection(false)
	environmentCreateAction.SetCreateEphemeralOnPrCreate(false)
	environmentCreateAction.SetDestroyEphemeralOnPrClose(false)

	return &CreateOptions{
		DeployOptions: *NewDeployOptions(""),

		EnvironmentCreateAction: *environmentCreateAction,

		genesisSourceOptions: *NewGenesisSourceOptions(),
	}
}

func (co *CreateOptions) UpdateCommandFlags(command *cobra.Command) {
	flags := command.Flags()

	k8sIntegration := co.KubernetesIntegration.Get()

	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the environment")
	flags.BoolVar(&co.WithDeploy, "deploy", co.WithDeploy, "Deploy the environment after creation")
	flags.StringVar(k8sIntegration, "k8s", *k8sIntegration, "Use a Kubernetes integration for the environment")

	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the project for the new environment")

	flags.StringToStringVar(co.Labels, "label", *co.Labels, "Set labels for the new environment (key=value)")

	ephemeralsK8sIntegration := co.EphemeralKubernetesIntegration.Get()
	flags.BoolVar(co.CreateEphemeralOnPrCreate, "create-ephemeral-on-pr", *co.CreateEphemeralOnPrCreate, "Create ephemeral environments when pull requests are created")
	flags.BoolVar(co.DestroyEphemeralOnPrClose, "destroy-ephemeral-on-pr-close", *co.DestroyEphemeralOnPrClose, "Destroys the created ephemerals when the pull request is closed (or merged)")
	flags.BoolVar(co.AutoDeployEphemeral, "auto-deploy-ephemerals", *co.AutoDeployEphemeral, "Auto deploy the created ephemerals")
	flags.BoolVar(co.TerminationProtection, "termination-protection", *co.TerminationProtection, "Prevent environment from being accidentally terminated")
	flags.StringVar(ephemeralsK8sIntegration, "ephemerals-k8s", *ephemeralsK8sIntegration, "The Kubernetes integration to be used for the ephemeral environments triggered by this environment")

	flags.StringVar(&co.EphemeralBranchWhitelist, "ephemeral-branch-whitelist", co.EphemeralBranchWhitelist, "Ephemeral branch whitelist regex when create-ephemeral-on-pr is set")

	co.DeployOptions.UpdateFlagSet(flags)

	co.genesisSourceOptions.updateCommandFlags(command, "creation")
}

func (co *CreateOptions) Validate() error {
	return co.genesisSourceOptions.validate()
}

func (co *CreateOptions) AttachGenesis() error {
	fromGit, fromGitSpec, fromString, fromTemplate, err := co.genesisSourceOptions.getGenesis()
	if err != nil {
		return err
	}

	co.Genesis = &sdk.EnvironmentCreateActionGenesis{
		FromGit:      fromGit,
		FromGitSpec:  fromGitSpec,
		FromTemplate: fromTemplate,
		FromString:   fromString,
	}

	return nil
}

func (co *CreateOptions) HandleError(cmd *cobra.Command, err error) error {
	var apiError api.Error

	if errors.As(err, &apiError) {
		return co.genesisSourceOptions.handleError(cmd, apiError)
	}

	return lib.FormatCommandError(cmd, err)
}

func Create(options *CreateOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentCreate(ctx)

	if len(*options.Labels) > 0 {
		options.EnvironmentCreateAction.SetLabels(*options.Labels)
	}

	if options.EphemeralBranchWhitelist != "" {
		primaryOptions := sdk.NewPrimaryOptionsCreate()
		primaryOptions.SetWhitelistEnabled(true)
		primaryOptions.SetBranchWhitelist(options.EphemeralBranchWhitelist)

		options.SetPrimaryOptions(*primaryOptions)
	}

	request = request.EnvironmentCreateAction(options.EnvironmentCreateAction)

	return request.Execute()
}
