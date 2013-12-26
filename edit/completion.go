package edit

import (
	"io/ioutil"
	"../parse"
)

type tokenPart struct {
	text string
	completed bool
}

type candidate struct {
	text string
	parts []tokenPart
}

func newCandidate() *candidate {
	return &candidate{}
}

func (c *candidate) push(tp tokenPart) {
	c.text += tp.text
	c.parts = append(c.parts, tp)
}

type completion struct {
	candidates []*candidate
	current int
}

func (c *completion) prev() {
	if c.current > 0 {
		c.current--
	}
}

func (c *completion) next() {
	if c.current < len(c.candidates) - 1 {
		c.current++
	}
}

func findCandidates(text string) (candidates []*candidate) {
	// Find last token
	l := parse.Lex("<completion>", text)
	var lastToken parse.Item
	for token := range l.Chan() {
		if token.Typ != parse.ItemEOF {
			lastToken = token
		}
	}
	prefix := lastToken.Val

	infos, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}
	for _, info := range infos {
		name := info.Name()
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			cand := newCandidate()
			cand.push(tokenPart{prefix, false})
			cand.push(tokenPart{name[len(prefix):], true})
			candidates = append(candidates, cand)
		}
	}
	return
}