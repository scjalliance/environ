package main

// Prefix is the prefix used for environment variables.
const Prefix = "env."

// Environment is a map of evironment variable names to environment values
type Environment map[string]string

// Parse will scan the provided argument set for arguments matching one of
// these forms:
//
//   env.NAME=VALUE
//   env.NAME=
//
// Any arguments not matching these patterns are returned in remaining.
func (e Environment) Parse(arguments []string) (remaining []string) {
	for _, arg := range arguments {
		if !e.ParseArg(arg) {
			remaining = append(remaining, arg)
		}
	}

	return
}

// ParseArg evaluates a single argument and determines whether it describes an
// environment variable. If it does then the variable is added to the
// environment and ParseArg returns true.
func (e Environment) ParseArg(arg string) bool {
	const (
		pre = len(Prefix)
		min = pre + 2
	)

	if len(arg) < min {
		return false
	}
	if arg[0:pre] != Prefix {
		return false
	}
	name := arg[pre:] // Leaves name with at least 2 characters
	if name[0] == '\x00' || name[0] == '=' {
		return false
	}

	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // Equals cannot be first
		if name[i] == '\x00' {
			// NUL is not permitted in environment variable names
			return false
		}
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	if !hasValue {
		return false
	}

	e[name] = value

	return true
}
