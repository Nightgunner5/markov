// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)
*/
package main

import (
	"math/rand"
	"strings"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of PrefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Chain     map[string]*Suffix
	PrefixLen int
}

type Suffix struct {
	T int
	W map[string]int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string]*Suffix), prefixLen}
}

// Build reads text from the provided []string and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(words []string) {
	p := make(Prefix, c.PrefixLen)
	for _, s := range words {
		c.increment(p.String(), s)
		p.Shift(s)
	}

	c.increment(p.String(), " ")
}

func (c *Chain) increment(prefix, suffix string) {
	s := c.Chain[prefix]
	if s == nil {
		s = &Suffix{0, make(map[string]int)}
		c.Chain[prefix] = s
	}

	s.T++
	s.W[suffix]++
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int, r *rand.Rand) string {
	p := make(Prefix, c.PrefixLen)

	var words []string
	for i := 0; i < n; i++ {
		suffix := c.Chain[p.String()]

		if suffix == nil || suffix.T == 0 {
			break
		}
		choice := r.Intn(suffix.T)
		for word, count := range suffix.W {
			choice -= count
			if choice < 0 {
				p.Shift(word)
				if word != " " {
					words = append(words, word)
				}
				break
			}
		}
	}
	return strings.Join(words, " ")
}
