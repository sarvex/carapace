package elvish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	vals := make([]string, len(values))
	for index, val := range values {
		if val.Description == "" {
			vals[index] = fmt.Sprintf(`edit:complex-candidate '%v' &display='%v'`, sanitizer.Replace(val.Value), sanitizer.Replace(val.Display))
		} else {
			vals[index] = fmt.Sprintf(`edit:complex-candidate '%v' &display='%v (%v)'`, sanitizer.Replace(val.Value), sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
		}
	}
	return strings.Join(vals, "\n")
}
