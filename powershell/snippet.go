package powershell

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	buf := new(bytes.Buffer)

	var subCommandCases bytes.Buffer
	generatePowerShellSubcommandCases(&subCommandCases, cmd, "")
	fmt.Fprintf(buf, powerShellCompletionTemplate, cmd.Name(), cmd.Name(), subCommandCases.String())

	return buf.String()
}

var powerShellCompletionTemplate = `using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Register-ArgumentCompleter -Native -CommandName '%s' -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements
    $command = @(
        '%s'
        for ($i = 1; $i -lt $commandElements.Count; $i++) {
            $element = $commandElements[$i]
            if ($element -isnot [StringConstantExpressionAst] -or
                $element.StringConstantType -ne [StringConstantType]::BareWord -or
                $element.Value.StartsWith('-')) {
                break
            }
            $element.Value
        }
    ) -join ';'
    $completions = @(switch ($command) {%s
    })
    $completions.Where{ $_.CompletionText -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}`

func nonCompletableFlag(flag *pflag.Flag) bool {
	return flag.Hidden || len(flag.Deprecated) > 0
}

func generatePowerShellSubcommandCases(out io.Writer, cmd *cobra.Command, previousCommandName string) {
	var cmdName string
	if previousCommandName == "" {
		cmdName = cmd.Name()
	} else {
		cmdName = fmt.Sprintf("%s;%s", previousCommandName, cmd.Name())
	}

	fmt.Fprintf(out, "\n        '%s' {", cmdName)

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if nonCompletableFlag(flag) {
			return
		}
		usage := escapeStringForPowerShell(flag.Usage)
		if len(flag.Shorthand) > 0 {
			fmt.Fprintf(out, "\n            [CompletionResult]::new('-%s', '%s', [CompletionResultType]::ParameterName, '%s')", flag.Shorthand, flag.Shorthand, usage)
		}
		fmt.Fprintf(out, "\n            [CompletionResult]::new('--%s', '%s', [CompletionResultType]::ParameterName, '%s')", flag.Name, flag.Name, usage)
	})

	for _, subCmd := range cmd.Commands() {
		usage := escapeStringForPowerShell(subCmd.Short)
		fmt.Fprintf(out, "\n            [CompletionResult]::new('%s', '%s', [CompletionResultType]::ParameterValue, '%s')", subCmd.Name(), subCmd.Name(), usage)
	}

	fmt.Fprint(out, "\n            break\n        }")

	for _, subCmd := range cmd.Commands() {
		generatePowerShellSubcommandCases(out, subCmd, cmdName)
	}
}

func escapeStringForPowerShell(s string) string {
	return strings.Replace(s, "'", "''", -1)
}
