package models

//Record ...
type Record struct {
	ID string `header:"ID"`
	//City   string `header:"City"`
	Price  float64 `header:"Price (USD)"`
	Seller string  `header:"Seller"`
}

type Records []Record

type UsedRecords Records
type NewRecords Records
