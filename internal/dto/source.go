package dto

import (
	"strconv"
	"strings"
	"unicode"
)

type ProductSource struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Brand    string `json:"brand"`
	Category string `json:"category"`
	Price    string `json:"price"`
	Stock    int    `json:"stock"`
}

type ClientSource struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Products  []uint `json:"products"`
}

func (p ProductSource) ParsePrice() float64 {
	var cleanPrice strings.Builder
	for _, ch := range p.Price {
		if unicode.IsDigit(ch) || ch == '.' {
			cleanPrice.WriteRune(ch)
		}
	}
	val, _ := strconv.ParseFloat(cleanPrice.String(), 64)
	return val
}
