# OpenCorporates

[![GoDoc](https://godoc.org/github.com/rvflash/opencorporates?status.svg)](https://godoc.org/github.com/rvflash/opencorporates)
[![Build Status](https://img.shields.io/travis/rvflash/opencorporates.svg)](https://travis-ci.org/rvflash/opencorporates)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/opencorporates.svg)](http://codecov.io/github/rvflash/opencorporates?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/opencorporates)](https://goreportcard.com/report/github.com/rvflash/opencorporates)


Golang interface for the OpenCorporates service.


### Installation

```bash
$ go get -u github.com/rvflash/opencorporates
```

### Usage

Package still in development, all the method are not implemented yet.

The import of the package and check of errors are ignored for the demo.


#### Search a company by its name.

```go
api := &opencorporates.API{}
list, _ := api.ByCompanyName("nautic motors evasion", "fr")
for _, company := range list {
    println(company.Name+" #"+ company.Number)
}
// Output: NAUTIC MOTOR'S EVASION 35 (810622795)
// SARL NAUTIC MOTOR'S EVASION (529591737)
```

#### Search a company by its number (identifier).

```go
api := &opencorporates.API{}
company, _ := api.ByCompanyNumber("529591737", "fr")
println(company.Name+" #"+ company.Number)
// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
```