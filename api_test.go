package main

import (
	"os"
	"testing"
)

const MAC_ADDRESS = "44:38:39:ff:ef:57"
const COMPANY_NAME = "Cumulus Networks, Inc"

func TestNew(t *testing.T) {
	cases := []struct {
		input string
		key   string
		error bool
	}{
		{"gibberish", "gibberish", false},
		{" gibberish with space ", "gibberish with space", false},
		{"  ", "", true},
		{"", "", true},
	}
	for _, c := range cases {
		m, err := New(c.input)
		if err != nil {
			if ! c.error {
				t.Errorf("returned an error %s when it was not expected for key \"%s\"", err.Error(), c.input)
			}
		} else {
			if c.error {
				t.Errorf("failed to return expected error for bad key \"%s\"", c.input)
			} else {
				if m.APIKey != c.key {
					t.Errorf("expected \"%s\" from key \"%s\", but got \"%s\"", c.key, c.input, m.APIKey)
				}
				if m.BaseURL != BaseURL {
					t.Errorf("wrong base-url: %s", m.BaseURL)
				}
			}
		}
	}
}

func TestParse(t *testing.T) {

	sample := `{
    "vendorDetails":{
        "oui":"443839",
        "isPrivate":false,
        "companyName":"Cumulus Networks, Inc",
        "companyAddress":"650 Castro Street, suite 120-245 Mountain View  CA  94041 US",
        "countryCode":"US"
    },
    "blockDetails":{
        "blockFound":true,
        "borderLeft":"443839000000",
        "borderRight":"443839FFFFFF",
        "blockSize":16777216,
        "assignmentBlockSize":"MA-L",
        "dateCreated":"2012-04-08",
        "dateUpdated":"2015-09-27"
    },
    "macAddressDetails":{
        "searchTerm":"44:38:39:ff:ef:57",
        "isValid":true,
        "virtualMachine":"Not detected",
        "applications":[
            "Multi-Chassis Link Aggregation (Cumulus Linux)"
        ],
        "transmissionType":"unicast",
        "administrationType":"UAA",
        "wiresharkNotes":"No details",
        "comment":""
    }
}` // Sample output from API documentation https://macaddress.io/api/documentation/output-format
	output, err := parse([]byte(sample))
	if err != nil {
		t.Errorf("error parsing sample output: %s", err.Error())
	} else {
		if output.VendorDetails.CompanyName != COMPANY_NAME {
			t.Errorf(
				"failed to extract correct company name; I got %s but was expecting %s",
				output.VendorDetails.CompanyName, COMPANY_NAME,
			)
		}
	}

}

func TestSource_Query(t *testing.T) {
	apiKey, envSet := os.LookupEnv("API_KEY")
	if envSet { // Only test if we have a key set.
		s, err := New(apiKey)
		if err == nil {
			name, err := s.Query(MAC_ADDRESS)
			if err != nil {
				t.Errorf("error querying the API: %s", err.Error())
			}
			if name != COMPANY_NAME {
				t.Errorf("Expected %s, but received %s", COMPANY_NAME, name)
			}
		} else {
			t.Errorf("error trying to initialise source: %s", err.Error())
		}

	}
}
