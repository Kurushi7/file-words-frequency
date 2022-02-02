package main

type Pair struct {
	Text      string `json:"text"`
	Occurence int    `json:"occurence"`
}

type PairList []Pair
