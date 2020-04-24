package parsers

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/price_detector/models"
)

type Avtopro struct {
}

func (ap *Avtopro) Parse(r io.Reader) (*models.Records, error) {
	records := make(models.Records, 0)

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create html dom: %w", err)
	}
	//TODO: check nil pointers
	table := doc.Find("table").Find("tbody")
	rows := table.Children()
	for row := rows.First(); row.Nodes != nil; row = row.Next() {
		cols := row.Children()
		record := models.Record{}
		record.Seller = row.AttrOr("data-seller-name", "N/A")
		for col := cols.First(); col.Nodes != nil; col = col.Next() {
			t, ok := col.Attr("data-type")
			if ok {
				switch t {
				case "price":
					if v, ok := col.Attr("data-value"); ok {
						record.Price = v
					}
				case "delivery":
					if v, ok := col.Attr("data-city"); ok {
						record.City = v
					}
				case "code":
					record.ID = strings.TrimSpace(col.Text())
				}
			}
		}
		records = append(records, record)
	}

	return &records, nil
}
