// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package opencorporates_test

import (
	"fmt"

	"github.com/rvflash/opencorporates"
)

// Caution: we voluntary ignore errors for the demo.
var api = &opencorporates.API{}

func ExampleAPI_CompanyByID() {
	company, _ := api.CompanyByID("529591737", "fr")
	fmt.Printf("%s (%s)\n", company.Name, company.Number)
	// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
}

func ExampleAPI_Companies() {
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
