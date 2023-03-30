package option

import (
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/pflag"
)

type Group struct {
	main  *pflag.Flag
	flags []*pflag.Flag
}

func (group *Group) GetFlags() []*pflag.Flag {
	return group.flags
}

func (group *Group) GetMainFlag() *pflag.Flag {
	return group.main
}

func (group *Group) AddFlag(flag *pflag.Flag) {
	group.flags = append(group.flags, flag)

	if group.main == nil {
		group.main = flag
	} else {
		copyFlagExtra(group.main, flag)
	}
}

func (group *Group) GetFlag(name string) *pflag.Flag {
	for _, flag := range group.flags {
		if flag.Name == name {
			return flag
		}
	}

	return nil
}

func (group *Group) IsChanged() bool {
	for _, flag := range group.flags {
		if flag.Changed {
			return true
		}
	}

	return false
}

func (group *Group) updateFlags(value string) error {
	for _, flag := range group.flags {
		flag.Changed = true

		err := flag.Value.Set(value)
		if err != nil {
			return err
		}

		if required, ok := flag.Annotations[string(util.FlagRequired)]; ok {
			required[0] = util.StrFalse
		}
	}

	return nil
}

func copyFlagExtra(src *pflag.Flag, dst *pflag.Flag) {
	if src.Annotations != nil {
		dst.Annotations = map[string][]string{}

		for k, v := range src.Annotations {
			dst.Annotations[k] = v
		}
	}
}
