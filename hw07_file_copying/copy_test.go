package main

import (
	"crypto/rand"
	"errors"
	rand2 "math/rand"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	fileSz              = 1024 * 1024 * 20
	defaultChunkCompare = 1024
)

type Result struct {
	sz     int64
	path   string
	chunks map[int][]byte
}

type TempFile struct {
	*os.File
	// Contents is map of few random chunks
	Contents map[int][]byte
}

func prepareTestData() *TempFile {
	f, err := os.CreateTemp("", "inputlarge.*.txt")
	if err != nil {
		return nil
	}
	buf := make([]byte, fileSz)
	// then we can call rand.Read.
	_, err = rand.Read(buf)
	if err != nil {
		return nil
	}
	f.Write(buf)
	m := getChunks(buf, fileSz, defaultChunkCompare)
	if m == nil {
		return nil
	}

	return &TempFile{
		File:     f,
		Contents: m,
	}
}

func getChunks(input []byte, totalSize, chunkSize int) map[int][]byte {
	if chunkSize == 0 || input == nil || totalSize == 0 {
		return nil
	}
	if chunkSize > totalSize {
		chunkSize = totalSize
	}
	if input != nil {
		// get three random chunks + 1 chunk from the start and 1 from the end
		one := rand2.Intn(totalSize - chunkSize)
		two := rand2.Intn(totalSize - chunkSize)
		three := rand2.Intn(totalSize - chunkSize)
		result := make(map[int][]byte, 5)
		result[0] = input[:chunkSize]
		result[one] = input[one : one+chunkSize]
		result[two] = input[two : two+chunkSize]
		result[three] = input[three : three+chunkSize]
		result[totalSize-chunkSize] = input[totalSize-chunkSize:]
		return result
	}
	return nil
}

func compareChunks(f *os.File, chunks map[int][]byte, chunkSize int) bool {
	for pos, chunk := range chunks {
		f.Seek(int64(pos), 0)
		readChunk := make([]byte, chunkSize)
		f.Read(readChunk)
		if !cmp.Equal(chunk, readChunk) {
			return false
		}
	}
	return true
}

func TestCopy(t *testing.T) {
	tmpFile := prepareTestData()
	defer os.Remove(tmpFile.Name())

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
			name:    "large file 1",
			inFile:  tmpFile.Name(),
			outFile: "./out1.txt",
			expectRes: Result{
				sz:     fileSz,
				path:   tmpFile.Name(),
				chunks: tmpFile.Contents,
			},
		},
	}

	for _, tc := range testcases {
		err := Copy(tc.inFile, tc.outFile, tc.offset, tc.limit)
		if tc.expectErr != nil {
			if !errors.Is(err, tc.expectErr) {
				t.Fail()
			}
		}
		if err != nil {
			t.Error("err must be nil\n")
		}
		if len(tc.expectRes.path) > 0 {
			f, errOpen := os.Open(tc.expectRes.path)
			if errOpen != nil {
				t.Fail()
			}
			defer f.Close()
			fStat, _ := f.Stat()
			if fStat.Size() != tc.expectRes.sz {
				t.Error("size not matching\n")
				t.Fail()
			}
			if !compareChunks(f, tc.expectRes.chunks, defaultChunkCompare) {
				t.Error("chunk not matched\n")
				t.Fail()
			}
		}
	}
}
