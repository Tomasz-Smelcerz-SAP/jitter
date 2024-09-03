package cmd

import "strings"

const (
	argPrefix         = "--"
	argValueSeparator = "="
)

type Arguments struct {
	defs []string
}

func (a *Arguments) Add(cliarg string) {
	a.defs = append(a.defs, cliarg)
}

func (a *Arguments) Get(name string) (string, bool) {
	if strings.HasPrefix(name, argPrefix) {
		name = name[2:]
	}

	for i, arg := range a.defs {
		if strings.Contains(arg, name) {
			return a.value(i), true
		}
	}
	return "", false
}

func (a *Arguments) value(idx int) string {
	val := a.defs[idx]
	if strings.Contains(val, argValueSeparator) {
		val = strings.Split(val, argValueSeparator)[1]
	}
	return val
}
