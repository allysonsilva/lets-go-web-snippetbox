package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Table-Driven Tests
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output (the want field).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// and anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}

	// // Initialize a new time.Time object and pass it to the humanDate function.
	// tm := time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// // Check that the output from the humanDate function is in the format we
	// // expect. If it isn't what we expect, use the t.Errorf() function to
	// // indicate that the test has failed and log the expected and actual
	// // values.
	// if hd != "17 Dec 2020 at 10:00" {
	// 	t.Errorf("want %q; got %q", "17 Dec 2020 at 10:00", hd)
	// }
}
