package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSedBasicSubstitute(t *testing.T) {
	content := "hello world\nfoo bar\nhello again\n"
	writeTestFile(t, "test_sed_basic.txt", content)
	defer os.Remove("test_sed_basic.txt")

	cmd := exec.Command("./gobox", "sed", "s/hello/hi/", "test_sed_basic.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hi world") {
		t.Errorf("Expected 'hi world' in output, got: %s", result)
	}
	if !strings.Contains(result, "hi again") {
		t.Errorf("Expected 'hi again' in output, got: %s", result)
	}
	if strings.Contains(result, "hello") {
		t.Errorf("Unexpected 'hello' in output: %s", result)
	}
}

func TestSedGlobalReplace(t *testing.T) {
	content := "foo foo foo\nbar baz\nfoo\n"
	writeTestFile(t, "test_sed_global.txt", content)
	defer os.Remove("test_sed_global.txt")

	cmd := exec.Command("./gobox", "sed", "s/foo/X/g", "test_sed_global.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "X X X") {
		t.Errorf("Expected 'X X X' in output, got: %s", result)
	}
	if strings.Contains(result, "foo") {
		t.Errorf("Unexpected 'foo' in output: %s", result)
	}
}

func TestSedIgnoreCase(t *testing.T) {
	content := "HELLO world\nHello Again\nhello\n"
	writeTestFile(t, "test_sed_case.txt", content)
	defer os.Remove("test_sed_case.txt")

	cmd := exec.Command("./gobox", "sed", "s/hello/hi/i", "test_sed_case.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hi world") {
		t.Errorf("Expected 'hi world' in output, got: %s", result)
	}
	if !strings.Contains(result, "hi Again") {
		t.Errorf("Expected 'hi Again' in output, got: %s", result)
	}
}

func TestSedQuietMode(t *testing.T) {
	content := "hello world\nfoo bar\nhello again\n"
	writeTestFile(t, "test_sed_quiet.txt", content)
	defer os.Remove("test_sed_quiet.txt")

	cmd := exec.Command("./gobox", "sed", "-n", "s/hello/hi/p", "test_sed_quiet.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d: %s", len(lines), result)
	}
	if !strings.Contains(result, "hi world") {
		t.Errorf("Expected 'hi world' in output, got: %s", result)
	}
	if !strings.Contains(result, "hi again") {
		t.Errorf("Expected 'hi again' in output, got: %s", result)
	}
}

func TestSedDelete(t *testing.T) {
	content := "line1\nDELETE_ME\nline3\nDELETE_ME\nline5\n"
	writeTestFile(t, "test_sed_delete.txt", content)
	defer os.Remove("test_sed_delete.txt")

	cmd := exec.Command("./gobox", "sed", "/DELETE_ME/d", "test_sed_delete.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if strings.Contains(result, "DELETE_ME") {
		t.Errorf("Unexpected 'DELETE_ME' in output: %s", result)
	}
	if !strings.Contains(result, "line1") {
		t.Errorf("Expected 'line1' in output, got: %s", result)
	}
	if !strings.Contains(result, "line3") {
		t.Errorf("Expected 'line3' in output, got: %s", result)
	}
}

func TestSedPrint(t *testing.T) {
	content := "line1\nMATCH\nline3\nMATCH\nline5\n"
	writeTestFile(t, "test_sed_print.txt", content)
	defer os.Remove("test_sed_print.txt")

	cmd := exec.Command("./gobox", "sed", "-n", "/MATCH/p", "test_sed_print.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "MATCH") {
		t.Errorf("Expected 'MATCH' in output, got: %s", result)
	}
	if strings.Contains(result, "line1") {
		t.Errorf("Unexpected 'line1' in output: %s", result)
	}
}

func TestSedInPlace(t *testing.T) {
	content := "old value\nkeep this\nold again\n"
	filename := "test_sed_inplace.txt"
	writeTestFile(t, filename, content)
	defer os.Remove(filename)
	defer os.Remove(filename + ".bak")

	cmd := exec.Command("./gobox", "sed", "-i.bak", "s/old/new/", filename)
	if err := cmd.Run(); err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	// Read modified file
	modified, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("cannot read modified file: %v", err)
	}

	result := string(modified)
	if !strings.Contains(result, "new value") {
		t.Errorf("Expected 'new value' in modified file, got: %s", result)
	}
	if !strings.Contains(result, "new again") {
		t.Errorf("Expected 'new again' in modified file, got: %s", result)
	}

	// Check backup exists
	backup, err := os.ReadFile(filename + ".bak")
	if err != nil {
		t.Fatalf("backup file not created: %v", err)
	}

	if !strings.Contains(string(backup), "old value") {
		t.Errorf("Backup should contain original content")
	}
}

func TestSedMultipleExpressions(t *testing.T) {
	content := "foo bar\nbaz qux\n"
	writeTestFile(t, "test_sed_multi.txt", content)
	defer os.Remove("test_sed_multi.txt")

	cmd := exec.Command("./gobox", "sed", "-e", "s/foo/FOO/", "-e", "s/qux/QUX/", "test_sed_multi.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "FOO bar") {
		t.Errorf("Expected 'FOO bar' in output, got: %s", result)
	}
	if !strings.Contains(result, "baz QUX") {
		t.Errorf("Expected 'baz QUX' in output, got: %s", result)
	}
}

