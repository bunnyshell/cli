package action

import (
	"errors"
	"fmt"
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

	OverrideClusterServer string
}

func (o *SSHOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&o.PodName, "pod", o.PodName, "Pod name in namespace/pod-name format")
	flags.StringVar(&o.Container, "container", o.Container, "Container name")
	flags.StringVar(&o.Shell, "shell", o.Shell, "Shell to use")

	flags.BoolVar(&o.NoTTY, "no-tty", o.NoTTY, "Do not allocate a TTY")
	flags.BoolVar(&o.NoBanner, "no-banner", o.NoBanner, "Do not show environment banner before ssh")

	flags.StringVar(&o.OverrideClusterServer, "override-kubeconfig-cluster-server", o.OverrideClusterServer, "Override kubeconfig cluster server with :port, host:port or scheme://host:port")
}

func (o *SSHOptions) MakeExecOptions(kubeConfig *environment.KubeConfigItem) *k8sExec.Options {
	motd := "/opt/bunnyshell/motd.txt"

	return &k8sExec.Options{
		TTY:   !o.NoTTY,
		Stdin: true,

		Command: []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf(
				"[ -f %[1]s ] && cat %[1]s; %[2]s",
				motd,
				o.Shell,
			),
		},

		KubeConfig: kubeConfig.Bytes,
	}
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	sshOptions := SSHOptions{
		Shell:                 "/bin/sh",
		OverrideClusterServer: "",
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

			kubeConfigOptions := environment.NewKubeConfigOptions(componentItem.GetEnvironment(), sshOptions.OverrideClusterServer)
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
				showBanner(cmd, sshOptions, componentItem.GetEnvironment())
			}

			return execCommand.Run()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.ServiceComponent.GetRequiredFlag("id"))

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
