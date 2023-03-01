package git

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

var prepareTemplate = fmt.Sprintf(`In order to clone repositories, run the following:
{{ range $repo := .PrepareManager.Repositories }}
    {{ if $.Options.DiscardHistory -}}
        git clone \
            --separate-git-dir=$(mktemp -u) \
            --depth 1 \
            {{ $repo }} {{ $.PrepareManager.GetDir $repo }} \
        && rm -f {{ $.PrepareManager.GetDir $repo }}/.git
    {{- else -}}
		git clone {{ $repo }} {{ $.PrepareManager.GetDir $repo }}
	{{- end }}
{{- end }}

For Remote Development:

    Component%[1]s Command
{{- range .Components }}
    {{ .Name }}%[1]s %[2]s rdev up --component {{ .Id }} --local-sync-path {{ $.PrepareManager.GetPath . }}
{{- end }}
`, "\t", build.Name)

type PrepareData struct {
	PrepareManager *PrepareManager

	Components []sdk.ComponentGitCollection

	Options PrepareOptions
}

type PrepareOptions struct {
	DiscardHistory bool
}

func NewPrepareOptions() *PrepareOptions {
	return &PrepareOptions{
		DiscardHistory: false,
	}
}

func (po *PrepareOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.BoolVar(&po.DiscardHistory, "discard-history", po.DiscardHistory, "Discard git history after clone")
}

func PrintPrepareInfo(components []sdk.ComponentGitCollection, options *PrepareOptions) error {
	if len(components) == 0 {
		return nil
	}

	prepareManager := NewPrepareManager()
	if err := prepareManager.AddComponents(components); err != nil {
		return err
	}

	tpl, err := template.New("git prepare").Parse(prepareTemplate)
	if err != nil {
		return err
	}

	prepareInfo := PrepareData{
		PrepareManager: prepareManager,
		Components:     components,

		Options: *options,
	}

	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
	if err = tpl.Execute(writer, prepareInfo); err != nil {
		return err
	}

	writer.Flush()

	return nil
}
