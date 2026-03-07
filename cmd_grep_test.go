package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGrepBasicMatch(t *testing.T) {
	// Create test file
	content := "hello world\nfoo bar\nhello again\n"
	writeTestFile(t, "test_basic.txt", content)
	defer os.Remove("test_basic.txt")

	// Run grep
	cmd := exec.Command("./gobox", "grep", "hello", "test_basic.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hello world") {
		t.Errorf("Expected 'hello world' in output, got: %s", result)
	}
	if !strings.Contains(result, "hello again") {
		t.Errorf("Expected 'hello again' in output, got: %s", result)
	}
	if strings.Contains(result, "foo bar") {
		t.Errorf("Unexpected 'foo bar' in output: %s", result)
	}
}

func TestGrepIgnoreCase(t *testing.T) {
	content := "HELLO world\nfoo BAR\nHello Again\n"
	writeTestFile(t, "test_case.txt", content)
	defer os.Remove("test_case.txt")

	cmd := exec.Command("./gobox", "grep", "-i", "hello", "test_case.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "HELLO world") {
		t.Errorf("Expected 'HELLO world' in output, got: %s", result)
	}
	if !strings.Contains(result, "Hello Again") {
		t.Errorf("Expected 'Hello Again' in output, got: %s", result)
	}
}

func TestGrepInvertMatch(t *testing.T) {
	content := "hello world\nfoo bar\nhello again\nbaz qux\n"
	writeTestFile(t, "test_invert.txt", content)
	defer os.Remove("test_invert.txt")

	cmd := exec.Command("./gobox", "grep", "-v", "hello", "test_invert.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if strings.Contains(result, "hello world") {
		t.Errorf("Unexpected 'hello world' in output: %s", result)
	}
	if strings.Contains(result, "hello again") {
		t.Errorf("Unexpected 'hello again' in output: %s", result)
	}
	if !strings.Contains(result, "foo bar") {
		t.Errorf("Expected 'foo bar' in output, got: %s", result)
	}
	if !strings.Contains(result, "baz qux") {
		t.Errorf("Expected 'baz qux' in output, got: %s", result)
	}
}

func TestGrepCount(t *testing.T) {
	content := "hello world\nfoo bar\nhello again\nhello third\n"
	writeTestFile(t, "test_count.txt", content)
	defer os.Remove("test_count.txt")

	cmd := exec.Command("./gobox", "grep", "-c", "hello", "test_count.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := strings.TrimSpace(string(output))
	// Count output includes filename when file is specified
	if result != "test_count.txt:3" && result != "3" {
		t.Errorf("Expected count 3 (with or without filename), got: %s", result)
	}
}

func TestGrepLineNumber(t *testing.T) {
	content := "first line\nsecond line with hello\nthird line\n"
	writeTestFile(t, "test_linenum.txt", content)
	defer os.Remove("test_linenum.txt")

	cmd := exec.Command("./gobox", "grep", "-n", "hello", "test_linenum.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "2:second line with hello") {
		t.Errorf("Expected line number 2 in output, got: %s", result)
	}
}

func TestGrepFixedString(t *testing.T) {
	content := "hello.world\nfoo bar\nhelloXworld\n"
	writeTestFile(t, "test_fixed.txt", content)
	defer os.Remove("test_fixed.txt")

	cmd := exec.Command("./gobox", "grep", "-F", "hello.world", "test_fixed.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hello.world") {
		t.Errorf("Expected 'hello.world' in output, got: %s", result)
	}
	if strings.Contains(result, "helloXworld") {
		t.Errorf("Unexpected 'helloXworld' in output: %s", result)
	}
}

func TestGrepRegex(t *testing.T) {
	content := "test123\nfoo456\ntest789\nbar\n"
	writeTestFile(t, "test_regex.txt", content)
	defer os.Remove("test_regex.txt")

	cmd := exec.Command("./gobox", "grep", "test[0-9]+", "test_regex.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "test123") {
		t.Errorf("Expected 'test123' in output, got: %s", result)
	}
	if !strings.Contains(result, "test789") {
		t.Errorf("Expected 'test789' in output, got: %s", result)
	}
	if strings.Contains(result, "foo456") {
		t.Errorf("Unexpected 'foo456' in output: %s", result)
	}
}

func TestGrepNoMatch(t *testing.T) {
	content := "hello world\nfoo bar\n"
	writeTestFile(t, "test_nomatch.txt", content)
	defer os.Remove("test_nomatch.txt")

	cmd := exec.Command("./gobox", "grep", "notfound", "test_nomatch.txt")
	output, err := cmd.Output()
	if err != nil {
		// Exit code 1 is expected for no matches
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				t.Fatalf("Expected exit code 1, got: %d", exitErr.ExitCode())
			}
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if len(output) != 0 {
		t.Errorf("Expected empty output for no matches, got: %s", string(output))
	}
}

func TestGrepRecursive(t *testing.T) {
	// Create test directory structure
	os.MkdirAll("testdir/subdir", 0755)
	defer os.RemoveAll("testdir")

	writeTestFile(t, "testdir/file1.txt", "hello world\n")
	writeTestFile(t, "testdir/subdir/file2.txt", "hello again\n")
	writeTestFile(t, "testdir/subdir/file3.txt", "goodbye\n")

	cmd := exec.Command("./gobox", "grep", "-r", "hello", "testdir")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hello world") {
		t.Errorf("Expected 'hello world' in output, got: %s", result)
	}
	if !strings.Contains(result, "hello again") {
		t.Errorf("Expected 'hello again' in output, got: %s", result)
	}
	if strings.Contains(result, "goodbye") {
		t.Errorf("Unexpected 'goodbye' in output: %s", result)
	}
}

func TestGrepStdin(t *testing.T) {
	cmd := exec.Command("./gobox", "grep", "test")
	cmd.Stdin = strings.NewReader("hello\ntest line\nworld\n")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("grep command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "test line") {
		t.Errorf("Expected 'test line' in output, got: %s", result)
	}
}

// Helper function to write test files
func writeTestFile(t *testing.T, filename, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file %s: %v", filename, err)
	}
}
