package nushell

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

type record struct {
	Value        string `json:"value"`
	Description  string `json:"description,omitempty"`
	nushellStyle `json:",inline"`
}

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
)

func sanitize(values []common.RawValue) []common.RawValue {
	for index, v := range values {
		(&values[index]).Value = sanitizer.Replace(v.Value)
		(&values[index]).Display = sanitizer.Replace(v.Display)
		(&values[index]).Description = sanitizer.Replace(v.Description)
	}
	return values
}

// ActionRawValues formats values for nushell
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	filtered := values.FilterPrefix(currentWord)
	sort.Sort(common.ByDisplay(filtered))

	vals := make([]record, len(filtered))
	for index, val := range sanitize(filtered) {
		if strings.ContainsAny(val.Value, ` {}()[]<>$&"|;#\`+"`") {
			val.Value = fmt.Sprintf("'%v'", val.Value)
		}

		if !nospace {
			val.Value = val.Value + " "
		}

		vals[index] = record{Value: val.Value, Description: val.TrimmedDescription(), nushellStyle: toNushellstyle(val.Style)}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}

type nushellStyle struct {
	Foreground string   `json:"foreground,omitempty"`
	Background string   `json:"background,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}

func toNushellstyle(_style string) nushellStyle {
	s := parseStyle(_style)

	attributes := make([]string, 0)
	if s.Blink {
		attributes = append(attributes, "slow_blink")
	}
	if s.Bold {
		attributes = append(attributes, "bold")
	}
	if s.Dim {
		attributes = append(attributes, "dim")
	}
	if s.Italic {
		attributes = append(attributes, "italic")
	}
	if s.Inverse {
		attributes = append(attributes, "negative")
	}
	if s.Underlined {
		attributes = append(attributes, "underlined")
	}

	replacer := strings.NewReplacer(
		"bright-black", "dark_grey",
		"black", "black",
		"bright-red", "red",
		"red", "dark_red",
		"bright-green", "green",
		"green", "dark_green",
		"bright-yellow", "yellow",
		"yellow", "dark_yellow",
		"bright-blue", "blue",
		"blue", "dark_blue",
		"bright-magenta", "magenta",
		"magenta", "dark_magenta",
		"bright-cyan", "cyan",
		"cyan", "dark_cyan",
		"bright-white", "white",
		"white", "grey",
	)

	// TODO colorx -> ansi_(x)
	// TODO #rrggbb -> rgb_(r,g,b)

	fg := ""
	if s.Foreground != nil {
		fg = replacer.Replace(s.Foreground.String())
	}

	bg := ""
	if s.Background != nil {
		bg = replacer.Replace(s.Background.String())
	}

	// TODO patch color names
	return nushellStyle{
		Foreground: fg,
		Background: bg,
		Attributes: attributes,
	}
}

// TODO copied from style package
func parseStyle(s string) ui.Style {
	stylings := make([]ui.Styling, 0)
	for _, word := range strings.Split(s, " ") {
		if styling := ui.ParseStyling(word); styling != nil {
			stylings = append(stylings, styling)
		}
	}
	return ui.ApplyStyling(ui.Style{}, stylings...)
}
