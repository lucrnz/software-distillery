package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"simple-downloader/internal/archive"
	"simple-downloader/internal/downloader"
	"simple-downloader/internal/version"
)

var (
	urlStr          string
	output          string
	quiet           bool
	expectedHash    string
	extractArchive  bool
	removeArchive   bool
	chdir           string
	stripComponents int
	connectTimeout  time.Duration
	maxTime         time.Duration
	userAgent       string
)

var rootCmd = &cobra.Command{
	Use:   "simple-downloader",
	Short: "Your Swiss-Army Knife for downloading files",
	Long: fmt.Sprintf(`Simple Downloader

Downloads a file from URL, following redirects.
Optionally extracts archives (zip, tar, tar.gz, tar.bz2, tar.xz, tar.zstd).

Copyright (c) %d Luciano Hillcoat.
This program is open-source and warranty-free, read more at: https://github.com/lucrnz/software-distillery/blob/main/LICENSE
`, time.Now().Year()),
	RunE:    run,
	Version: version.Print(),
}

func init() {
	rootCmd.Flags().StringVarP(&urlStr, "url", "U", "", "The URL to download (required)")
	rootCmd.Flags().StringVarP(&output, "output", "O", "", "The name for the file to write it as")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Does not show any progress or output")
	rootCmd.Flags().StringVarP(&expectedHash, "hash", "H", "", "User-provided SHA256 hex sum hash to check against the resulting file")
	rootCmd.Flags().BoolVarP(&extractArchive, "extract-archive", "x", false, "Extract the downloaded archive")
	rootCmd.Flags().BoolVar(&removeArchive, "remove-archive", true, "Delete archive file after successful extraction")
	rootCmd.Flags().StringVarP(&chdir, "chdir", "C", "", "Change working directory before any operation (panics if directory doesn't exist)")
	rootCmd.Flags().IntVar(&stripComponents, "extract-strip-components", 0, "Strip N leading components from file names during extraction")
	rootCmd.Flags().DurationVar(&connectTimeout, "connect-timeout", 300*time.Second, "Maximum time for connection establishment")
	rootCmd.Flags().DurationVarP(&maxTime, "max-time", "m", 0, "Maximum total time for the entire operation (0 = unlimited)")
	rootCmd.Flags().StringVar(&userAgent, "user-agent", version.UserAgent(), "User-Agent header to send with HTTP requests")

	rootCmd.MarkFlagRequired("url")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Change directory first if specified
	if chdir != "" {
		if err := os.Chdir(chdir); err != nil {
			return fmt.Errorf("failed to change directory to %q: %w", chdir, err)
		}
	}

	// Normalize URL
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "https://" + urlStr
	}

	// Determine output filename
	if output == "" {
		parsedURL := urlStr
		if idx := strings.LastIndex(parsedURL, "/"); idx != -1 {
			output = parsedURL[idx+1:]
		}
		if output == "" || output == "/" {
			output = "download"
		}
		// Strip query string if present
		if idx := strings.Index(output, "?"); idx != -1 {
			output = output[:idx]
		}
	}

	// Cannot extract when outputting to stdout
	if extractArchive && output == "-" {
		return fmt.Errorf("cannot extract archive when output is stdout (-)")
	}

	// Perform download
	opts := downloader.Options{
		URL:            urlStr,
		Output:         output,
		Quiet:          quiet,
		ExpectedHash:   expectedHash,
		ConnectTimeout: connectTimeout,
		MaxTime:        maxTime,
		UserAgent:      userAgent,
	}

	_, err := downloader.Download(opts)
	if err != nil {
		return err
	}

	// Extract archive if requested
	if extractArchive {
		if !quiet {
			fmt.Printf("Detecting archive type...\n")
		}

		archiveType, err := archive.Detect(output)
		if err != nil {
			return fmt.Errorf("error detecting archive type: %w", err)
		}

		if archiveType == archive.Unknown {
			return fmt.Errorf("unknown or unsupported archive format")
		}

		if !quiet {
			fmt.Printf("Detected archive type: %s\n", archiveType)
			fmt.Printf("Extracting...\n")
		}

		opts := archive.ExtractOptions{
			StripComponents: stripComponents,
		}
		if err := archive.Extract(output, archiveType, opts); err != nil {
			return fmt.Errorf("error extracting archive: %w", err)
		}

		if !quiet {
			fmt.Printf("✅ Extraction complete\n")
		}

		if removeArchive {
			if err := os.Remove(output); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove archive file: %v\n", err)
			} else if !quiet {
				fmt.Printf("Removed archive file: %s\n", output)
			}
		}
	}

	return nil
}
