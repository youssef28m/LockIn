package blocker

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/youssef28m/LockIn/internal/validator"
)

// Test domain validation
func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		domain   string
		expected bool
		name     string
	}{
		{"google.com", true, "valid domain"},
		{"github.com", true, "valid domain with dash"},
		{"sub.domain.com", true, "subdomain"},
		{"invalid", false, "no TLD"},
		{"google", false, "single word"},
		{"http://google.com", false, "with protocol"},
		{"google.com/path", false, "with path"},
		{"", false, "empty domain"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validator.IsValidDomain(test.domain)
			if result != test.expected {
				t.Errorf("IsValidDomain(%s) = %v, expected %v", test.domain, result, test.expected)
			}
		})
	}
}

// Test BlockSite - checks if entry is added to hosts file
func TestBlockSite(t *testing.T) {
	// Create a temporary hosts file for testing
	tempHostsPath := "test_hosts.txt"
	defer os.Remove(tempHostsPath)

	// Create test file
	os.WriteFile(tempHostsPath, []byte("127.0.0.1 localhost\n"), 0644)

	// Temporarily override hostsPath
	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	testDomain := "test.example.com"
	entry := "127.0.0.1    " + testDomain

	// Block the site
	err := BlockSite(testDomain)
	if err != nil {
		t.Fatalf("Failed to block site: %v", err)
	}

	// Verify entry was added
	content, _ := os.ReadFile(tempHostsPath)
	if !strings.Contains(string(content), entry) {
		t.Errorf("Expected entry '%s' not found in hosts file", entry)
	}
}

// Test blocking duplicate site
func TestBlockSiteDuplicate(t *testing.T) {
	tempHostsPath := "test_hosts_dup.txt"
	defer os.Remove(tempHostsPath)

	os.WriteFile(tempHostsPath, []byte("127.0.0.1 localhost\n"), 0644)

	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	testDomain := "duplicate.example.com"

	// Block twice
	err := BlockSite(testDomain)
	if err != nil {
		t.Fatalf("First block failed: %v", err)
	}

	err = BlockSite(testDomain)
	if err != nil {
		t.Fatalf("Second block failed: %v", err)
	}

	// Check that it's only added once
	content, _ := os.ReadFile(tempHostsPath)
	count := strings.Count(string(content), "duplicate.example.com")
	if count != 1 {
		t.Errorf("Expected domain to appear once, but appeared %d times", count)
	}
}

// Test UnblockSite
func TestUnblockSite(t *testing.T) {
	tempHostsPath := "test_hosts_unblock.txt"
	defer os.Remove(tempHostsPath)

	testDomain := "unblock.example.com"
	entry := "127.0.0.1    " + testDomain

	os.WriteFile(tempHostsPath, []byte("127.0.0.1 localhost\n"+entry+"\n"), 0644)

	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	// Unblock the site
	err := UnblockSite(testDomain)
	if err != nil {
		t.Fatalf("Failed to unblock site: %v", err)
	}

	// Verify entry was removed
	content, _ := os.ReadFile(tempHostsPath)
	if strings.Contains(string(content), testDomain) {
		t.Errorf("Domain '%s' still found in hosts file after unblock", testDomain)
	}
}

// Test unblocking non-existent site
func TestUnblockNonExistent(t *testing.T) {
	tempHostsPath := "test_hosts_nonexist.txt"
	defer os.Remove(tempHostsPath)

	os.WriteFile(tempHostsPath, []byte("127.0.0.1 localhost\n"), 0644)

	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	// Try to unblock a site that doesn't exist
	err := UnblockSite("nonexistent.example.com")
	if err != nil {
		t.Fatalf("Unblock non-existent should not error: %v", err)
	}
}

// Benchmark tests
func BenchmarkBlockSite(b *testing.B) {
	tempHostsPath := "bench_hosts.txt"
	defer os.Remove(tempHostsPath)

	os.WriteFile(tempHostsPath, []byte("127.0.0.1 localhost\n"), 0644)

	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain := fmt.Sprintf("bench%d.example.com", i)
		BlockSite(domain)
	}
}

func BenchmarkUnblockSite(b *testing.B) {
	tempHostsPath := "bench_hosts_unblock.txt"
	defer os.Remove(tempHostsPath)

	// Setup with many blocked sites
	content := "127.0.0.1 localhost\n"
	for i := 0; i < 100; i++ {
		content += fmt.Sprintf("127.0.0.1    bench%d.example.com\n", i)
	}
	os.WriteFile(tempHostsPath, []byte(content), 0644)

	oldHostsPath := hostsPath
	hostsPath = tempHostsPath
	defer func() { hostsPath = oldHostsPath }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain := fmt.Sprintf("bench%d.example.com", i%100)
		UnblockSite(domain)
	}
}
