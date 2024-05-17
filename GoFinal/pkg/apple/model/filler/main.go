package filler

import (
	model "github.com/kim0111/GoMidterm/pkg/apple/model"
)

func PopulateDatabase(models model.Models) error {
	for _, product := range products {
		models.Products.Insert(&product)
	}

	return nil
}

var products = []model.Products{
	{Title: "iPhone 17 pro MAX", Description: "A new iPhone", ForWhatCountry: "JP", Price: 1499},
	{Title: "MacBook Pro M4", Description: "Power of m4 CPU", ForWhatCountry: "USA", Price: 3500},
	{Title: "airPods Max v2", Description: "Amazing sound", ForWhatCountry: "EU", Price: 700},
	{Title: "iMac", Description: "Nice PC", ForWhatCountry: "USA", Price: 899},
	{Title: "apple TV", Description: "Only sub use", ForWhatCountry: "CH", Price: 399},
}
