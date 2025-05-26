package component_debug

import (
	"github.com/spf13/pflag"
)

type SSHOptions struct {
	Shell string

	NoTTY    bool
	NoBanner bool

	OverrideClusterServer string
}

func (o *SSHOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&o.Shell, "shell", o.Shell, "Shell to use")

	flags.BoolVar(&o.NoTTY, "no-tty", o.NoTTY, "Do not allocate a TTY")
	flags.BoolVar(&o.NoBanner, "no-banner", o.NoBanner, "Do not show environment banner before ssh")
}
