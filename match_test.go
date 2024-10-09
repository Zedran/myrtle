package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

/* Tests the ParseQuery functions. The function should return the slice of non-nil Match struct
 * pointers equal in length to the number of sets in a sample string. The length of Lines 1 and 2
 * must be 69 chars long.
 */
func TestParseQuery(t *testing.T) {
	const LINE_LEN = 69

	// test sample must contain a newline count equal to lines count (hence nl at the end)
	const sample = `ISS (ZARYA)             
1 25544U 98067A   22014.20078024 -.00001581  00000+0 -20061-4 0  9991
2 25544  51.6452  19.1428 0006828  17.5887  10.3753 15.49476744321309
SWISSCUBE               
1 35932U 09051B   22013.55765441  .00000268  00000+0  71136-4 0  9999
2 35932  98.5837 225.0161 0007892 155.7494 204.4076 14.56655971653547
`

	// Mock response needs to have 0d0a line breaks
	mockResp := &http.Response{
		Body: io.NopCloser(bytes.NewBuffer([]byte(strings.Replace(sample, "\n", "\r\n", -1)))),
	}
	defer mockResp.Body.Close()

	output, err := ParseQuery(mockResp)
	if err != nil {
		t.Fatalf("Error when parsing response: %v", err)
	}

	// every set is a 3-liner, therefore
	if len(output) != strings.Count(sample, "\n")/3 {
		t.Fatalf("Improper output slice length: %d", len(output))
	}

	// testing the length of Lines 1 and 2
	for i := range output {
		if len(output[i].Line1) != LINE_LEN {
			t.Fatalf("Improper line 1 length for struct %d", i+1)
		} else if len(output[i].Line2) != LINE_LEN {
			t.Fatalf("Improper line 2 length for struct %d", i+1)
		}
	}
}
