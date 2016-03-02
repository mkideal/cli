package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// InstallShell install bash_completion
func genBashCompletion(root *Command) (*bytes.Buffer, error) {
	buff := bytes.NewBufferString("")
	t, err := template.New("bash_completion").Parse(shellTemplateText)
	if err != nil {
		return nil, err
	}
	return buff, t.Execute(buff, struct {
		Cli        string
		CompleteFn string
	}{Cli: root.Name, CompleteFn: "#TODO"})
}

func InstallBashCompletion(root *Command) error {
	if root.Name == "" {
		return fmt.Errorf("root command name is empty")
	}
	filename := filepath.Join(os.Getenv("HOME"), "."+root.Name+"_compeltion")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	buff, err := genBashCompletion(root)
	if err != nil {
		return err
	}
	if _, err := file.Write(buff.Bytes()); err != nil {
		return err
	}

	//TODO: write .bashrc or .bash_profile or others

	return nil
}

const shellTemplateText = `# {{.Cli}} command completion script

COMP_WORDBREAKS=${COMP_WORDBREAKS/=/}
COMP_WORDBREAKS=${COMP_WORDBREAKS/@/}
export COMP_WORDBREAKS

__complete_fn() {
#COMP_CWORD
#COMP_LINE
#COMP_POINT
#COMP_WORDS
#{{.Cli}} completion -- "${COMP_WORDS[@]}"
{{.CompleteFn}}
}

if type complete &>/dev/null; then
  _{{.Cli}}_completion () {
    local si="$IFS"
    IFS=$'\n' COMPREPLY=($(__complete_fn \
                           2>/dev/null)) || return $?
    IFS="$si"
  }
  complete -F _{{.Cli}}_completion {{.Cli}}
elif type compdef &>/dev/null; then
  _{{.Cli}}_completion() {
    si=$IFS
    compadd -- $(COMP_CWORD=$((CURRENT-1)) \
                 COMP_LINE=$BUFFER \
                 COMP_POINT=0 \
                 COMP_WORDS="${words[@]}" \
				 __complete_fn \
                 2>/dev/null)
    IFS=$si
  }
  compdef _{{.Cli}}_completion {{.Cli}}
elif type compctl &>/dev/null; then
  _{{.Cli}}_completion () {
    local cword line point words si
    read -Ac words
    read -cn cword
    let cword-=1
    read -l line
    read -ln point
    si="$IFS"
    IFS=$'\n' reply=($(COMP_CWORD="$cword" \
                       COMP_LINE="$line" \
                       COMP_POINT="$point" \
                       COMP_WORDS="${words[@]}" \
					   __complete_fn \
                       2>/dev/null)) || return $?
    IFS="$si"
  }
  compctl -K _{{.Cli}}_completion {{.Cli}}
fi
`
