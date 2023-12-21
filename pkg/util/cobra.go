package util

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type BoolFlagType string

const (
	HelpTemplate = "bns_annotation_help_template"
)

const (
	FlagHidden     BoolFlagType = "bns_annotation_hidden"
	FlagRequired   BoolFlagType = cobra.BashCompOneRequiredFlag
	FlagDirname    BoolFlagType = cobra.BashCompSubdirsInDir
	FlagAllowBlank BoolFlagType = "bns_annotation_allow_blank"
)

func HasHelp(flag *pflag.Flag) bool {
	if flag.Annotations == nil {
		return false
	}

	return len(flag.Annotations[HelpTemplate]) > 0
}

func GetHelp(flag *pflag.Flag) string {
	if !HasHelp(flag) {
		return flag.Usage
	}

	return strings.Join(flag.Annotations[HelpTemplate], ", ")
}

func IsHidden(flag *pflag.Flag) bool {
	if flag.Annotations == nil {
		return false
	}

	if flag.Annotations[string(FlagHidden)] == nil {
		return false
	}

	return flag.Annotations[string(FlagHidden)][0] == StrTrue
}

func MarkFlag(flag *pflag.Flag, flagTypes ...BoolFlagType) *pflag.Flag {
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}

	for _, flagType := range flagTypes {
		flag.Annotations[string(flagType)] = []string{StrTrue}
	}

	return flag
}

func MarkFlagRequiredWithHelp(flag *pflag.Flag, helpTemplate string) *pflag.Flag {
	AppendFlagHelp(flag, helpTemplate)

	MarkFlag(flag, FlagRequired)

	return flag
}

func AppendFlagHelp(flag *pflag.Flag, helpTemplate string) *pflag.Flag {
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}

	if flag.Annotations[HelpTemplate] == nil {
		flag.Annotations[HelpTemplate] = []string{}
	}

	flag.Annotations[HelpTemplate] = append(flag.Annotations[HelpTemplate], helpTemplate)

	return flag
}

func GetFlagBoolAnnotation(flag *pflag.Flag, annotation BoolFlagType) bool {
	if flag.Annotations == nil {
		return false
	}

	if flag.Annotations[string(annotation)] == nil {
		return false
	}

	return flag.Annotations[string(annotation)][0] == StrTrue
}
