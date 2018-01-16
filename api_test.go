// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.
// Caution: we voluntary ignore errors for the demo.
package opencorporates_test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/rvflash/opencorporates"
)

func ExampleClient_CompanyByID() {
	api := opencorporates.API()
	company, _ := api.CompanyByID("529591737", "fr")
	fmt.Printf("%s (%s)\n", company.Name, company.Number)
	// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
}

func ExampleClient_Companies() {
	api := opencorporates.API()
	it := api.Companies("nautic motors evasion", "fr")
	for {
		company, err := it.Next()
		if err != nil {
			if err != opencorporates.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Printf("%s (%s)\n", company.Name, company.Number)
	}
	// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
	// NAUTIC MOTOR'S EVASION 35 (810622795)
}

func ExampleClient_UseClient() {
	api := opencorporates.API()
	proxyURL, _ := url.Parse("http://127.0.0.1:80")
	myClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	_, err := api.UseClient(myClient).CompanyByID("529591737", "fr")
	fmt.Println(err)
	// Output: Get https://api.opencorporates.com/v0.4/companies/fr/529591737?sparse=true: proxyconnect tcp: dial tcp 127.0.0.1:80: getsockopt: connection refused
}
