package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var usageStr = `Simple downloader

Downloads a file from URL, following redirects.

Usage: %s [flags]
`

func main() {
	urlStr := ""
	output := ""
	quiet := false
	expectedHash := ""

	flag.StringVar(&urlStr, "url", "", "The URL to download (required)")
	flag.StringVar(&urlStr, "U", "", "")
	flag.StringVar(&output, "output", "", "The name for the file to write it as")
	flag.StringVar(&output, "O", "", "")
	flag.BoolVar(&quiet, "quiet", false, "Does not show any progress or output (disabled by default)")
	flag.BoolVar(&quiet, "q", false, "")
	flag.StringVar(&expectedHash, "hash", "", "User-provided SHA256 hex sum hash to check against the resulting file")
	flag.StringVar(&expectedHash, "H", "", "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageStr, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if urlStr == "" {
		fmt.Fprintln(os.Stderr, "Error: --url is required")
		flag.Usage()
		os.Exit(1)
	}

	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "https://" + urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	if output == "" {
		output = filepath.Base(u.Path)
		if output == "" || output == "/" {
			output = "download"
		}
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: HTTP %s\n", resp.Status)
		os.Exit(1)
	}

	contentLength := resp.ContentLength

	var writer io.Writer
	if output == "-" {
		writer = os.Stdout
	} else {
		file, err := os.Create(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	downloadWithProgress(writer, resp.Body, contentLength, output, quiet, expectedHash)
}

// humanReadableBytes formats bytes to a human-readable string (B, KB, MB, GB)
func humanReadableBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	var exp int = 0
	val := float64(bytes)
	for val >= unit && exp < 3 {
		val /= unit
		exp++
	}
	units := []string{"B", "KB", "MB", "GB"}
	return fmt.Sprintf("%.1f %s", val, units[exp])
}

// downloadWithProgress reads from reader in chunks and writes to writer, showing real-time progress
// throttled to update every 500ms, with optional hash verification
func downloadWithProgress(writer io.Writer, reader io.Reader, total int64, outName string, quiet bool, expectedHash string) {
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
				fmt.Printf("\nError reading: %v\n", err)
				os.Exit(1)
			}
			break
		}
		if hasher != nil {
			hasher.Write(buf[:n])
		}
		n2, err := writer.Write(buf[:n])
		if err != nil || n2 != n {
			fmt.Printf("\nError writing: %v\n", err)
			os.Exit(1)
		}
		downloaded += int64(n)
		if !quiet {
			if time.Since(lastUpdate) >= updateInterval {
				if total == -1 {
					fmt.Printf("\rDownloaded: %s...", humanReadableBytes(downloaded))
				} else {
					percent := float64(downloaded) / float64(total) * 100
					fmt.Printf("\rProgress: %.1f%% (%s/%s)", percent, humanReadableBytes(downloaded), humanReadableBytes(total))
				}
				lastUpdate = time.Now()
			}
		}
	}

	// Hash verification
	if expectedHash != "" {
		sum := hasher.Sum(nil)
		computed := hex.EncodeToString(sum)
		if computed != expectedHash {
			if !quiet {
				fmt.Printf("\n❌ Error: Invalid Sha256 sum\n")
			}
			os.Exit(1)
		}
		if !quiet {
			fmt.Printf("\n✅ Sha256 sum hash matches\n")
		}
	}

	// Final message
	if !quiet {
		sizeStr := humanReadableBytes(downloaded)
		if total != -1 {
			sizeStr = humanReadableBytes(total)
		}
		if outName == "-" {
			fmt.Printf("\nDownloaded %s\n", sizeStr)
		} else {
			fmt.Printf("\nDownloaded %s to %s\n", sizeStr, outName)
		}
	}
}
