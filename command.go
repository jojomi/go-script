package script

import (
	"strings"
)

type Command interface {
	Binary() string
	Args() []string
	Add(input string)
	AddAll(input ...string)
	String() string
}

func isWrapped(source, s string) bool {
	return strings.HasPrefix(source, s) && strings.HasSuffix(source, s)
}

// SplitCommand helper splits a string to command and arbitrarily many args.
// Does handle bash-like escaping (\) and string delimiters " and '.
func SplitCommand(input string) (command string, args []string) {
	quotes := []string{`"`, `'`}

	var (
		ok     bool
		length int
		value  string
		index  = 0
	)
	args = make([]string, 0)

outerloop:
	for {
		if index >= len(input) {
			break
		}

		ok, length, _ = parseWhitespace(input[index:])
		if ok {
			index += length
			continue
		}

		for _, quote := range quotes {
			ok, length, value = parseQuoted(input[index:], quote, `\`+quote)
			if ok {
				if command == "" {
					command = value
				} else {
					args = append(args, value)
				}
				index += length
				continue outerloop
			}
		}

		ok, length, value = parseUnquoted(input[index:])
		if ok {
			if command == "" {
				command = value
			} else {
				args = append(args, value)
			}
			index += length
			continue
		}
	}
	return
}

func parseQuoted(input, quoteString, escapeString string) (ok bool, length int, value string) {
	if !strings.HasPrefix(input, quoteString) {
		return
	}

	length = len(quoteString)
	for {
		if length >= len(input) {
			break
		}
		// escaped quoteString? (continue!)
		if strings.HasPrefix(input[length:], escapeString) {
			length += len(escapeString)
			value += quoteString
		}
		// quoteString (end!)
		if strings.HasPrefix(input[length:], quoteString) {
			length += len(quoteString)
			ok = true
			return
		}

		// otherwise inner content
		value += input[length : length+1]
		length++
	}

	return ok, length, value
}

func parseUnquoted(input string) (ok bool, length int, value string) {
	length = 0
	for {
		if length >= len(input) {
			ok = true
			return
		}
		// whitespace (end!) // TODO all whitespace!
		if strings.HasPrefix(input[length:], " ") {
			length++
			ok = true
			return
		}

		// otherwise inner content
		value += input[length : length+1]
		length++
	}
}

func parseWhitespace(input string) (ok bool, length int, value string) {
	length = 0
	for {
		if length >= len(input) {
			break
		}
		// no whitespace (end!) // TODO all whitespace!
		if !strings.HasPrefix(input[length:], " ") {
			ok = length > 0
			return
		}

		// otherwise inner content (whitespace)
		value += input[length : length+1]
		length++
	}

	return ok, length, value
}
