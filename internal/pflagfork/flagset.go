package pflagfork

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
}

func (f FlagSet) IsPosix() bool {
	if method := reflect.ValueOf(f.FlagSet).MethodByName("IsPosix"); method.IsValid() {
		if values := method.Call([]reflect.Value{}); len(values) == 1 && values[0].Kind() == reflect.Bool {
			return values[0].Bool()
		}
	}
	return true
}

func (f FlagSet) IsShorthandSeries(arg string) bool {
	if len(arg) < 3 || !strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") || !f.IsPosix() {
		return false
	}

	flag := f.ShorthandLookup(string(arg[1]))
	if flag == nil {
		return false
	}

	switch {
	case (Flag{flag}).IsOptarg() && arg[2] == byte(flag.OptargDelimiter):
		return false

	case (Flag{flag}).TakesValue():
		return false

	default:
		return true
	}
}

func (f FlagSet) IsMutuallyExclusive(flag *pflag.Flag) bool {
	if groups, ok := flag.Annotations["cobra_annotation_mutually_exclusive"]; ok {
		for _, group := range groups {
			for _, name := range strings.Split(group, " ") {
				if other := f.Lookup(name); other != nil && other.Changed {
					return true
				}
			}
		}
	}
	return false
}

func (f *FlagSet) VisitAll(fn func(*Flag)) {
	f.FlagSet.VisitAll(func(flag *pflag.Flag) {
		fn(&Flag{flag})
	})

}

func (fs FlagSet) shorthandLookupArg(arg string) (result *Flag) {
	for i := 1; i < len(arg); i++ {
		f := fs.ShorthandLookup(string(arg[i]))
		if f == nil {
			return
		}
		result = &Flag{f}
		switch {
		case result.IsOptarg() && i < len(arg)-1 && arg[i] == byte(result.OptargDelimiter()):
			return
		case result.TakesValue():
			return
		}
	}
	return
}

func (fs FlagSet) LookupArg(arg string) (result *Flag) {
	switch {
	case fs.IsShorthandSeries(arg):
		return fs.shorthandLookupArg(arg)

	default:
		fs.VisitAll(func(f *Flag) {
			if result != nil {
				return
			}

			if f.NameMatches(arg) {
				result = f
			}
		})
	}
	return
}
