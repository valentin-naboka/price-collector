package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/price-collector/processor"
)

func main() {
	f, err := excelize.OpenFile("input.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	sheetName := f.GetSheetName(1)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	partIDs := make([]string, len(rows))
	for i, r := range rows {
		if i == 0 {
			continue
		}
		partIDs[i] = r[1]
	}

	proc := processor.NewAvtoproPool()
	prices, _ := proc.Do(partIDs)
	for _, price := range prices {
		if price.UsedMin == 0 && price.UsedMax == 0 && price.NewMin == 0 && price.NewMax == 0 {
			continue
		} else {
			if price.UsedMin != 0 {
				f.SetCellFloat(sheetName, "C"+strconv.Itoa(price.Row+1), price.UsedMin, 0, 64)
			}
			if price.UsedMax != 0 {
				f.SetCellFloat(sheetName, "D"+strconv.Itoa(price.Row+1), price.UsedMax, 0, 64)
			}
			if price.NewMin != 0 {
				f.SetCellFloat(sheetName, "E"+strconv.Itoa(price.Row+1), price.NewMin, 0, 64)
			}
			if price.NewMax != 0 {
				f.SetCellFloat(sheetName, "F"+strconv.Itoa(price.Row+1), price.NewMax, 0, 64)
			}
		}
	}
	f.SaveAs("output.xlsx")
}
