// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package opencorporates

import (
	"encoding/json"
	"errors"
)

// EOF indicates if the iterator is done.
var EOF = errors.New("no more items in iterator")

// Pager contains information about an iterator's paging state.
type Pager struct {
	// Current page number
	curPage,
	// The maximum number of page
	maxPage,
	// Size of the full request
	TotalSize,
	// Number of item per page
	perPage,
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

// Remaining returns the number of items available.
func (i *Pager) Remaining() int {
	if i.curPage == 0 || i.TotalSize <= 0 {
		return 0
	}
	if i.perPage > i.TotalSize {
		return i.TotalSize - i.pos
	}
	return i.TotalSize - ((i.perPage * (i.curPage - 1)) + i.pos)
}

// Pageable is implemented by iterators that support paging.
type Pageable interface {
	Info() *Pager
}

// CompanyIterator iterates threw a list of companies.
type CompanyIterator struct {
	api  *API
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
		// Initializes the iterator.
		ci.resp, ci.page, ci.err = ci.api.companies(ci.name, ci.jurisdiction, ci.page.curPage)
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

func (api *API) companies(name, jurisdiction string, page int) (res []Company, info *Pager, err error) {
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
		curPage:   data.Results.Page,
		maxPage:   data.Results.NbPage,
		perPage:   data.Results.PerPage,
		TotalSize: data.Results.Nb,
	}
	res = make([]Company, min(info.perPage, info.TotalSize))
	for i, c := range data.Results.Companies {
		res[i] = c.Company
	}
	return
}

// Pager implements the Pageable interface.
// It returns information associated with the iterator.
func (ci *CompanyIterator) Info() *Pager {
	return ci.page
}
