package scraper

import (
	"bytes"
	"strconv"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
	"github.com/PuerkitoBio/goquery"
)

func ScrapeCISMangers(html []byte) ([]*models.CISManager, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		return nil, err
	}

	var managers []*models.CISManager
	doc.Find("select[name='MANCO_ID'] option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")

		integerVal, err := strconv.Atoi(value)

		if exists && err == nil && s.Text() != "" {
			manager := &models.CISManager{
				ID:   integerVal,
				Name: s.Text(),
			}
			managers = append(managers, manager)
		}
	})
	return managers, nil
}
