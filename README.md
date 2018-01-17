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

Package still in development, all the methods are not implemented yet.

The import of the package and errors check are voluntary ignored for the demo.

See the test files for more example.

#### Search a company by its name or jurisdication .

```go
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
```

#### Search a company by its number (identifier).

```go
api := opencorporates.API()
company, _ := api.CompanyByID("529591737", "fr")
println(company.Name+" #"+ company.Number)
// Output: SARL NAUTIC MOTOR'S EVASION (529591737)
```
