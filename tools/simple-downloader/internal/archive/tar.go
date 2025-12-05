package archive

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

// extractGzipTar extracts a .tar.gz archive
func extractGzipTar(path string, opts ExtractOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	return extractTar(gzr, opts)
}

// extractBzip2Tar extracts a .tar.bz2 archive
func extractBzip2Tar(path string, opts ExtractOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	bzr := bzip2.NewReader(f)
	return extractTar(bzr, opts)
}

// extractXzTar extracts a .tar.xz archive
func extractXzTar(path string, opts ExtractOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	xzr, err := xz.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	return extractTar(xzr, opts)
}

// extractZstdTar extracts a .tar.zstd archive
func extractZstdTar(path string, opts ExtractOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	zstdr, err := zstd.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create zstd reader: %w", err)
	}
	defer zstdr.Close()

	return extractTar(zstdr, opts)
}


