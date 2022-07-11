package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
	ErrEmptyFilePath         = errors.New("empty file path")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateInput(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	if fStat.Mode().IsDir() || !fStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	fw, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	if _, err := f.Seek(offset, 0); err != nil {
		return err
	}

	bar := initBar(fStat.Size(), offset, limit)
	defer bar.Finish()

	if limit == 0 {
		if _, err := io.Copy(bar.NewProxyWriter(fw), f); err != nil {
			return err
		}
	} else {
		if _, err := io.CopyN(bar.NewProxyWriter(fw), f, limit); err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}

	return fw.Sync()
}

func validateInput(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyFilePath
	}

	if offset < 0 {
		return ErrNegativeOffset
	}

	if limit < 0 {
		return ErrNegativeLimit
	}

	return nil
}

func initBar(fileSize, offset, limit int64) *pb.ProgressBar {
	fileSize -= offset

	if limit == 0 || limit > fileSize {
		limit = fileSize
	}

	return pb.Simple.Start64(limit)
}
