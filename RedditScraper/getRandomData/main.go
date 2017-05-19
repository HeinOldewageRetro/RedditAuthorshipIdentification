package main

import (
	"RedditScraper/util"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

var numPages *int = flag.Int("numPages", 1, "The number of pages on the front page to get")

func main() {
	flag.Parse()
	AllAuthors := make([]string, 0)
	after := ""
	for k := 0; k < *numPages; k++ {
		var a []string
		a, after = util.GetAuthorNames(after)
		AllAuthors = append(AllAuthors, a...)
	}

	fmt.Println("Num Authors", len(AllAuthors))

	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(file)

	for _, author := range AllAuthors {
		c := util.GetAllComments(author)
		for _, comment := range c {
			w.Write([]string{author, comment})
		}
	}

	fmt.Println("Done Downloading, writing")
	w.Flush()
	file.Close()
}
