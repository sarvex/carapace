package cmd

import (
    "github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:     "action",
	Short:   "action example",
	Aliases: []string{"alias"},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringP("files", "f", "", "files flag")
	actionCmd.Flags().StringP("groups", "g", "", "groups flag")
	actionCmd.Flags().String("hosts", "", "hosts flag")
	actionCmd.Flags().StringP("message", "m", "", "message flag")
	actionCmd.Flags().StringP("net_interfaces", "n", "", "net_interfaces flag")
	actionCmd.Flags().StringP("options", "o", "", "options flag")
	actionCmd.Flags().String("path_files", "", "path_files flag")
	actionCmd.Flags().StringP("users", "u", "", "users flag")
	actionCmd.Flags().StringP("values", "v", "", "values flag")
	actionCmd.Flags().StringP("values_described", "d", "", "values with description flag")
	actionCmd.Flags().StringP("custom", "c", "", "custom flag")
	actionCmd.Flags().String("multi_parts", "", "multi_parts flag")

	carapace.Gen(actionCmd).FlagCompletion(carapace.ActionMap{
		"files":            carapace.ActionFiles(".go"),
		"groups":           carapace.ActionGroups(),
		"hosts":            carapace.ActionHosts(),
		"message":          carapace.ActionMessage("message example"),
		"net_interfaces":   carapace.ActionNetInterfaces(),
		"options":          carapace.ActionOptions(),
		"path_files":       carapace.ActionPathFiles(""),
		"users":            carapace.ActionUsers(),
		"values":           carapace.ActionValues("values", "example"),
		"values_described": carapace.ActionValuesDescribed("values", "valueDescription", "example", "exampleDescription"),
		"custom":           carapace.Action{Zsh: "_most_recent_file 2"},
		"multi_parts":      carapace.ActionMultiParts('/', "multi/parts", "multi/parts/example", "multi/parts/test", "example/parts"),
	})
}
