package scraper

import (
	"bytes"
	"fmt"
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

func ScrapeFunds(html []byte, managerId int) ([]*models.Fund, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		return nil, fmt.Errorf("failed to read document from html: %s", err)
	}

	var funds []*models.Fund
	doc.Find("select[name='TrustNo'] option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")

		intVal, err := strconv.Atoi(value)

		if exists && err == nil && s.Text() != "" {
			fund := &models.Fund{
				TrustNo:       intVal,
				Name:          s.Text(),
				SecondaryName: "",
				ManagerID:     managerId,
			}
			funds = append(funds, fund)
		}
	})

	return funds, nil
}
