package port_forward

import (
	"os"

	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		resourcePath string
		podName      string
	)

	command := &cobra.Command{
		Use:     "start mappings...",
		Short:   "Starts the port forwarding for the given mappings.",
		Example: "start 8080:80 3306 :9003",

		ValidArgsFunction: cobra.NoFileCompletions,

		Run: func(cmd *cobra.Command, portMappings []string) {
			root := cmd.Root()
			root.SetArgs(append([]string{
				"components", "port-forward",
				"--id", settings.Profile.Context.ServiceComponent,
				"--resource", resourcePath,
				"--pod", podName,
			}, portMappings...))

			if err := root.Execute(); err != nil {
				os.Exit(1)
			}
		},
	}

	flags := command.Flags()

	componentFlag := options.ServiceComponent.AddFlag("component", "Service Component")
	flags.AddFlag(componentFlag)
	_ = command.MarkFlagRequired(componentFlag.Name)

	flags.StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	flags.StringVar(&podName, "pod", "", "The resource pod to forward ports to.")

	mainCmd.AddCommand(command)
}
