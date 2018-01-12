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
	// ByNameURL is the path to search a company by name or jurisdiction.
	ByNameURL = "companies/search"
	// ByNumberURL is the path to search a company by identifier.
	ByNumberURL = "companies/%s/%s?sparse=true"
)

// Error messages.
var (
	// ErrMethod is the error for unknown method call.
	ErrMethod = errors.New("unknown method call")
	// ErrJurisdiction is the error returned if the jurisdiction code is missing.
	ErrJurisdiction = errors.New("missing jurisdiction code")
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
	CountryCode     string  `json:"country_code,omitempty"`
	jurisdiction    string  `json:"jurisdiction_code,omitempty"`
	CreationDate    Date    `json:"incorporation_date"`
	DissolutionDate Date    `json:"dissolution_date,omitempty"`
	Address         Address `json:"registered_address,omitempty"`
}

// Companies returns an iterator of companies with its name or / and in this jurisdiction.
func (api *API) Companies(name, jurisdiction string) *CompanyIterator {
	return &CompanyIterator{
		api:          api,
		page:         NewPager(1),
		name:         name,
		jurisdiction: jurisdiction,
	}
}

// CompanyByID returns the company by its identifier and jurisdiction code.
// companies/fr/529591737
func (api *API) CompanyByID(id, jurisdiction string) (c Company, err error) {
	url, err := api.url(ByNumberURL, id, jurisdiction)
	if err != nil {
		return
	}
	resp, err := api.call(url)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

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

func (api *API) call(url string) (resp *http.Response, err error) {
	if api.http == nil {
		// Default client.
		api.http = http.DefaultClient
	}
	resp, err = api.http.Get(url)
	if err != nil {
		return
	}
	// Only accepts valid response.
	if resp.StatusCode > http.StatusBadRequest {
		if resp.StatusCode < http.StatusInternalServerError {
			// Response contains the search response.
			type Response struct {
				Err struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			var res Response
			if jErr := json.NewDecoder(resp.Body).Decode(&res); jErr == nil {
				err = errors.New(res.Err.Message)
			}
		}
		if err == nil {
			err = errors.New(resp.Status)
		}
		_ = resp.Body.Close()
	}
	return
}

func (api *API) url(method string, param ...interface{}) (string, error) {
	if api.Version == "" {
		// Default value.
		api.Version = Version
	}
	switch method {
	case ByNumberURL:
		var id, jurisdiction string
		if len(param) == 2 {
			id = url.QueryEscape(param[0].(string))
			jurisdiction = url.QueryEscape(param[1].(string))
		}
		if _, err := strconv.Atoi(id); err != nil {
			return "", err
		}
		if jurisdiction == "" {
			return "", ErrJurisdiction
		}
		return fmt.Sprintf(URL+ByNumberURL, api.Version, jurisdiction, id), nil
	case ByNameURL:
		var q, jurisdiction string
		var page int
		switch len(param) {
		case 3:
			page = param[2].(int)
			fallthrough
		case 2:
			jurisdiction = url.QueryEscape(param[1].(string))
			fallthrough
		case 1:
			q = url.QueryEscape(param[0].(string))
		}
		query := URL + ByNameURL + "?page=%d&order=score&q=%s"
		if jurisdiction != "" {
			return fmt.Sprintf(query+"&jurisdiction_code=%s", api.Version, page, q, jurisdiction), nil
		}
		return fmt.Sprintf(query, api.Version, page, q), nil
	}
	return "", ErrMethod
}

// UseClient allows to use your own HTTP client to request the API.
func (api *API) UseClient(http Getter) *API {
	api.http = http
	return api
}
