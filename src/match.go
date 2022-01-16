package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/* Match struct contains a TLE set split into separate lines of text. */
type Match struct {
	Title, Line1, Line2 string
}

const (
	// API URL to be formatted with query value type and a value
	URL       string = "https://celestrak.com/NORAD/elements/gp.php?%s&FORMAT=TLE"

	// The length of the TLE Title Line
	TITLE_LEN int    = 24

	// Minimum query length - to avoid downloading of half the database
	MIN_QLEN  int    =  3
)

var (
	// The error returned from the Query function if both name and catnr are of zero length.
	errEmptyQuery = errors.New("both query values are empty")

	// The error returned from the Query function if neither of the query strings are longer than MIN_QLEN
	errShortQuery = errors.New("query value is too short")
)

/* Queries the API with object name or NORAD catalogue number. If both values are not of zero length,
 * the catalogue number is preferred. The result is a list of pointers to Match structs containing 
 * results.
 */
func Query(client *http.Client, name, catnr string) ([]*Match, error) {
	if len(name) < MIN_QLEN && len(catnr) < MIN_QLEN {
		return nil, errShortQuery
	}

	var queryValue string
	
	if len(catnr) != 0 {
		queryValue = "CATNR=" + catnr
	} else if len(name) != 0 {
		queryValue = "NAME=" + name
	} else {
		return nil, errEmptyQuery
	}

	resp, err := client.Get(fmt.Sprintf(URL, queryValue))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseQuery(resp)
}

/* Parses the response from the API and returns the results as a list of pointers to Match structs. */
func ParseQuery(resp *http.Response) ([]*Match, error) {
	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines   := strings.Split(string(stream), "\r\n")

	matches := make([]*Match, len(lines) / 3)

	iM := 0

	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == TITLE_LEN {
			matches[iM] = &Match{
				Title: lines[i    ],
				Line1: lines[i + 1],
				Line2: lines[i + 2],
			}
			
			iM++
			i += 2
		}
	}

	return matches, nil
}
