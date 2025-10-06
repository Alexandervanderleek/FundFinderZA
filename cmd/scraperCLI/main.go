package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Alexandervanderleek/FundFinderZA/internal/database"
	"github.com/Alexandervanderleek/FundFinderZA/internal/scraper"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load env file.")
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	dbConfig := &database.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMode"),
	}

	newDb, dbErr := database.NewDB(dbConfig)

	if dbErr != nil {
		log.Println("Failed to connect to the database!")
	}

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

	if err := newDb.SaveCisManagers(managers); err != nil {
		log.Println("Failed to store the managers %w", err)
	}

	if err != nil {
		fmt.Println("Error scraping managers:", err)
		return
	}

	for _, manager := range managers {
		fmt.Println("Manager:", manager.Name, "ID:", manager.ID)
	}
}
