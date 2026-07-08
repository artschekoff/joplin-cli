package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flagDoc struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type commandDoc struct {
	Name    string    `json:"name"`
	Short   string    `json:"short"`
	Long    string    `json:"long,omitempty"`
	Flags   []flagDoc `json:"flags"`
	Output  string    `json:"output,omitempty"`
	Example string    `json:"example,omitempty"`
}

type interfaceDoc struct {
	Binary   string       `json:"binary"`
	Version  string       `json:"version"`
	Commands []commandDoc `json:"commands"`
}

// collect walks the tree, appending a doc for every runnable (leaf) command.
// prefix is the space-joined path of ancestor command names below the root.
func collect(cmd *cobra.Command, prefix string, out *[]commandDoc) {
	name := strings.TrimSpace(prefix + " " + cmd.Name())
	if cmd.Runnable() && cmd.Name() != "help" {
		cd := commandDoc{
			Name:    name,
			Short:   cmd.Short,
			Long:    cmd.Long,
			Output:  cmd.Annotations["output"],
			Example: cmd.Annotations["example"],
		}
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			cd.Flags = append(cd.Flags, flagDoc{
				Name:        f.Name,
				Shorthand:   f.Shorthand,
				Type:        f.Value.Type(),
				Default:     f.DefValue,
				Description: f.Usage,
			})
		})
		*out = append(*out, cd)
	}
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		collect(c, name, out)
	}
}

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe",
		Short: "Emit the full CLI interface as JSON for machine consumption",
		Long: `Prints a JSON document describing every command: fully-qualified name,
short/long help, input flags (name, shorthand, type, default, description),
the output JSON shape, and an example invocation. Consumed by LLM agents that
need a schema of the CLI without parsing text help.

Output: JSON object with "binary", "version", "commands" fields.
Exit codes: 0 always (unless the command tree itself is malformed).`,
		Annotations: map[string]string{
			"output":  `{"binary":"joplin-cli","version":"string","commands":[{"name":"note search","short":"string","long":"string","flags":[{"name":"limit","type":"int","default":"100","description":"string"}],"output":"json shape string","example":"string"}]}`,
			"example": `joplin-cli describe | jq '.commands[].name'`,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			doc := interfaceDoc{Binary: root.Name(), Version: root.Version}
			for _, c := range root.Commands() {
				if c.Hidden {
					continue
				}
				collect(c, "", &doc.Commands)
			}
			out, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}
}
