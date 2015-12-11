package main

import (
	"os"
	"bufio"
	"log"
	"bytes"
	"flag"
	"github.com/cloudflare/ahocorasick"
	"fmt"
)

func normalizeBytes(s []byte) []byte {
	return bytes.ToLower(bytes.TrimSpace(s))
}

func main() {
	var dictFilename string
	flag.StringVar(&dictFilename, "dict", "dict.txt", "dictionary of most frequent words")
	var bigFilename string
	flag.StringVar(&bigFilename, "big", "big.txt", "file to analyze frequency of words")
	flag.Parse()

	dict := make([][]byte, 0)
	{
		log.Println("Starting reading dictionary")
		dictFile, err := os.Open(dictFilename)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(dictFile)
		for scanner.Scan() {
			line := normalizeBytes(scanner.Bytes())
			if len(line) > 2 {
				dict = append(dict, line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		log.Println("Finished parsing dictionary")
	}

	dictFreq := make([]int64, len(dict))
	aho := ahocorasick.NewMatcher(dict)
	{
		log.Println("Starting the big file")
		big, err := os.Open(bigFilename)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(big)
		counter := 0
		for scanner.Scan() {
			line := normalizeBytes(scanner.Bytes())
			hits := aho.Match(line)

			for _, v := range hits {
				dictFreq[v]++
			}

			if counter & ((1<<17)-1) == 0 {
				log.Printf("%d\n", counter)
			}

			counter++
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		log.Println("Finished reading the big file")
	}

	for i, word := range dict {
		if dictFreq[i] > 0 {
			fmt.Printf("%d,%s\n", dictFreq[i], word)
		}
	}
}