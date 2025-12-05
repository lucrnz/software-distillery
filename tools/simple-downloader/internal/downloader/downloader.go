package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"simple-downloader/internal/util"
)

// Options configures the download behavior
type Options struct {
	URL            string
	Output         string // Output file path, or "-" for stdout
	Quiet          bool
	ExpectedHash   string        // SHA256 hex string to verify against
	ConnectTimeout time.Duration // Maximum time for connection establishment
	MaxTime        time.Duration // Maximum total time for the entire operation (0 = unlimited)
	UserAgent      string        // User-Agent header to send with HTTP requests
}

// Result contains the outcome of a download
type Result struct {
	BytesDownloaded int64
	HashMatched     bool
}

// Download fetches a URL and writes it to the specified output
func Download(opts Options) (*Result, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: opts.ConnectTimeout,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
	}

	if opts.MaxTime > 0 {
		client.Timeout = opts.MaxTime
	}

	req, err := http.NewRequest("GET", opts.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}

	var writer io.Writer
	if opts.Output == "-" {
		writer = os.Stdout
	} else {
		file, err := os.Create(opts.Output)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %w", err)
		}
		defer file.Close()
		writer = file
	}

	return downloadWithProgress(writer, resp.Body, resp.ContentLength, opts.Output, opts.Quiet, opts.ExpectedHash)
}

// downloadWithProgress reads from reader in chunks and writes to writer, showing real-time progress
// throttled to update every 500ms, with optional hash verification
func downloadWithProgress(writer io.Writer, reader io.Reader, total int64, outName string, quiet bool, expectedHash string) (*Result, error) {
	updateInterval := 500 * time.Millisecond
	lastUpdate := time.Now()
	var downloaded int64
	buf := make([]byte, 4096)

	var hasher hash.Hash
	if expectedHash != "" {
		hasher = sha256.New()
	}

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("error reading: %w", err)
			}
			break
		}
		if hasher != nil {
			hasher.Write(buf[:n])
		}
		n2, err := writer.Write(buf[:n])
		if err != nil || n2 != n {
			return nil, fmt.Errorf("error writing: %w", err)
		}
		downloaded += int64(n)
		if !quiet {
			if time.Since(lastUpdate) >= updateInterval {
				if total == -1 {
					fmt.Fprintf(os.Stderr, "\rDownloaded: %s...", util.HumanReadableBytes(downloaded))
				} else {
					percent := float64(downloaded) / float64(total) * 100
					fmt.Fprintf(os.Stderr, "\rProgress: %.1f%% (%s/%s)", percent, util.HumanReadableBytes(downloaded), util.HumanReadableBytes(total))
				}
				lastUpdate = time.Now()
			}
		}
	}

	result := &Result{
		BytesDownloaded: downloaded,
		HashMatched:     true,
	}

	// Hash verification
	if expectedHash != "" {
		sum := hasher.Sum(nil)
		computed := hex.EncodeToString(sum)
		if computed != expectedHash {
			result.HashMatched = false
			if !quiet {
				fmt.Printf("\n❌ error: invalid SHA-256 sum\n")
			}
			return result, fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, computed)
		}
		if !quiet {
			fmt.Printf("\n✅ SHA-256 sum hash matches\n")
		}
	}

	// Final message
	if !quiet {
		sizeStr := util.HumanReadableBytes(downloaded)
		if total != -1 {
			sizeStr = util.HumanReadableBytes(total)
		}
		if outName == "-" {
			fmt.Printf("\nDownloaded %s\n", sizeStr)
		} else {
			fmt.Printf("\nDownloaded %s to %s\n", sizeStr, outName)
		}
	}

	return result, nil
}
