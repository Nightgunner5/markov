package markov

import (
	//"math/rand"
	"strings"
	"testing"
)

var (
	// From http://en.wikipedia.org/wiki/Potato
	Corpus = splitSentences(strings.Fields(`The potato is a starchy, tuberous crop from the perennial Solanum tuberosum of the Solanaceae family (also known as the nightshades). The word may refer to the plant itself as well as the edible tuber. In the region of the Andes, there are some other closely related cultivated potato species. Potatoes were introduced outside the Andes region four centuries ago, and have become an integral part of much of the world's cuisine. It is the world's fourth-largest food crop, following rice, wheat and maize. Long-term storage of potatoes requires specialised care in cold warehouses.
Wild potato species occur throughout the Americas, from the United States to southern Chile. The potato was originally believed to have been domesticated independently in multiple locations, but later genetic testing of the wide variety of cultivars and wild species proved a single origin for potatoes in the area of present-day southern Peru and extreme northwestern Bolivia (from a species in the Solanum brevicaule complex), where they were domesticated 7,000–10,000 years ago. Following centuries of selective breeding, there are now over a thousand different types of potatoes. Of these subspecies, a variety that at one point grew in the Chiloé Archipelago (the potato's south-central Chilean sub-center of origin) left its germplasm on over 99% of the cultivated potatoes worldwide.
The annual diet of an average global citizen in the first decade of the 21st century included about 33 kg (73 lb) of potato. However, the local importance of potato is extremely variable and rapidly changing. It remains an essential crop in Europe (especially eastern and central Europe), where per capita production is still the highest in the world, but the most rapid expansion over the past few decades has occurred in southern and eastern Asia. China is now the world's largest potato-producing country, and nearly a third of the world's potatoes are harvested in China and India.
The English word potato comes from Spanish patata (the name used in Spain). The Spanish Royal Academy says the Spanish word is a compound of the Taino batata (sweet potato) and the Quechua papa (potato). The name potato originally referred to a type of sweet potato rather than the other way around, although there is actually no close relationship between the two plants. The English confused the two plants one for the other. In many of the chronicles detailing agriculture and plants, no distinction is made between the two. The 16th-century English herbalist John Gerard used the terms "bastard potatoes" and "Virginia potatoes" for this species, and referred to sweet potatoes as "common potatoes". Potatoes are occasionally referred to as "Irish potatoes" or "white potatoes" in the United States, to distinguish them from sweet potatoes.
The name spud for a small potato comes from the digging of soil (or a hole) prior to the planting of potatoes. The word has an unknown origin and was originally (c. 1440) used as a term for a short knife or dagger, probably related to Dutch spyd and/or the Latin "spad-" root meaning "sword"; cf. Spanish "espada", English "spade" and "spadroon". The word spud traces back to the 16th century. It subsequently transferred over to a variety of digging tools. Around 1845 it transferred over to the tuber itself. The origin of "spud" has erroneously been attributed to a 19th century activist group dedicated to keeping the potato out of Britain, calling itself The Society for the Prevention of an Unwholesome Diet. It was Mario Pei's 1949 The Story of Language that can be blamed for the false origin. Pei writes, "the potato, for its part, was in disrepute some centuries ago. Some Englishmen who did not fancy potatoes formed a Society for the Prevention of Unwholesome Diet. The initials of the main words in this title gave rise to spud." Like most other pre-20th century acronymic origins, this one is false.`))
)

func splitSentences(words []string) [][]string {
	var sentences [][]string

outer:
	for len(words) > 0 {
		for i, word := range words {
			if word[len(word)-1] == '.' {
				sentences = append(sentences, words[:i+1])
				words = words[i+1:]
				continue outer
			}
		}

		sentences = append(sentences, words)
		words = nil
	}
}

func BenchmarkStringTableBuild1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := NewSTChain()

		c.Build(Corpus)
	}
}

func BenchmarkOriginalBuild1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := NewChain(PrefixLength)

		c.Build(Corpus)
	}
}
