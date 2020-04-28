package processor

import (
	"bytes"
	"log"
	"math"
	"sync"

	"github.com/price-collector/data"
	"github.com/price-collector/models"
	"github.com/price-collector/parsers"
)

type Price struct {
	Row     int
	UsedMin float64
	UsedMax float64
	NewMin  float64
	NewMax  float64
}

type part struct {
	Idx int
	ID  string
}

type WorkerPool interface {
	Do(partsIDs []string) ([]Price, error)
}

func NewAvtoproPool() WorkerPool {
	return &avtoproPool{}
}

type avtoproPool struct {
}

const GoRoutinesNum int = 30

//TODO: return error
func (ap *avtoproPool) Do(partsID []string) ([]Price, error) {
	input := make(chan part)
	output := make(chan Price)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(GoRoutinesNum)
		for i := 0; i < GoRoutinesNum; i++ {
			go func() {
				for part := range input {
					used, new := getPartPrices(part.ID)
					usedMin, usedMax := getMinMax(models.Records(used))
					newMin, newMax := getMinMax(models.Records(new))
					//TODO: remove
					//fmt.Printf("id %s min used %.0f, max used %.0f, min new %.0f, max new %.0f\n", part.ID, usedMin, usedMax, newMin, newMax)

					p := Price{}
					p.Row = part.Idx
					p.UsedMin = usedMin
					p.UsedMax = usedMax
					p.NewMin = newMin
					p.NewMax = newMax
					output <- p
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(output)
	}()

	prices := make([]Price, len(partsID))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		i := 0
		for o := range output {
			prices[i] = o
			i++
		}
		wg.Done()
	}()

	for i, partID := range partsID {
		input <- part{i, partID}
	}
	close(input)
	wg.Wait()
	return prices, nil
}

func getMinMax(r models.Records) (float64, float64) {
	min := math.MaxFloat64
	max := 0.0
	for _, v := range r {
		if v.Price < min {
			min = v.Price
		}

		if v.Price > max {
			max = v.Price
		}
	}
	if min == math.MaxFloat64 {
		min = 0
	}
	return float64(min), float64(max)
}

func getPartPrices(id string) (models.UsedRecords, models.NewRecords) {
	pages, err := data.FetchPages("https://avto.pro/api/v1/search/query", data.Request{id, 1})
	if err != nil {
		log.Fatal(err)
	}
	ap := parsers.Avtopro{}
	usedRecords := make(models.UsedRecords, 0)
	newRecords := make(models.NewRecords, 0)
	for _, p := range *pages {
		usedR, newR, err := ap.Parse(bytes.NewReader(p), id)
		if err != nil {
			log.Fatal(err)
		}

		usedRecords = append(usedRecords, *usedR...)
		newRecords = append(newRecords, *newR...)
	}
	return usedRecords, newRecords
}
