package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	infoSrcFile, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if !infoSrcFile.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > infoSrcFile.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || offset+limit > infoSrcFile.Size() {
		limit = infoSrcFile.Size() - offset
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	barReader := bar.NewProxyReader(srcFile)
	_, err = io.CopyN(dstFile, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
