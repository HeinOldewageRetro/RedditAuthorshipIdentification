package main

import (
	"RedditScraper/util"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	AllAuthors := getAllAuthorsForByteWave()

	fmt.Println("Num Authors", len(AllAuthors))

	FilteredAuthors := make(map[string]bool)
	for _, author := range AllAuthors {
		FilteredAuthors[author] = true
	}

	file, err := os.Create("bytewave.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(file)

	for author := range FilteredAuthors {
		fmt.Println("Getting comments for", author)
		c := util.GetAllComments(author)
		for _, comment := range c {
			w.Write([]string{author, comment})
		}
	}

	fmt.Println("Done Downloading, writing")
	w.Flush()
	file.Close()
}

func getAllAuthorsForByteWave() (AllAuthors []string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	AllAuthors = make([]string, 0)
	after := ""
	more := true
	for more {
		var a []string
		a, after, more = util.GetBytewaveAuthorNames(after)
		AllAuthors = append(AllAuthors, a...)
	}
	return AllAuthors
}
