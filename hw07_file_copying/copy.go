package main

import (
	"errors"
	"github.com/cheggaaa/pb"
	"io"
	"math"
	"os"
)

const chunkSizeDefault = 4096

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var sz int64
	sz, err = validateFile(f)
	if err != nil {
		return err
	}

	toFile, errCreate := os.Create(toPath)
	if errCreate != nil {
		return errCreate
	}
	toFile.Chmod(0644 | os.ModeAppend)
	defer toFile.Close()

	currentOffset := offset
	chunkSize := int64(chunkSizeDefault)

	//handling limit
	if limit > 0 && limit < chunkSize {
		chunkSize = limit
	}
	if limit > 0 && limit+offset < sz {
		sz = limit + offset
	}
	var written int64
	count := int(math.Ceil(float64(sz) / float64(chunkSize)))
	bar := pb.StartNew(count)
	for {
		if sz-currentOffset < chunkSize {
			chunkSize = sz - currentOffset
		}

		f.Seek(currentOffset, 0)
		writeOffset := currentOffset - offset
		toFile.Seek(writeOffset, 0)
		wr, errW := io.CopyN(toFile, f, chunkSize)
		if errW != nil {
			return errW
		}
		written = written + wr
		bar.Increment()
		if written >= sz-currentOffset {
			break
		}
		currentOffset += chunkSize
	}

	return nil
}

func validateFile(f *os.File) (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	if stat.IsDir() {
		// TODO may be separate error for this
		return 0, ErrUnsupportedFile
	}
	sz := stat.Size()

	if offset > sz {
		return sz, ErrOffsetExceedsFileSize
	}
	return sz, nil
}
