package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/price_detector/data"
	"github.com/price_detector/models"
	"github.com/price_detector/parsers"
	"golang.org/x/net/html"
)

func main() {
	pages, err := data.FetchPages("https://avto.pro/api/v1/search/query", data.Request{"8200385222", 1})
	if err != nil {
		log.Fatal(err)
	}
	ap := parsers.Avtopro{}
	records := make(models.Records, 0)
	for _, p := range *pages {
		r, err := ap.Parse(bytes.NewReader(p))
		if err != nil {
			log.Fatal(err)
		}
		//TODO: remove duplicates
		//TODO: min price, max price, others in ascending order
		records = append(records, *r...)
	}

	for _, r := range records {
		fmt.Printf("%v\n", r)
	}
}

func traverse(n *html.Node) string {
	if n.Data == "table" {
		return n.Attr[0].Val
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := traverse(c); res != "" {
			return res
		}
	}
	return ""
}
