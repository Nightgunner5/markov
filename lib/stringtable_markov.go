package markov

import (
	"math/rand"
	"strings"
)

const PrefixLength = 2

type STPrefix [PrefixLength]int

func (p *STPrefix) shift(i int) {
	copy((*p)[:], (*p)[1:])
	(*p)[PrefixLength-1] = i
}

type STSuffix struct {
	Count int
	Words map[int]int
}

type STChain struct {
	Strings []string
	Chain   map[STPrefix]STSuffix
	cache   map[string]int
}

func NewSTChain() *STChain {
	return &STChain{nil, make(map[STPrefix]STSuffix), nil}
}

func (c *STChain) Build(words []string) {
	if len(c.Strings) == 0 {
		c.Strings = []string{""}
	}

	if c.cache == nil {
		c.cache = make(map[string]int, len(c.Strings))
		for i, s := range c.Strings {
			c.cache[s] = i
		}
	}

	var p STPrefix
	for _, s := range words {
		i, ok := c.cache[s]
		if !ok {
			i = len(c.Strings)
			c.Strings = append(c.Strings, string([]byte(s)))
			c.cache[s] = i
		}

		c.increment(p, i)
		p.shift(i)
	}

	c.increment(p, 0)
}

func (c *STChain) increment(p STPrefix, s int) {
	suffix := c.Chain[p]
	if suffix.Words == nil {
		suffix.Words = make(map[int]int)
	}
	suffix.Words[s]++
	suffix.Count++
	c.Chain[p] = suffix
}

func (c *STChain) Generate(max int, r *rand.Rand) string {
	return strings.Join(c.GenerateSlice(max, r), " ")
}

func (c *STChain) GenerateSlice(max int, r *rand.Rand) []string {
	var s []string
	var p STPrefix

	for i := 0; i < max; i++ {
		suffix := c.Chain[p]
		if suffix.Count == 0 {
			return s
		}

		j := r.Intn(suffix.Count)

		for word, freq := range suffix.Words {
			j -= freq
			if j < 0 {
				if word == 0 {
					return s
				}

				p.shift(word)
				s = append(s, c.Strings[word])
			}
		}
	}

	return s
}
