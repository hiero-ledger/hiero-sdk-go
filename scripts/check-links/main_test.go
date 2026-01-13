package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveInlineCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no inline code",
			input:    "This is a normal line",
			expected: "This is a normal line",
		},
		{
			name:     "single inline code",
			input:    "Use `go run main.go` to start",
			expected: "Use  to start",
		},
		{
			name:     "multiple inline codes",
			input:    "Run `cmd1` and `cmd2` together",
			expected: "Run  and  together",
		},
		{
			name:     "inline code with URL",
			input:    "See `[link](http://example.com)` for details",
			expected: "See  for details",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only backticks",
			input:    "``",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeInlineCode(tt.input)
			if result != tt.expected {
				t.Errorf("removeInlineCode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractLinksFromFile(t *testing.T) {
	// Create a temporary markdown file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	content := `# Test Document

This is a [normal link](https://example.com).

Here's an image: ![alt text](./image.png)

` + "```" + `go
// This link should be ignored: [code link](http://ignored.com)
` + "```" + `

Another [link](./relative/path.md) here.

Inline code with link: ` + "`[ignored](http://skip.me)`" + `

[Link with title](https://example.org "Title")

mailto should be skipped: [email](mailto:test@example.com)
`

	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	links, err := extractLinksFromFile(testFile)
	if err != nil {
		t.Fatalf("extractLinksFromFile failed: %v", err)
	}

	// Expected links (excluding code blocks, inline code, and mailto)
	expectedURLs := map[string]bool{
		"https://example.com": true,
		"./image.png":         true,
		"./relative/path.md":  true,
		"https://example.org": true,
	}

	// Links that should NOT be extracted
	excludedURLs := map[string]bool{
		"http://ignored.com":      true,
		"http://skip.me":          true,
		"mailto:test@example.com": true,
	}

	for _, link := range links {
		if excludedURLs[link.URL] {
			t.Errorf("extractLinksFromFile extracted excluded URL: %s", link.URL)
		}
		delete(expectedURLs, link.URL)
	}

	for url := range expectedURLs {
		t.Errorf("extractLinksFromFile missed expected URL: %s", url)
	}
}

func TestFindMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"docs",
		"src",
		".git",
		"node_modules",
		"vendor",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create markdown files
	mdFiles := []string{
		"README.md",
		"docs/guide.md",
		"src/notes.MD",        // uppercase extension
		".git/config.md",      // should be skipped
		"node_modules/pkg.md", // should be skipped
		"vendor/lib.md",       // should be skipped
	}

	for _, f := range mdFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("# Test"), 0600); err != nil {
			t.Fatalf("Failed to create file %s: %v", f, err)
		}
	}

	// Create a non-markdown file
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0600); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	files, err := findMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatalf("findMarkdownFiles failed: %v", err)
	}

	// Should find: README.md, docs/guide.md, src/notes.MD
	// Should NOT find: .git/config.md, node_modules/pkg.md, vendor/lib.md, main.go
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("findMarkdownFiles found %d files, want %d", len(files), expectedCount)
		for _, f := range files {
			t.Logf("  Found: %s", f)
		}
	}

	// Verify hidden/excluded directories are skipped
	for _, f := range files {
		rel, _ := filepath.Rel(tmpDir, f)
		if strings.HasPrefix(rel, ".git") ||
			strings.HasPrefix(rel, "node_modules") ||
			strings.HasPrefix(rel, "vendor") {
			t.Errorf("findMarkdownFiles should have skipped: %s", rel)
		}
	}
}

func TestCheckRelativePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs directory: %v", err)
	}

	testMD := filepath.Join(docsDir, "test.md")
	if err := os.WriteFile(testMD, []byte("# Test"), 0600); err != nil {
		t.Fatalf("Failed to create test.md: %v", err)
	}

	existingFile := filepath.Join(docsDir, "existing.md")
	if err := os.WriteFile(existingFile, []byte("# Existing"), 0600); err != nil {
		t.Fatalf("Failed to create existing.md: %v", err)
	}

	tests := []struct {
		name      string
		relPath   string
		wantError bool
	}{
		{
			name:      "existing file",
			relPath:   "./existing.md",
			wantError: false,
		},
		{
			name:      "non-existing file",
			relPath:   "./missing.md",
			wantError: true,
		},
		{
			name:      "anchor only",
			relPath:   "#section",
			wantError: false, // Should be handled by checkLinks, not here
		},
		{
			name:      "file with anchor",
			relPath:   "./existing.md#section",
			wantError: false,
		},
		{
			name:      "missing file with anchor",
			relPath:   "./missing.md#section",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason := checkRelativePath(tmpDir, testMD, tt.relPath, false)
			hasError := reason != ""
			if hasError != tt.wantError {
				t.Errorf("checkRelativePath(%q) returned %q, wantError=%v", tt.relPath, reason, tt.wantError)
			}
		})
	}
}

func TestStatusDescriptions(t *testing.T) {
	// Verify common status codes have descriptions
	expectedCodes := []int{400, 401, 403, 404, 500, 502, 503, 504}

	for _, code := range expectedCodes {
		if desc, ok := StatusDescriptions[code]; !ok || desc == "" {
			t.Errorf("StatusDescriptions missing or empty for code %d", code)
		}
	}
}
