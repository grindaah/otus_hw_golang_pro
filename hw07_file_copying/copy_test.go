package main

import (
	"crypto/rand"
	"errors"
	"os"
	"testing"
)

type Result struct {
	sz   int64
	path string
}

func prepareTestData() {
	f, err := os.Create("./testdata/input_large.txt")
	if err != nil {
		return
	}
	buf := make([]byte, 1024*1024*20)
	// then we can call rand.Read.
	_, err = rand.Read(buf)
	if err != nil {
		return
	}
	f.Write(buf)
}

func deleteTestData() {
	_ = os.Remove("./testdata/input_large.txt")
}

func TestCopy(t *testing.T) {
	prepareTestData()
	defer deleteTestData()

	testcases := []struct {
		name      string
		limit     int64
		offset    int64
		inFile    string
		outFile   string
		expectErr error
		expectRes Result
	}{
		{
			name:      "check offset error",
			limit:     1000,
			offset:    100000,
			inFile:    "./testdata/input.txt",
			outFile:   "./out.txt",
			expectErr: ErrOffsetExceedsFileSize,
		},
		{
			name:      "check error on directory",
			limit:     1000,
			offset:    100000,
			inFile:    "./testdata",
			outFile:   "./out.txt",
			expectErr: ErrUnsupportedFile,
		},
		{
			name:    "check default offset is 0",
			inFile:  "./testdata/input.txt",
			outFile: "./out1.txt",
			expectRes: Result{
				sz:   6617,
				path: "./out.txt",
			},
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
			if tc.expectRes.path {
				f, errOpen := tc.expectRes.path
			}
		}
	}
}
