// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package opencorporates

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EOF indicates if the iterator is done.
var EOF = errors.New("no more items in iterator")

// Pager contains information about an iterator's paging state.
type Pager struct {
	// Single disables the automatic pagination.
	Single bool
	// Current page number
	curPage,
	// The maximum number of page for the request
	maxPage,
	// Number of item per page
	perPage,
	// Size of the full request
	total,
	// Current position in the buffer
	pos int
}

// NewPager exposes some statistics about the result.
func NewPager(start int) *Pager {
	if start <= 0 {
		start = 1
	}
	return &Pager{curPage: start}
}

// CurrentPage returns the current page number.
func (i *Pager) CurrentPage() int {
	return i.curPage
}

// Remaining returns the number of items available.
func (i *Pager) Remaining() int {
	if i.curPage == 0 || i.total <= 0 {
		return 0
	}
	if i.perPage > i.total {
		return i.total - i.pos
	}
	return i.total - ((i.perPage * (i.curPage - 1)) + i.pos)
}

// TotalCount returns the total number of items.
func (i *Pager) TotalCount() int {
	return i.total
}

// TotalPage returns the maximum number of page.
func (i *Pager) TotalPage() int {
	return i.maxPage
}

// Pageable is implemented by iterators that support paging.
type Pageable interface {
	Info() *Pager
}

// CompanyIterator iterates threw a list of companies.
type CompanyIterator struct {
	api  *client
	page *Pager
	name,
	jurisdiction string
	resp []Company
	err  error
}

// Next tries to return the next company in the iterator.
func (ci *CompanyIterator) Next() (Company, error) {
	var c Company
	if ci.page.pos == ci.page.perPage && ci.err == nil {
		if !ci.page.Single || ci.page.pos == 0 {
			if ci.page.pos > 0 {
				// Go to the next page
				ci.page.curPage++
			}
			// Retrieves the results
			ci.resp, ci.page, ci.err = ci.api.companies(ci.name, ci.jurisdiction, ci.page.curPage)
		} else {
			// No more results
			ci.page.curPage = 0
		}
	}
	if ci.err != nil {
		return c, ci.err
	}
	if ci.page.Remaining() == 0 {
		// No more company to iterate
		ci.err = EOF
	} else {
		c = ci.resp[ci.page.pos]
		ci.page.pos++
	}
	return c, ci.err
}

func (api *client) companies(name, jurisdiction string, page int) (res []Company, info *Pager, err error) {
	url, err := api.url(ByNameURL, name, jurisdiction, page)
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
			Companies []struct {
				Company Company `json:"company"`
			} `json:"companies"`
			Page    int `json:"page"`
			NbPage  int `json:"total_pages"`
			PerPage int `json:"per_page"`
			Nb      int `json:"total_count"`
		} `json:"results"`
	}
	var data Response
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	var min = func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	info = &Pager{
		curPage: data.Results.Page,
		maxPage: data.Results.NbPage,
		perPage: data.Results.PerPage,
		total:   data.Results.Nb,
	}
	res = make([]Company, min(info.perPage, info.total))
	for i, c := range data.Results.Companies {
		res[i] = c.Company
	}
	return
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
	Jurisdiction    string  `json:"jurisdiction_code,omitempty"`
	CreationDate    Date    `json:"incorporation_date"`
	DissolutionDate Date    `json:"dissolution_date,omitempty"`
	Address         Address `json:"registered_address,omitempty"`
}

// Pager implements the Pageable interface.
// It returns information associated with the iterator.
func (ci *CompanyIterator) Info() *Pager {
	return ci.page
}
