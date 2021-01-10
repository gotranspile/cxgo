package stdio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name   string
		format string
		words  []formatWord
	}{
		{
			name:   "plain",
			format: `some string here`,
			words: []formatWord{
				{Str: "some string here"},
			},
		},
		{
			name:   "string",
			format: `%s`,
			words: []formatWord{
				{Str: "%s", Verb: true},
			},
		},
		{
			name:   "wstring",
			format: `%S`,
			words: []formatWord{
				{Str: "%S", Verb: true},
			},
		},
		{
			name:   "int",
			format: `%d`,
			words: []formatWord{
				{Str: "%d", Verb: true},
			},
		},
		{
			name:   "verb repeat",
			format: `%dd`,
			words: []formatWord{
				{Str: "%d", Verb: true},
				{Str: "d"},
			},
		},
		{
			name:   "mixed start",
			format: `%d num`,
			words: []formatWord{
				{Str: "%d", Verb: true},
				{Str: " num"},
			},
		},
		{
			name:   "mixed start 2",
			format: `%dnum`,
			words: []formatWord{
				{Str: "%d", Verb: true},
				{Str: "num"},
			},
		},
		{
			name:   "mixed middle",
			format: `v = %d num`,
			words: []formatWord{
				{Str: "v = "},
				{Str: "%d", Verb: true},
				{Str: " num"},
			},
		},
		{
			name:   "mixed middle 2",
			format: `v=%dnum`,
			words: []formatWord{
				{Str: "v="},
				{Str: "%d", Verb: true},
				{Str: "num"},
			},
		},
		{
			name:   "mixed end",
			format: `v = %d`,
			words: []formatWord{
				{Str: "v = "},
				{Str: "%d", Verb: true},
			},
		},
		{
			name:   "mixed end 2",
			format: `v=%d`,
			words: []formatWord{
				{Str: "v="},
				{Str: "%d", Verb: true},
			},
		},
		{
			name:   "percents",
			format: `%d%% = %%%d %%`,
			words: []formatWord{
				{Str: "%d", Verb: true},
				{Str: "%%", Verb: true},
				{Str: " = "},
				{Str: "%%", Verb: true},
				{Str: "%d", Verb: true},
				{Str: " "},
				{Str: "%%", Verb: true},
			},
		},
		{
			name:   "percents 2",
			format: `%%%d = %d%% `,
			words: []formatWord{
				{Str: "%%", Verb: true},
				{Str: "%d", Verb: true},
				{Str: " = "},
				{Str: "%d", Verb: true},
				{Str: "%%", Verb: true},
				{Str: " "},
			},
		},
		{
			name:   "all types",
			format: "%%%c%d%ld%10d%010d%x%o%#x%#o%4.2f%+.0e%E%*d%s%S%%",
			words: []formatWord{
				{Str: "%%", Verb: true},
				{Str: "%c", Verb: true},
				{Str: "%d", Verb: true},
				{Str: "%ld", Verb: true},
				{Str: "%10d", Verb: true},
				{Str: "%010d", Verb: true},
				{Str: "%x", Verb: true},
				{Str: "%o", Verb: true},
				{Str: "%#x", Verb: true},
				{Str: "%#o", Verb: true},
				{Str: "%4.2f", Verb: true},
				{Str: "%+.0e", Verb: true},
				{Str: "%E", Verb: true},
				{Str: "%*d", Verb: true},
				{Str: "%s", Verb: true},
				{Str: "%S", Verb: true},
				{Str: "%%", Verb: true},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			words := parseFormat(c.format)
			require.Equal(t, c.words, words)
		})
	}
}
