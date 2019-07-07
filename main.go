package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Default separator between MAC address and company name in the output.
const DefaultSeparator = "="

// The Environmental key to use for the API Key.
const APIKeyEnv = "API_KEY"

func main() {
	var apiKey string
	var separator string
	// API Key as argument.
	flag.StringVar(
		&apiKey,
		"key",
		"",
		fmt.Sprintf("API key to use with the source (has precedence over %s env)", APIKeyEnv),
	)
	// Custom separator between MAC address and company name at printout.
	flag.StringVar(
		&separator,
		"separator",
		DefaultSeparator,
		"Separator to use between MAC and company name in printout. Empty string to print only the company name",
	)
	// Do not print the MAC address; only the company name.
	var nameOnly = flag.Bool(
		"name-only",
		false,
		"Print the company name only in the output. Same as empty separator args",
	)
	flag.Parse()

	// An empty separator is taken as indication as not to print the MAC address.
	// It would make little sense to print them attached and
	// it saves to pass an extra parameter to the _query_ function.
	if *nameOnly {
		separator = ""
	}

	// The --key parameter takes precedence over the env variable as more explicit.
	if len(apiKey) == 0 {
		apiKey = strings.TrimSpace(os.Getenv(APIKeyEnv))
	}

	// The _New_ function already takes care of empty keys, but I want to print a more helpful error message.
	if len(apiKey) == 0 {
		printErrorMessage(
			fmt.Sprintf(
				"a string API Key from the --key argument or the %s environment variable is required",
				APIKeyEnv,
			),
			true,
		)
	}

	// Create a new source (could be easily converted to an interface, with multiple sources).
	source, err := New(apiKey)
	if err != nil {
		printErrorMessage(err.Error(), true)
	}
	// https://golang.org/pkg/flag/#NArg
	// If there are not args left we expect a pipe
	if flag.NArg() == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			query(source, strings.TrimSpace(line), separator)
		}
	} else { // We process all the mac addresses given (more than one allowed)
		for _, mac := range flag.Args() {
			query(source, strings.TrimSpace(mac), separator)
		}
	}
}

// Query the source and prints an error or the company name
func query(source Source, mac string, separator string) {
	if len(mac) > 0 {
		output, err := source.Query(mac)
		if err != nil {
			printErrorMessage(err.Error(), false)
		}
		if len(separator) > 0 {
			output = fmt.Sprintf("%s%s%s", mac, separator, output)
		}
		if len(output) > 0 { // Do not just print an empty line
			fmt.Fprintln(os.Stdout, output)
		}
	}
}

// printErrorMessage is a simple wrapper to print out error messages to stderr.
// An alternative to log.Fatal with log.SetFlags(0).
func printErrorMessage(text string, fatal bool) {
	fmt.Fprint(os.Stderr, text)
	if fatal {
		os.Exit(1)
	}
}
