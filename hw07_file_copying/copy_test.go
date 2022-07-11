package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const inputFile = "./testdata/input.txt"

	fStat, err := os.Stat(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                  string
		giveFrom, giveTo      string
		giveLimit, giveOffset int64
		wantErr               error
		wantFileSize          int64
	}{
		{
			name:    "empty from and to file paths",
			wantErr: ErrEmptyFilePath,
		},
		{
			name:    "empty from file path",
			giveTo:  "/tmp/test.txt",
			wantErr: ErrEmptyFilePath,
		},
		{
			name:     "empty to file path",
			giveFrom: inputFile,
			wantErr:  ErrEmptyFilePath,
		},
		{
			name:       "negative offset",
			giveFrom:   inputFile,
			giveTo:     "/tmp/file.txt",
			giveOffset: -1,
			wantErr:    ErrNegativeOffset,
		},
		{
			name:      "negative limit",
			giveFrom:  inputFile,
			giveTo:    "/tmp/file.txt",
			giveLimit: -1,
			wantErr:   ErrNegativeLimit,
		},
		{
			name:     "unsupported file",
			giveFrom: "/dev/urandom",
			giveTo:   "/tmp/large.txt",
			wantErr:  ErrUnsupportedFile,
		},
		{
			name:     "from directory",
			giveFrom: "/tmp",
			giveTo:   "/tmp/dir.txt",
			wantErr:  ErrUnsupportedFile,
		},
		{
			name:       "offset larger then file size",
			giveFrom:   inputFile,
			giveTo:     "/tmp/file.txt",
			giveOffset: fStat.Size() + 10,
			wantErr:    ErrOffsetExceedsFileSize,
		},
		{
			name:         "limit more than file size",
			giveFrom:     inputFile,
			giveTo:       "/tmp/large_limit.txt",
			giveLimit:    fStat.Size() + fStat.Size(),
			wantFileSize: fStat.Size(),
		},
		{
			name:         "offset and limit are equal to file size",
			giveFrom:     inputFile,
			giveTo:       "./equal_offset_and_limit.txt",
			giveOffset:   fStat.Size(),
			giveLimit:    fStat.Size(),
			wantFileSize: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := os.Remove(tc.giveTo); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			}()

			err := Copy(tc.giveFrom, tc.giveTo, tc.giveOffset, tc.giveLimit)

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)

			stat, err := os.Stat(tc.giveTo)

			require.NoError(t, err)
			require.Equal(t, tc.wantFileSize, stat.Size())
		})
	}
}
