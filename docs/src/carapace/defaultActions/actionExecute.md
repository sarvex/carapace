# ActionExecute

[`ActionExecute`] invokes an internal command and parses it's output using [`ActionImport`].

> Cobra commands can only be executed **once** so be sure each invocation uses a new instance.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	cmd := &cobra.Command{
		Use: "embedded",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.Flags().Bool("embedded-flag", false, "embedded flag")

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionValues("embeddedPositional1", "embeddedP1"),
		carapace.ActionValues("embeddedPositional2", "embeddedP2"),
	)

	return carapace.ActionExecute(cmd)
})
````

![](./actionExecute.cast)

[`ActionExecute`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionExecute
[`ActionImport`]:../defaultActions/actionImport.md