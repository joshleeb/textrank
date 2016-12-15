package textrank

import "testing"

func TestTokenizeSentences(t *testing.T) {
	cases := map[string]struct {
		Text      string
		Sentences []string
	}{
		"comma":     {"one, two", []string{"one, two"}},
		"empty":     {"", []string{}},
		"single":    {"word", []string{"word"}},
		"spaces":    {" ", []string{}},
		"untrimmed": {"one.    ", []string{"one"}},
		"double": {
			"a sentence. Now a second",
			[]string{"a sentence", "Now a second"},
		},
		"multiple": {
			"one sentence. Two sentence. more sentences",
			[]string{"one sentence", "Two sentence", "more sentences"},
		},
	}

	for k, tc := range cases {
		tokens := TokenizeSentences(tc.Text)
		if !eqStringSlices(tokens, tc.Sentences) {
			t.Errorf("%s: sentences = %#v, expected %#v", k, tokens, tc.Sentences)
		}
	}
}

func TestTokenizeWords(t *testing.T) {
	cases := map[string]struct {
		Text  string
		Words []string
	}{
		"apostrophes":     {"we've", []string{"weve"}},
		"comma":           {"some, word", []string{"some", "word"}},
		"double":          {"some word", []string{"some", "word"}},
		"empty":           {"", []string{}},
		"hyphen":          {"some-word", []string{"some", "word"}},
		"multiple":        {"a some word", []string{"a", "some", "word"}},
		"period":          {"some. word", []string{"some", "word"}},
		"preiod no space": {"some.word", []string{"some", "word"}},
		"single":          {"word", []string{"word"}},
		"spaces":          {" ", []string{}},
		"untrimmed":       {"  spaces  ", []string{"spaces"}},
	}

	for k, tc := range cases {
		tokens := TokenizeWords(tc.Text)
		if !eqStringSlices(tokens, tc.Words) {
			t.Errorf("%s: words = %#v, expected %#v", k, tokens, tc.Words)
		}
	}
}

func eqStringSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil || len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
