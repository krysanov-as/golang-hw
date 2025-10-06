package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	src := filepath.Join("testdata", "input.txt")

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatalf("cannot stat source file: %v", err)
	}
	srcSize := srcInfo.Size()

	tests := []struct {
		name      string
		src       string
		offset    int64
		limit     int64
		expectErr error
	}{
		{
			name:      "full_copy",
			src:       src,
			offset:    0,
			limit:     0,
			expectErr: nil,
		},
		{
			name:      "copy_with_offset_and_limit",
			src:       src,
			offset:    5,
			limit:     10,
			expectErr: nil,
		},
		{
			name:      "limit_bigger_than_file",
			src:       src,
			offset:    0,
			limit:     999999,
			expectErr: nil,
		},
		{
			name:      "offset_exceeds_file_size",
			src:       src,
			offset:    srcSize + 10,
			limit:     0,
			expectErr: ErrOffsetExceedsFileSize,
		},
		{
			name:      "unsupported_file",
			src:       "testdata",
			offset:    0,
			limit:     0,
			expectErr: ErrUnsupportedFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := filepath.Join(os.TempDir(), "test_"+tt.name)
			defer os.Remove(dst)

			err := Copy(tt.src, dst, tt.offset, tt.limit)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			dstInfo, err := os.Stat(dst)
			if err != nil {
				t.Fatalf("cannot stat destination: %v", err)
			}

			var expectedSize int64
			switch {
			case tt.offset > srcSize:
				expectedSize = 0
			case tt.limit == 0 || tt.offset+tt.limit > srcSize:
				expectedSize = srcSize - tt.offset
			default:
				expectedSize = tt.limit
			}

			if dstInfo.Size() != expectedSize {
				t.Errorf("unexpected file size: got %d, want %d", dstInfo.Size(), expectedSize)
			}
		})
	}
}
