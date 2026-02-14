package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(from)
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileInfo.Size() - offset
	}

	reader, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if _, err := reader.Seek(offset, 0); err != nil {
		return ErrUnsupportedFile
	}

	limitReader := io.LimitReader(reader, limit)
	defer reader.Close()
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(limitReader)

	writer, err := os.OpenFile(toPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}

	defer writer.Close()

	io.Copy(writer, barReader)

	// finish bar
	bar.Finish()
	return nil
}
