package action

import (
	"errors"
	"strings"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	k8sExec "bunnyshell.com/cli/pkg/k8s/kubectl/exec"
	k8sWizard "bunnyshell.com/cli/pkg/wizard/k8s"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var errInvalidPodName = errors.New("invalid pod name")

type SSHOptions struct {
	Namespace string
	PodName   string
	Container string

	Shell string

	NoTTY    bool
	NoBanner bool
}

func (o *SSHOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&o.PodName, "pod", o.PodName, "Pod name in namespace/pod-name format")
	flags.StringVar(&o.Container, "container", o.Container, "Container name")
	flags.StringVar(&o.Shell, "shell", o.Shell, "Shell to use")

	flags.BoolVar(&o.NoTTY, "no-tty", o.NoTTY, "Do not allocate a TTY")
	flags.BoolVar(&o.NoBanner, "no-banner", o.NoBanner, "Do not show environment banner before ssh")
}

func (o *SSHOptions) MakeExecOptions(kubeConfig *environment.KubeConfigItem) *k8sExec.Options {
	return &k8sExec.Options{
		TTY:   !o.NoTTY,
		Stdin: true,

		Command: []string{o.Shell},

		KubeConfig: kubeConfig.Bytes,
	}
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	sshOptions := SSHOptions{
		Shell: "/bin/sh",
	}

	command := &cobra.Command{
		Use: "ssh",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			itemOptions := component.NewItemOptions(settings.Profile.Context.ServiceComponent)
			componentItem, err := component.Get(itemOptions)
			if err != nil {
				return err
			}

			kubeConfigOptions := environment.NewKubeConfigOptions(componentItem.GetEnvironment())
			kubeConfig, err := environment.KubeConfig(kubeConfigOptions)
			if err != nil {
				return err
			}

			execCommand, err := k8sExec.Exec(sshOptions.MakeExecOptions(kubeConfig))
			if err != nil {
				return err
			}

			podListOptions := &k8sWizard.PodListOptions{
				Component: componentItem.GetId(),
			}
			containerListOptions := &k8sWizard.ContainerListOptions{
				Client: execCommand.PodClient,
			}

			if err = ensureContainerSelected(&sshOptions, podListOptions, containerListOptions); err != nil {
				return err
			}

			execCommand.Namespace = sshOptions.Namespace
			execCommand.PodName = sshOptions.PodName
			execCommand.ContainerName = sshOptions.Container

			if err = execCommand.Validate(); err != nil {
				return err
			}

			if !sshOptions.NoBanner {
				showBanner(cmd, sshOptions, componentItem.GetId())
			}

			return execCommand.Run()
		},
	}

	flags := command.Flags()

	idFlag := options.ServiceComponent.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	sshOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func ensureContainerSelected(
	sshOptions *SSHOptions,
	podListOptions *k8sWizard.PodListOptions,
	containerListOptions *k8sWizard.ContainerListOptions,
) error {
	if sshOptions.PodName != "" {
		parts := strings.Split(sshOptions.PodName, "/")

		if len(parts) != 2 {
			return errInvalidPodName
		}

		sshOptions.Namespace = parts[0]
		sshOptions.PodName = parts[1]
	} else {
		resource, err := k8sWizard.PodSelect(podListOptions)
		if err != nil {
			return err
		}

		sshOptions.PodName = resource.GetName()
		sshOptions.Namespace = resource.GetNamespace()
	}

	if sshOptions.Container == "" {
		containerListOptions.Namespace = sshOptions.Namespace
		containerListOptions.PodName = sshOptions.PodName

		container, err := k8sWizard.ContainerSelect(containerListOptions)
		if err != nil {
			return err
		}

		sshOptions.Container = container.Name
	}

	return nil
}

func showBanner(cmd *cobra.Command, sshOptions SSHOptions, environmentID string) {
	bnsColor := color.New(color.BgBlue).Add(color.FgWhite).Add(color.Bold)

	cmd.Printf(
		"Environment: %s Pod: %s Container: %s\n",
		bnsColor.Sprint(environmentID),
		bnsColor.Sprint(sshOptions.Namespace)+"/"+getPodName(bnsColor, sshOptions),
		bnsColor.Sprint(sshOptions.Container),
	)
}

func getPodName(bnsColor *color.Color, sshOptions SSHOptions) string {
	parts := strings.Split(sshOptions.PodName, "-")

	for index, part := range parts {
		parts[index] = bnsColor.Sprint(part)
	}

	return strings.Join(parts, "-")
}
