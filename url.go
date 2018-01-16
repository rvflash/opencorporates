// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package opencorporates

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

func companyNamedURL(version string, param ...interface{}) (string, error) {
	var (
		q, jurisdiction string
		page            int
	)
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
		return fmt.Sprintf(query+"&jurisdiction_code=%s", version, page, q, jurisdiction), nil
	}
	return fmt.Sprintf(query, version, page, q), nil
}

func companyNumberURL(version string, param ...interface{}) (string, error) {
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
	return fmt.Sprintf(URL+ByNumberURL, version, jurisdiction, id), nil
}

func (api *client) url(method string, param ...interface{}) (s string, err error) {
	if api.Version == "" {
		// Default value.
		api.Version = Version
	}
	// By method calls
	switch method {
	case ByNumberURL:
		s, err = companyNumberURL(api.Version, param...)
	case ByNameURL:
		s, err = companyNamedURL(api.Version, param...)
	default:
		err = ErrMethod
	}
	if err != nil {
		return
	}
	if api.Token != "" {
		// Adds the API key to increase the usage limits
		if strings.Contains(s, "?") {
			s += "&"
		} else {
			s += "?"
		}
		s += "&api_token=" + url.QueryEscape(api.Token)
	}
	return
}
