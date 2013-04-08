package main

import (
	"compress/gzip"
	"encoding/gob"
	markov "github.com/Nightgunner5/markov/lib"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var Chain *markov.STChain

func load(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	Chain = markov.NewSTChain()
	err = gob.NewDecoder(g).Decode(&Chain)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.ParseInt(r.URL.Path[1:], 10, 64)
	if err != nil {
		http.Redirect(w, r, "/"+strconv.FormatInt(rand.Int63(), 10), 307)
		return
	}

	random := rand.New(rand.NewSource(i))

	w.Write([]byte(Chain.Generate(1000, random)))
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ", os.Args[0], " filename.gz")
	}

	load(os.Args[1])

	http.HandleFunc("/", handler)

	log.Print("Now listening on port 20098...")
	log.Fatal(http.ListenAndServe(":20098", nil))
}
