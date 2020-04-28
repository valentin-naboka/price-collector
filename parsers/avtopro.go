package parsers

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/price-collector/models"
)

type Avtopro struct {
}

func (ap *Avtopro) Parse(r io.Reader, partID string) (*models.UsedRecords, *models.NewRecords, error) {
	usedRecords := make(models.UsedRecords, 0)
	newRecords := make(models.NewRecords, 0)

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create html dom: %w", err)
	}

	table := doc.Find("table").Find("tbody")
	rows := table.Children()
	for row := rows.First(); row.Nodes != nil; row = row.Next() {
		cols := row.Children()
		record := models.Record{}
		record.Seller = row.AttrOr("data-seller-name", "N/A")
		usedPart := false

		for col := cols.First(); col.Nodes != nil; col = col.Next() {
			if t, ok := col.Attr("data-type"); ok {
				switch t {
				case "price":
					if v, ok := col.Attr("data-value"); ok {
						f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
						if err != nil {
							log.Printf("failed to convert price %s to float: %s", v, err)
							break
						}
						record.Price = f
					}
				case "code":
					record.ID = strings.TrimSpace(col.Text())
				}
			}

			if v, ok := col.Attr("data-sub-title"); ok && v == "Б/У" {
				usedPart = true
			}
		}
		if record.ID == partID {
			if usedPart {
				usedRecords = append(usedRecords, record)
			} else {
				newRecords = append(newRecords, record)
			}
		}
	}

	return &usedRecords, &newRecords, nil
}
