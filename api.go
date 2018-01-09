// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package opencorporates is an unofficial Golang API client for the OpenCorporates.
// http://api.opencorporates.com/documentation/API-Reference
package opencorporates

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// Version is the API default version
	Version = "0.4"
	// URL is the base URL of the OpenCorporates service.
	URL = "https://api.opencorporates.com/v%s/"
)

// Method calls
const (
	// ByNameURL is the URL to search a company by name.
	ByNameURL = "companies/search"
	// ByNumberURL is the URL to search a company by identifier.
	ByNumberURL = "companies/%s/%s"
)

// Error messages.
var (
	// ErrMethod is the error for unknown method call.
	ErrMethod = errors.New("unknown method call")
	// ErrMissingCountry is the error returned if the jurisdiction code is missing.
	ErrMissingCountry = errors.New("missing jurisdiction code")
)

// Date represents a date without time.
type Date struct {
	Time time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Date) UnmarshalJSON(b []byte) (err error) {
	data := strings.Trim(string(b), "\"")
	if data == "null" {
		return nil
	}
	m.Time, err = time.Parse("2006-01-02", data)
	return err
}

// Getter represents the mean to do a HTTP get.
type Getter interface {
	Get(url string) (*http.Response, error)
}

// API represents the API client.
type API struct {
	Version string
	http    Getter
}

// Address represents the company's address.
type Address struct {
	Street     string `json:"street_address"`
	City       string `json:"locality"`
	Region     string `json:"region,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// String implements the strings.Stringer interface.
func (addr Address) String() string {
	return fmt.Sprintf(
		"%s, %s, %s, %s %s",
		addr.Street, addr.City, addr.Region, addr.PostalCode, addr.Country,
	)
}

// Company represents a company.
type Company struct {
	Name            string  `json:"name"`
	Kind            string  `json:"company_type"`
	Number          string  `json:"company_number"`
	CountryCode     string  `json:"jurisdiction_code,omitempty"`
	CreationDate    Date    `json:"incorporation_date"`
	DissolutionDate Date    `json:"dissolution_date,omitempty"`
	Address         Address `json:"registered_address,omitempty"`
}

// ByCompanyNumber returns the company by its identifier and jurisdiction code.
// companies/fr/529591737
func (api *API) ByCompanyNumber(n, code string) (c Company, err error) {
	url, err := api.url(ByNumberURL, n, code)
	if err != nil {
		return
	}
	resp, err := api.call(url)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// Only accepts valid response.
	if resp.StatusCode > http.StatusNotFound {
		err = errors.New(resp.Status)
		return
	}

	// Response contains the search response.
	type Response struct {
		Results struct {
			Company Company `json:"company"`
		} `json:"results"`
	}
	var res Response
	if err = json.NewDecoder(resp.Body).Decode(&res); err == nil {
		c = res.Results.Company
	}
	return
}

// ByCompanyName calls the search method of the API to lookup companies by name.
// The second parameter, optional, allows to filter by jurisdiction code.
// companies/search?q=nautic+motors+evasion&jurisdiction_code=fr
func (api *API) ByCompanyName(q, code string) ([]Company, error) {
	url, err := api.url(ByNameURL, q, code)
	if err != nil {
		return nil, err
	}
	resp, err := api.call(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Response contains the search response.
	type Response struct {
		Results struct {
			Companies []struct {
				Company Company `json:"company,omitempty"`
			} `json:"companies,omitempty"`
		} `json:"results"`
	}
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	var list []Company
	for _, c := range res.Results.Companies {
		list = append(list, c.Company)
	}
	return list, nil
}

func (api *API) call(url string) (*http.Response, error) {
	if api.http == nil {
		// Default client.
		api.http = http.DefaultClient
	}
	return api.http.Get(url)
}

func (api *API) url(method string, param ...string) (string, error) {
	if api.Version == "" {
		// Default value.
		api.Version = Version
	}
	switch method {
	case ByNumberURL:
		var number, country string
		if len(param) == 2 {
			number = url.QueryEscape(param[0])
			country = url.QueryEscape(param[1])
		}
		if _, err := strconv.Atoi(number); err != nil {
			return "", err
		}
		if country == "" {
			return "", ErrMissingCountry
		}
		return fmt.Sprintf(URL+ByNumberURL, api.Version, country, number), nil
	case ByNameURL:
		var query, country string
		switch len(param) {
		case 2:
			country = url.QueryEscape(param[1])
			fallthrough
		case 1:
			query = url.QueryEscape(param[0])
		}
		if country != "" {
			return fmt.Sprintf(URL+ByNameURL+"?q=%s&jurisdiction_code=%s", api.Version, query, country), nil
		}
		return fmt.Sprintf(URL+ByNameURL+"?q=%s", api.Version, query), nil
	}
	return "", ErrMethod
}

// UseClient allows to use your own HTTP client to request the API.
func (api *API) UseClient(http Getter) *API {
	api.http = http
	return api
}
