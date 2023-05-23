package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"price-collector/processor"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	println("Starting...")
	list, err := filepath.Glob("./input/*.xlsx")
	if err != nil {
		log.Fatalf("error walking input dir: %s", err)
	}

	if err = os.RemoveAll("output"); err != nil {
		log.Fatalf("failed to remove dir: %s", err)
	}

	if err = os.Mkdir("output", os.ModePerm); err != nil {
		log.Fatalf("failed to create dir: %s", err)
	}

	for _, v := range list {
		_, filename := filepath.Split(v)
		f, err := excelize.OpenFile(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		sheetName := f.GetSheetName(1)
		rows := f.GetRows(sheetName)

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
					f.SetCellValue(sheetName, "C"+strconv.Itoa(price.Row+1), price.UsedMin)
				}
				if price.UsedMax != 0 {
					f.SetCellValue(sheetName, "D"+strconv.Itoa(price.Row+1), price.UsedMax)
				}
				if price.NewMin != 0 {
					f.SetCellValue(sheetName, "E"+strconv.Itoa(price.Row+1), price.NewMin)
				}
				if price.NewMax != 0 {
					f.SetCellValue(sheetName, "F"+strconv.Itoa(price.Row+1), price.NewMax)
				}
			}
		}
		f.SaveAs("./output/" + filename)
	}
	println("Done!")
	fmt.Scanln()
}
