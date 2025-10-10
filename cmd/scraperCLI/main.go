package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Alexandervanderleek/FundFinderZA/internal/database"
	"github.com/Alexandervanderleek/FundFinderZA/internal/scraper"
	"github.com/joho/godotenv"
)

func main() {

	scrapeManco := flag.Bool("manco", false, "Scrape & Save CIS managers only")
	scrapeFunds := flag.Bool("funds", false, "Scrape & Save funds for saved managers")
	scrapePrices := flag.Bool("prices", false, "Scrape & Save all fund prices per class of fund")
	mancoIDs := flag.String("manco-ids", "", "Comman-seperated list of manco Ids to scrape and update.")

	flag.Parse()

	if !*scrapeFunds && !*scrapeManco && !*scrapePrices {
		log.Println("Usage: scraperCLI -manco | -funds | -prices [-manco-ids=0303,0037]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env file")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalln("Failed to parse database port")
	}

	dbConfig := &database.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMode"),
	}

	newDb, err := database.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}

	httpClient := scraper.NewClient(
		scraper.WithRetries(1),
		scraper.WithUserAgent("MyCustomUserAgent/1.0"),
		scraper.WithTimeout(30*time.Second))

	if *scrapeManco {
		if err := scrapeFundManagers(httpClient, newDb); err != nil {
			log.Fatalf("Failed to scrape fund managers: %s", err)
		}
	}

	if *scrapeFunds {
		if err := scrapeFundsForMangers(httpClient, newDb, mancoIDs); err != nil {
			log.Fatalf("Failed to scrape funds for managers: %s", err)
		}
	}

	if *scrapePrices {
		if err := ScrapeHistoricalPrices(httpClient, newDb); err != nil {
			log.Fatalf("Failed to scrape historical prices: %s", err)
		}
	}
}

func scrapeFundManagers(client *scraper.Client, db *database.DB) error {
	log.Println("Fetching CIS managers...")

	byteBody, err := client.Get("https://funds.profiledata.co.za/aci/ASISA/HistPriceLookUp.aspx")
	if err != nil {
		return fmt.Errorf("error fetching page: %s", err)
	}

	cisMangers, err := scraper.ScrapeCISMangers(byteBody)

	if err != nil {
		return fmt.Errorf("error scraping managers from page %s", err)
	}

	if err := db.SaveCISManagers(cisMangers); err != nil {
		return fmt.Errorf("error saving scraped cisManger: %s", err)
	}

	log.Printf("Succesfully saved %d fund managers\n", len(cisMangers))
	return nil
}

func scrapeFundsForMangers(client *scraper.Client, db *database.DB, mancoIds *string) error {
	log.Println("Fetching funds for managers...")

	var managerIdsToProcess []int
	if *mancoIds != "" {
		idStrs := strings.SplitSeq(*mancoIds, ",")
		for idStr := range idStrs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("error invalid manco id: %s : %w", idStr, err)
			}
			managerIdsToProcess = append(managerIdsToProcess, id)
		}
		log.Printf("Processing funds for %d specific managers\n", len(managerIdsToProcess))
	} else {
		mancoMangers, err := db.GetAllCISManagers()
		if err != nil {
			return fmt.Errorf("error getting cismanagers from db: %s", err)
		}

		for _, manager := range mancoMangers {
			managerIdsToProcess = append(managerIdsToProcess, manager.ID)
		}
		log.Printf("Processing funds for all %d managers\n", len(managerIdsToProcess))
	}

	for i, managerID := range managerIdsToProcess {
		log.Printf("[%d/%d] Processing manager ID: %d \n", i+1, len(managerIdsToProcess), managerID)

		initialHTML, err := client.Get("https://funds.profiledata.co.za/aci/ASISA/HistPriceLookUp.aspx")

		if err != nil {
			return fmt.Errorf("error fetching initial page: %s", err)
		}

		viewState, err := scraper.ExtractViewStateData(initialHTML)
		if err != nil {
			return fmt.Errorf("error extracting the view state: %s", err)
		}

		formData := scraper.BuildFormData(viewState, managerID)

		fundHtml, err := client.Post("https://funds.profiledata.co.za/aci/ASISA/HistPriceLookUp.aspx", formData)
		if err != nil {
			return fmt.Errorf("error posting form for manager: %s", err)
		}

		funds, err := scraper.ScrapeFunds(fundHtml, managerID)
		if err != nil {
			return fmt.Errorf("error scraping funds from html for ID - %d : %s", managerID, err)
		}

		if len(funds) > 0 {
			if err := db.SaveFunds(funds); err != nil {
				return fmt.Errorf("error saving funds : %s", err)
			}
		} else {
			log.Printf("No funds found")
		}

		time.Sleep(1 * time.Second)
	}
	log.Println("Completed the scraping of funds")
	return nil
}

func ScrapeHistoricalPrices(client *scraper.Client, db *database.DB) error {
	log.Println("Scraping Historical prices...")

	url := "https://funds.profiledata.co.za/aci/ASISA/LatestPrices.aspx"

	byteBody, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error getting latest price page: %s", err)
	}

	currentPriceDate, err := scraper.ScrapeCurrentPriceAndCostData(byteBody)
	if err != nil {
		return fmt.Errorf("error parsing scraped html: %s", err)
	}

	log.Printf("Scraped %d fund classes from prices page \n", len(currentPriceDate))

	matchedCount := 0
	unmatchedFunds := make([]string, 0)
	savedPrices := 0
	savedCosts := 0

	for _, data := range currentPriceDate {
		fundID, matchedName, err := db.FuzzyMatchFundName(data.FundClass.FundName)
		if err != nil {
			return fmt.Errorf("error fuzzy matching for fund: %s, %s", data.FundClass.FundName, err)
		}

		if fundID == 0 {
			unmatchedFunds = append(unmatchedFunds, fmt.Sprintf("%s %s", data.FundClass.FundName, data.FundClass.ClassName))
			continue
		}

		if matchedName != data.FundClass.FundName {
			log.Printf("Fuzzy Matched: %s -> %s \n", data.FundClass.FundName, matchedName)
		}

		matchedCount++
		data.FundClass.FundID = fundID

		if err := db.SaveFundClass(data.FundClass); err != nil {
			log.Printf("Error saving fund class for %s %s: %v\n",
				data.FundClass.FundName, data.FundClass.ClassName, err)
			continue
		}

		if data.Costs.TICDate != nil {
			data.Costs.FundClassID = data.FundClass.ID
			if err := db.SaveFundClassCosts(data.Costs); err != nil {
				log.Printf("Error saving costs for %s %s: %v\n",
					data.FundClass.FundName, data.FundClass.ClassName, err)
			} else {
				savedCosts++
			}
		}

		if data.Price.PriceDate != nil && data.Price.NAV != nil {
			data.Price.FundClassID = data.FundClass.ID
			if err := db.SaveFundClassPrice(data.Price); err != nil {
				log.Printf("Error saving price for %s %s: %v\n",
					data.FundClass.FundName, data.FundClass.ClassName, err)
			} else {
				savedPrices++
			}
		}
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("The following are considered unmatched...")

	for i, fund := range unmatchedFunds {
		fmt.Printf("Fund %d: %s\n", i, fund)
	}

	return nil

}
