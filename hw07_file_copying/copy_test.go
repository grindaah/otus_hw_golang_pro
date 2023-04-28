package main

import (
	"errors"
	"testing"
)

func TestCopy(t *testing.T) {
	testcases := []struct {
		name      string
		limit     int64
		offset    int64
		inFile    string
		outFile   string
		expectErr error
	}{
		{
			name:      "check offset error",
			limit:     1000,
			offset:    100000,
			inFile:    "./testdata/input.txt",
			outFile:   "./out.txt",
			expectErr: ErrOffsetExceedsFileSize,
		},
	}

	for _, tc := range testcases {
		err := Copy(tc.inFile, tc.outFile, tc.limit, tc.offset)
		if tc.expectErr != nil {
			if errors.Is(err, tc.expectErr) {
				t.Log("PASSED", tc.name)
			}
		} else {
			if err != nil {
				t.Error("err must be nil")
			}
		}
	}
}
