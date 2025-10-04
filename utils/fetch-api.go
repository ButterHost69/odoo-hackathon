package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ButterHost69/odoo-hackathon/errs"
)

type Country struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

func GetCurrencyUsingCountryName(countryName string) (string, error) {
	apiURL := "https://restcountries.com/v3.1/all?fields=name,currencies"

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("[utils.GetCurrencyUsingCountryName] failed to make API request::  ", err)
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[utils.GetCurrencyUsingCountryName] API returned a non-200 status code:  ", resp.StatusCode)
		return "", fmt.Errorf("API returned a non-200 status code: %d", resp.StatusCode)
	}
	
	fmt.Println("[utils.GetCurrencyUsingCountryName] Fetched Data for Country", countryName)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[utils.GetCurrencyUsingCountryName] failed to read response body: ", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("[utils.GetCurrencyUsingCountryName] Read Data for Country", countryName)

	var countries []Country
	if err := json.Unmarshal(body, &countries); err != nil {
		fmt.Println("[utils.GetCurrencyUsingCountryName] failed to parse JSON response", err)
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	fmt.Println("[utils.GetCurrencyUsingCountryName] Parsed Data for Country", countryName)

	// Iterate over the slice of countries to find a match.
	for _, country := range countries {
		// Use a case-insensitive comparison for better matching.
		if strings.EqualFold(country.Name.Common, countryName) {
			for _, currencyInfo := range country.Currencies {
				fmt.Printf("[LOG] [utils.GetCurrencyUsingCountryName] Currency For Country %s -: %s\n", countryName, currencyInfo.Symbol)
				
					return currencyInfo.Symbol, nil // Return the symbol.
			}
			fmt.Println("[utils.GetCurrencyUsingCountryName] no currency information found for", countryName)
			return "", fmt.Errorf("no currency information found for %s", countryName)
		}
	}

	fmt.Println("[utils.GetCurrencyUsingCountryName] [LOG] Country Not Found", countryName)
	return "", errs.ErrCountryNotFound
}