package main

import (
	"fmt"
	"time"

	"github.com/Alexandervanderleek/FundFinderZA/internal/scraper"
)

func main() {

	newClient := scraper.NewClient(
		scraper.WithRetries(1),
		scraper.WithUserAgent("MyCustomUserAgent/1.0"),
		scraper.WithTimeout(30*time.Second))

	byteBody, err := newClient.Get("https://funds.profiledata.co.za/aci/ASISA/HistPriceLookUp.aspx")

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	managers, err := scraper.ScrapeCISMangers(byteBody)

	if err != nil {
		fmt.Println("Error scraping managers:", err)
		return
	}

	for _, manager := range managers {
		fmt.Println("Manager:", manager.Name, "ID:", manager.ID)
	}
}
