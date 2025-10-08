package scraper

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
	"github.com/PuerkitoBio/goquery"
)

type ViewStateData struct {
	ViewState          string
	ViewStateGenerator string
	EventValidation    string
}

func ExtractViewStateData(html []byte) (*ViewStateData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		return nil, fmt.Errorf("error when parsing html %s", err)
	}

	viewState, _ := doc.Find("input[name='__VIEWSTATE']").Attr("value")
	viewStateGen, _ := doc.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventValidation, _ := doc.Find("input[name='__EVENTVALIDATION']").Attr("value")

	return &ViewStateData{
		ViewState:          viewState,
		ViewStateGenerator: viewStateGen,
		EventValidation:    eventValidation,
	}, nil
}

func BuildFormData(viewStateDate *ViewStateData, mancoId int) url.Values {
	formData := url.Values{}
	formData.Set("__VIEWSTATE", viewStateDate.ViewState)
	formData.Set("__VIEWSTATEGENERATOR", viewStateDate.ViewStateGenerator)
	formData.Set("__EVENTVALIDATION", viewStateDate.EventValidation)
	formData.Set("MANCO_ID", fmt.Sprintf("%04d", mancoId))
	return formData
}

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
