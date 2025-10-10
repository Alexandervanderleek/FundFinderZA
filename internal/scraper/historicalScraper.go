package scraper

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
	"github.com/PuerkitoBio/goquery"
)

type FundPricingData struct {
	FundClass *models.FundClass
	Costs     *models.FundClassCost
	Price     *models.FundClassPrice
}

func ScrapeCurrentPriceAndCostData(html []byte) ([]*FundPricingData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("error reading in html document: %s", err)
	}

	var results []*FundPricingData
	currentCategory := ""

	doc.Find("#dataTable tr").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("sectorrow") {
			categoryText := s.Find("td").First().Text()
			currentCategory = strings.TrimSpace(categoryText)
			return
		}

		if !s.HasClass("fundrow") {
			return
		}

		tds := s.Find("td")
		if tds.Length() < 11 {
			return
		}

		fundNameFull := strings.TrimSpace(tds.Eq(0).Find("div.fundname").Text())
		if fundNameFull == "" {
			return
		}

		fundName, className := parseFundNameAndClass(fundNameFull)

		targetMarket := strings.TrimSpace(tds.Eq(2).Text())

		addFee := strings.TrimSpace(tds.Eq(1).Text()) == "yes"

		maxInitFeeStr := strings.TrimSpace(tds.Eq(3).Text())
		maxInitFee := parsePercentage(maxInitFeeStr)

		ticDate := parseDate(strings.TrimSpace(tds.Eq(4).Text()))

		terPerfCompStr := strings.TrimSpace(tds.Eq(5).Text())
		terPerfComp := parsePercentage(terPerfCompStr)

		terStr := strings.TrimSpace(tds.Eq(6).Text())
		ter := parsePercentage(terStr)

		tcStr := strings.TrimSpace(tds.Eq(7).Text())
		tc := parsePercentage(tcStr)

		ticStr := strings.TrimSpace(tds.Eq(8).Text())
		tic := parsePercentage(ticStr)

		priceDate := parseDate(strings.TrimSpace(tds.Eq(9).Text()))

		navStr := strings.TrimSpace(tds.Eq(10).Text())
		nav := parseDecimal(navStr)

		data := &FundPricingData{
			FundClass: &models.FundClass{
				FundName:     fundName,
				ClassName:    className,
				TargetMarket: targetMarket,
				AddFee:       addFee,
				MaxInitFee:   maxInitFee,
				Category:     currentCategory,
			},
			Costs: &models.FundClassCost{
				TICDate:     ticDate,
				TERPerfComp: terPerfComp,
				TER:         ter,
				TC:          tc,
				TIC:         tic,
			},
			Price: &models.FundClassPrice{
				PriceDate: priceDate,
				NAV:       nav,
			},
		}

		results = append(results, data)
	})
	return results, nil
}

func parseFundNameAndClass(fundNameFull string) (string, string) {
	classPatterns := []string{
		" Class ",
		" class "}

	for _, pattern := range classPatterns {
		if idx := strings.Index(fundNameFull, pattern); idx != -1 {
			fundName := strings.TrimSpace(fundNameFull[:idx])
			className := strings.TrimSpace(fundNameFull[idx+len(pattern):])
			return fundName, "Class " + className
		}
	}

	return fundNameFull, ""
}

func parsePercentage(percentage string) *float64 {
	percentage = strings.TrimSpace(percentage)
	if percentage == "n/a" || percentage == "" {
		return nil
	}

	percentage = strings.TrimSuffix(percentage, "%")
	val, err := strconv.ParseFloat(percentage, 64)

	if err != nil {
		return nil
	}

	return &val
}

func parseDecimal(decimal string) *float64 {
	decimal = strings.TrimSpace(decimal)

	if decimal == "n/a" || decimal == "" {
		return nil
	}

	val, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return nil
	}

	return &val
}

func parseDate(date string) *string {
	date = strings.TrimSpace(date)
	if date == "n/a" || date == "" {
		return nil
	}

	if strings.Contains(date, "/") {
		parts := strings.Split(date, "/")
		if len(parts) == 3 {
			day := parts[0]
			month := parts[1]
			year := parts[2]

			if len(year) == 2 {
				yearInt, _ := strconv.Atoi(year)

				if yearInt >= 0 && yearInt <= 50 {
					year = "20" + year
				} else {
					year = "19" + year
				}
			}

			result := fmt.Sprintf("%s-%s-%s", year, month, day)
			return &result
		}
	}

	if len(date) == 5 {
		monthMap := map[string]string{
			"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
			"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
			"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
		}

		monthStr := date[:3]
		yearStr := date[3:]

		if month, ok := monthMap[monthStr]; ok {
			yearInt, _ := strconv.Atoi(yearStr)
			if yearInt >= 0 && yearInt <= 50 {
				yearStr = "20" + yearStr
			} else {
				yearStr = "19" + yearStr
			}
			result := fmt.Sprintf("%s-%s-01", yearStr, month)
			return &result
		}
	}
	return nil
}
