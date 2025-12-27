package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// BrokenLink represents a link that failed validation
type BrokenLink struct {
	File   string
	Line   int
	URL    string
	Reason string
}

var (
	// Regex to match markdown links: [text](url)
	mdLinkRegex = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	// Regex to match image links: ![alt](url)
	imgLinkRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

// StatusDescriptions maps HTTP status codes to descriptions
var StatusDescriptions = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

func main() {
	var rootDir string
	var timeout int
	var concurrency int
	var checkExternal bool
	var verbose bool

	flag.StringVar(&rootDir, "root", ".", "Root directory to scan for markdown files")
	flag.IntVar(&timeout, "timeout", 10, "HTTP request timeout in seconds")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of concurrent HTTP requests")
	flag.BoolVar(&checkExternal, "external", true, "Check external URLs")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.Parse()

	// Convert to absolute path
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving root directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Scanning for markdown files in %s...\n\n", absRoot)

	// Find all markdown files
	mdFiles, err := findMarkdownFiles(absRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding markdown files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d markdown files.\n\n", len(mdFiles))

	// Extract all links
	allLinks := make(map[string][]LinkInfo) // URL -> list of (file, line) where it appears
	for _, mdFile := range mdFiles {
		links, err := extractLinksFromFile(mdFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not read %s: %v\n", mdFile, err)
			continue
		}
		for _, link := range links {
			allLinks[link.URL] = append(allLinks[link.URL], LinkInfo{File: mdFile, Line: link.Line})
		}
	}

	fmt.Printf("Found %d unique links to check.\n\n", len(allLinks))

	// Check links
	brokenLinks := checkLinks(absRoot, allLinks, timeout, concurrency, checkExternal, verbose)

	// Print results
	fmt.Println()
	if len(brokenLinks) == 0 {
		fmt.Println("No broken links found.")
		os.Exit(0)
	}

	fmt.Printf("Found %d broken links:\n\n", len(brokenLinks))
	for _, bl := range brokenLinks {
		relFile, _ := filepath.Rel(absRoot, bl.File)
		fmt.Printf("[BROKEN] %s:%d\n  URL: %s\n  Reason: %s\n\n", relFile, bl.Line, bl.URL, bl.Reason)
	}

	os.Exit(1)
}

// LinkInfo stores information about where a link was found
type LinkInfo struct {
	File string
	Line int
	URL  string
}

// findMarkdownFiles recursively finds all .md files in the given directory
func findMarkdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden directories and common non-doc directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		// Only include .md files
		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// extractLinksFromFile extracts all links from a markdown file
func extractLinksFromFile(filePath string) ([]LinkInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []LinkInfo
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip code blocks (simple detection)
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			continue
		}

		// Extract markdown links
		for _, match := range mdLinkRegex.FindAllStringSubmatch(line, -1) {
			if len(match) >= 3 {
				url := strings.TrimSpace(match[2])
				// Skip mailto links and empty links
				if !strings.HasPrefix(url, "mailto:") && url != "" {
					// Handle links with titles: [text](url "title")
					if idx := strings.Index(url, " "); idx != -1 {
						url = url[:idx]
					}
					if idx := strings.Index(url, "\t"); idx != -1 {
						url = url[:idx]
					}
					links = append(links, LinkInfo{URL: url, Line: lineNum})
				}
			}
		}

		// Extract image links
		for _, match := range imgLinkRegex.FindAllStringSubmatch(line, -1) {
			if len(match) >= 3 {
				url := strings.TrimSpace(match[2])
				if url != "" {
					// Handle links with titles
					if idx := strings.Index(url, " "); idx != -1 {
						url = url[:idx]
					}
					if idx := strings.Index(url, "\t"); idx != -1 {
						url = url[:idx]
					}
					links = append(links, LinkInfo{URL: url, Line: lineNum})
				}
			}
		}
	}

	return links, scanner.Err()
}

// checkLinks validates all links and returns broken ones
func checkLinks(rootDir string, links map[string][]LinkInfo, timeout, concurrency int, checkExternal, verbose bool) []BrokenLink {
	var brokenLinks []BrokenLink
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Semaphore for concurrency control
	sem := make(chan struct{}, concurrency)

	for url, locations := range links {
		wg.Add(1)
		go func(url string, locations []LinkInfo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Get first location for this URL
			firstLoc := locations[0]

			// Skip anchor-only links
			if strings.HasPrefix(url, "#") {
				if verbose {
					fmt.Printf("[SKIP] %s (anchor link)\n", url)
				}
				return
			}

			var reason string

			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				// External URL
				if !checkExternal {
					if verbose {
						fmt.Printf("[SKIP] %s (external check disabled)\n", url)
					}
					return
				}
				reason = checkExternalURL(client, url, verbose)
			} else {
				// Relative path
				reason = checkRelativePath(rootDir, firstLoc.File, url, verbose)
			}

			if reason != "" {
				mu.Lock()
				for _, loc := range locations {
					brokenLinks = append(brokenLinks, BrokenLink{
						File:   loc.File,
						Line:   loc.Line,
						URL:    url,
						Reason: reason,
					})
				}
				mu.Unlock()
			}
		}(url, locations)
	}

	wg.Wait()
	return brokenLinks
}

// checkExternalURL checks if an external URL is accessible
func checkExternalURL(client *http.Client, url string, verbose bool) string {
	if verbose {
		fmt.Printf("[CHECK] %s\n", url)
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	// Try HEAD request first
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return fmt.Sprintf("Invalid URL: %v", err)
	}

	// Set a user agent to avoid being blocked by some servers
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkChecker/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		// Try GET if HEAD fails
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkChecker/1.0)")
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Sprintf("Request failed: %v", err)
		}
	}
	defer resp.Body.Close()

	// Consider 4xx and 5xx as broken (except 405 which some servers return for HEAD)
	if resp.StatusCode >= 400 && resp.StatusCode != 405 {
		desc := StatusDescriptions[resp.StatusCode]
		if desc == "" {
			desc = "Error"
		}
		return fmt.Sprintf("%d %s", resp.StatusCode, desc)
	}

	return ""
}

// checkRelativePath checks if a relative file path exists
func checkRelativePath(rootDir, mdFile, relPath string, verbose bool) string {
	if verbose {
		fmt.Printf("[CHECK] %s (relative to %s)\n", relPath, mdFile)
	}

	// Remove any anchor from the path
	cleanPath := relPath
	if idx := strings.Index(cleanPath, "#"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	// Skip empty paths (pure anchor references)
	if cleanPath == "" {
		return ""
	}

	// Get the directory of the markdown file
	mdDir := filepath.Dir(mdFile)

	// Resolve the relative path
	fullPath := filepath.Join(mdDir, cleanPath)

	// Check if file/directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "File not found"
	}

	return ""
}
