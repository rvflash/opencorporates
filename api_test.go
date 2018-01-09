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

func ExampleByCompanyName() {
	list, _ := api.ByCompanyName("nautic motors evasion", "fr")
	for _, company := range list {
		fmt.Printf("%s (%s)\n", company.Name, company.Number)
	}
	// Output: NAUTIC MOTOR'S EVASION 35 (810622795)
	// SARL NAUTIC MOTOR'S EVASION (529591737)
}

func ExampleByCompanyNumber() {
	company, _ := api.ByCompanyNumber("529591737", "fr")
	fmt.Printf("%s (%s)\n", company.Name, company.Number)
	// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
}