func TestSedStdin(t *testing.T) {
	cmd := exec.Command("./gobox", "sed", "s/foo/bar/")
	cmd.Stdin = strings.NewReader("hello foo world\nfoo again\n")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "hello bar world") {
		t.Errorf("Expected 'hello bar world' in output, got: %s", result)
	}
	if !strings.Contains(result, "bar again") {
		t.Errorf("Expected 'bar again' in output, got: %s", result)
	}
}

func TestSedBackreference(t *testing.T) {
	content := "John Doe\nJane Smith\n"
	writeTestFile(t, "test_sed_backref.txt", content)
	defer os.Remove("test_sed_backref.txt")

	// Go regex uses ${1}, ${2} for backreferences, but we support \1, \2 syntax
	cmd := exec.Command("./gobox", "sed", `s/([A-Za-z]+) ([A-Za-z]+)/${2}, ${1}/`, "test_sed_backref.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "Doe, John") {
		t.Errorf("Expected 'Doe, John' in output, got: %s", result)
	}
	if !strings.Contains(result, "Smith, Jane") {
		t.Errorf("Expected 'Smith, Jane' in output, got: %s", result)
	}
}

func TestSedNthReplacement(t *testing.T) {
	content := "foo foo foo\nbar foo baz\n"
	writeTestFile(t, "test_sed_nth.txt", content)
	defer os.Remove("test_sed_nth.txt")

	cmd := exec.Command("./gobox", "sed", "s/foo/X/2", "test_sed_nth.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	// First line: foo X foo (2nd occurrence replaced)
	if !strings.Contains(lines[0], "foo X foo") {
		t.Errorf("Expected 'foo X foo', got: %s", lines[0])
	}
	// Second line: bar X baz (only one foo, so 2nd doesn't exist, no change)
	if !strings.Contains(lines[1], "bar foo baz") {
		t.Errorf("Expected 'bar foo baz', got: %s", lines[1])
	}
}

func TestSedPrintLineNumber(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_linenum.txt", content)
	defer os.Remove("test_sed_linenum.txt")

	cmd := exec.Command("./gobox", "sed", "=", "test_sed_linenum.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "1") {
		t.Errorf("Expected line number 1 in output, got: %s", result)
	}
	if !strings.Contains(result, "2") {
		t.Errorf("Expected line number 2 in output, got: %s", result)
	}
	if !strings.Contains(result, "3") {
		t.Errorf("Expected line number 3 in output, got: %s", result)
	}
}

func TestSedInsert(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_insert.txt", content)
	defer os.Remove("test_sed_insert.txt")

	cmd := exec.Command("./gobox", "sed", "/line2/i\\INSERTED", "test_sed_insert.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "INSERTED") {
		t.Errorf("Expected 'INSERTED' in output, got: %s", result)
	}
	// INSERTED should come before line2
	idxInserted := strings.Index(result, "INSERTED")
	idxLine2 := strings.Index(result, "line2")
	if idxInserted >= idxLine2 {
		t.Errorf("INSERTED should come before line2, got: %s", result)
	}
}

func TestSedInsertNumeric(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_insert_num.txt", content)
	defer os.Remove("test_sed_insert_num.txt")

	cmd := exec.Command("./gobox", "sed", "2i\\BEFORE_LINE_2", "test_sed_insert_num.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "BEFORE_LINE_2") {
		t.Errorf("Expected 'BEFORE_LINE_2' in output, got: %s", result)
	}
}

func TestSedAppend(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_append.txt", content)
	defer os.Remove("test_sed_append.txt")

	cmd := exec.Command("./gobox", "sed", "/line2/a\\APPENDED", "test_sed_append.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "APPENDED") {
		t.Errorf("Expected 'APPENDED' in output, got: %s", result)
	}
	// APPENDED should come after line2
	idxLine2 := strings.Index(result, "line2")
	idxAppended := strings.Index(result, "APPENDED")
	if idxAppended <= idxLine2 {
		t.Errorf("APPENDED should come after line2, got: %s", result)
	}
}

func TestSedChange(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_change.txt", content)
	defer os.Remove("test_sed_change.txt")

	cmd := exec.Command("./gobox", "sed", "/line2/c\\CHANGED", "test_sed_change.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "CHANGED") {
		t.Errorf("Expected 'CHANGED' in output, got: %s", result)
	}
	if strings.Contains(result, "line2") {
		t.Errorf("line2 should be replaced, got: %s", result)
	}
}

func TestSedChangeNumeric(t *testing.T) {
	content := "line1\nline2\nline3\n"
	writeTestFile(t, "test_sed_change_num.txt", content)
	defer os.Remove("test_sed_change_num.txt")

	cmd := exec.Command("./gobox", "sed", "2c\\REPLACED_LINE_2", "test_sed_change_num.txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("sed command failed: %v", err)
	}

	result := string(output)
	if !strings.Contains(result, "REPLACED_LINE_2") {
		t.Errorf("Expected 'REPLACED_LINE_2' in output, got: %s", result)
	}
	if strings.Contains(result, "line2") {
		t.Errorf("line2 should be replaced, got: %s", result)
	}
}
