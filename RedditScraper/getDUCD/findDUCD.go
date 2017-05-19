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
	file, err := os.Create("36055512.csv")
	if err != nil {
		panic(err)
	}
	author := "36055512"
	w := csv.NewWriter(file)
	c := util.GetAllComments(author)
	for _, comment := range c {
		w.Write([]string{author, comment})
	}

	fmt.Println("Done Downloading, writing")
	w.Flush()
	file.Close()
}
