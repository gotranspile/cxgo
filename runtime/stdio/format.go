package stdio

import "strings"

type formatWord struct {
	Str  string
	Verb bool
}

func parseFormat(format string) []formatWord {
	var (
		words       []formatWord
		placeholder = false
		cur         []rune
	)
	reset := func() {
		cur = cur[:0]
		placeholder = false
	}
	pushWord := func(verb bool) {
		if len(cur) != 0 {
			words = append(words, formatWord{Str: string(cur), Verb: verb})
		}
		reset()
	}
	for _, b := range format {
		if b == '%' {
			if placeholder {
				if len(cur) == 1 { // %%
					cur = append(cur, '%')
					pushWord(true) // fake verb
					continue
				}
				// end the current placeholder, start a new one
				pushWord(true)
			}
			// start a new placeholder
			pushWord(false)
			placeholder = true
			cur = append(cur, '%')
			continue
		}
		if !placeholder {
			// continue the string
			cur = append(cur, b)
			continue
		}
		// in placeholder
		switch b {
		default:
			if strings.IndexRune("1234567890#+-*. l", b) >= 0 {
				// consider a part of the placeholder
				cur = append(cur, b)
				continue
			}
			// other rune: stop the placeholder
			pushWord(true)
			cur = append(cur, b)
			continue
		case 'i': // signed, allows octal and hex
		case 'd': // signed, only decimal
		case 'u': // unsigned, only decimal
		case 'o': // octal
		case 'x': // hex
		case 'f', 'e', 'g', 'F', 'E', 'G': // float
		case 'c': // char
		case 'p': // ptr
		case 's': // string
		case 'S': // wstring
		}
		cur = append(cur, b)
		pushWord(true)
	}
	pushWord(placeholder)
	return words
}
