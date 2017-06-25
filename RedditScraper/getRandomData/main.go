package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/HeinOldewageRetro/RedditAuthorshipIdentification/RedditScraper/util"
)

var numPages *int = flag.Int("numPages", 1, "The number of pages on the front page to get")
var outPutFile *string = flag.String("o", "data.csv", "The file to sump the comments in. Won't contain duplicates")

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

	file, err := os.OpenFile(*outPutFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	names := make(map[string]struct{})
	r := csv.NewReader(file)
	for line, err := r.Read(); err == nil; line, err = r.Read() {
		names[line[0]] = struct{}{}
	}

	w := csv.NewWriter(file)

	for _, author := range AllAuthors {
		if _, ok := names[author]; ok {
			continue
		}
		names[author] = struct{}{}

		c := util.GetAllComments(author)
		for _, comment := range c {
			w.Write([]string{author, comment})
		}
	}

	fmt.Println("Done Downloading, writing")
	w.Flush()
	file.Close()
}
