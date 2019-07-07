package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Entry point for macaddress.io API.
const BaseURL = `https://api.macaddress.io/v1`
// API query key for output format.
const OutputKey = "output"
// API query key for the MAC address.
const SearchKey = "search"
// Connection timeout.
const Timeout = time.Duration(time.Second*10) // Ten seconds timeout.
// Header for API Key as specified by the API.
const ApiKeyHeader = "X-Authentication-Token"

// Contains the basic elements for the query.
type Source struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// Output describes the expected JSON structure of the returning body.
type Output struct {
	VendorDetails Vendor `json: "vendorDetails"` // The company name is in the vendor details field
}

// Vendor describes the relevant JSON elements we are interested in.
type Vendor struct {
	CompanyName string `json: "companyName"` // We only care for company name as for now.
}

// New returns a struct to be used to query the MacAddress.io API
func New(apiKey string) (s Source, err error) {
	key := strings.TrimSpace(apiKey)
	if len(key) == 0 {
		err = errors.New("empty API key")
	} else {
		s.APIKey = key
		s.BaseURL = BaseURL
		s.Timeout = Timeout
	}
	return s, err
}

// Query connects to API entry point and tries to parse the response for the company name.
func (s *Source) Query(mac string) (string,error) {
	var companyName string
	var err error
	parameters := url.Values{}
	parameters.Set(OutputKey, "json")
	parameters.Set(SearchKey, mac)
	url, err := url.Parse(s.BaseURL)
	if err != nil {
		return companyName, err
	}
	client := &http.Client{
		Timeout: s.Timeout, // Timeout for the default transport is not very sophisticated, but will do in this case.
	}
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return companyName, err
	}
	request.URL.RawQuery = parameters.Encode()
	request.Header.Add(ApiKeyHeader, s.APIKey)
	response, err := client.Do(request)
	if err != nil {
		return companyName, newError("failed to retrieve API response", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return companyName, newError("failed to retrieve the response body", err)
	}
	if response.StatusCode == 200 {
		response.Body.Close()
		output, parseError := parse(body)
		if parseError == nil {
			companyName = output.VendorDetails.CompanyName // Only if everything went well we se the name
		} else {
			err = newError("failed to parse body for company name", parseError)
		}
	} else {
		switch response.StatusCode {
		// https://macaddress.io/api/documentation/error-codes
		// I would not use the "error" field in the JSON response body because:
		// 1. not documented
		// 2. not used consistently (422 doesn't return an JSON body, for example)
		case 400:
			err = newError("invalid parameters: %s", nil)
		case 401:
			err = newError("invalid API key", nil)
		case 402:
			err = newError("no available requests for the day left", nil)
		case 422:
			err = newError(fmt.Sprintf("invalid MAC or OUI address: %s", mac), nil)
		case 429:
			err = newError("too many requests; try again later", nil)
		case 500:
			err = newError("internal server error", nil)
		default:
			err = newError(fmt.Sprintf("unknown error code: %d", response.StatusCode),  nil)
		}
	}
	return companyName, err
}

// Parse will just return the result of JSON unmarshal, but this way it's
// easier to unit test.
func parse(body []byte) (Output, error) {
	var output Output
	err := json.Unmarshal(body, &output)
	return output, err
}

// newError is just a wrapper for assembling the error text
func newError(text string, err error) error {
	if err != nil {
		text = fmt.Sprintf("%s: %s", text, err.Error())
	}
	return errors.New(text)
}