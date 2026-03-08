package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// grepCmd implements a basic subset of grep functionality
func grepCmd(args []string) error {
	fsFlags := flag.NewFlagSet("grep", flag.ContinueOnError)
	ignoreCase := fsFlags.Bool("i", false, "ignore case")
	invert := fsFlags.Bool("v", false, "invert match (show non-matching lines)")
	count := fsFlags.Bool("c", false, "show count of matching lines only")
	lineNumber := fsFlags.Bool("n", false, "show line numbers")
	recursive := fsFlags.Bool("r", false, "recursive search in directories")
	fixedString := fsFlags.Bool("F", false, "interpret pattern as fixed string (not regex)")
	help := fsFlags.Bool("help", false, "show help")

	fsFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gobox grep [OPTIONS] PATTERN [FILE...]")
		fmt.Fprintln(os.Stderr, "Search for PATTERN in each FILE or standard input.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		fsFlags.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  gobox grep \"error\" /var/log/syslog")
		fmt.Fprintln(os.Stderr, "  gobox grep -i -r \"TODO\" /path/to/code")
		fmt.Fprintln(os.Stderr, "  gobox grep -v \"^#\" config.txt")
		fmt.Fprintln(os.Stderr, "  cat file.txt | gobox grep \"pattern\"")
	}

	if err := fsFlags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if *help {
		fsFlags.Usage()
		return nil
	}

	if fsFlags.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "grep: PATTERN is required")
		fsFlags.Usage()
		return fmt.Errorf("pattern required")
	}

	pattern := fsFlags.Arg(0)
	files := fsFlags.Args()[1:]

	// Compile regex or use fixed string matching
	var regex *regexp.Regexp
	if !*fixedString {
		var err error
		if *ignoreCase {
			regex, err = regexp.Compile("(?i)" + pattern)
		} else {
			regex, err = regexp.Compile(pattern)
		}
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	// If no files specified, read from stdin
	if len(files) == 0 {
		if err := grepReader(os.Stdin, pattern, regex, *ignoreCase, *invert, *count, *lineNumber, *fixedString, ""); err != nil {
			return err
		}
		return nil
	}

	// Process files
	totalMatches := 0
	for _, file := range files {
		if *recursive {
			err := filepath.WalkDir(file, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil // Skip errors, continue
				}
				if d.IsDir() {
					return nil
				}
				return grepFile(path, pattern, regex, *ignoreCase, *invert, *count, *lineNumber, *fixedString, &totalMatches)
			})
			if err != nil {
				return err
			}
		} else {
			if err := grepFile(file, pattern, regex, *ignoreCase, *invert, *count, *lineNumber, *fixedString, &totalMatches); err != nil {
				return err
			}
		}
	}

	return nil
}

func grepFile(path, pattern string, regex *regexp.Regexp, ignoreCase, invert, count, lineNumber, fixedString bool, totalMatches *int) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	return grepReader(file, pattern, regex, ignoreCase, invert, count, lineNumber, fixedString, path)
}

func grepReader(r io.Reader, pattern string, regex *regexp.Regexp, ignoreCase, invert, count, lineNumber, fixedString bool, filename string) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	matches := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		var matched bool
		if fixedString {
			// Fixed string matching
			if ignoreCase {
				matched = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
			} else {
				matched = strings.Contains(line, pattern)
			}
		} else {
			// Regex matching
			matched = regex.MatchString(line)
		}

		// Invert match if -v is specified
		if invert {
			matched = !matched
		}

		if matched {
			matches++
			if !count {
				if filename != "" {
					fmt.Printf("%s:", filename)
				}
				if lineNumber {
					fmt.Printf("%d:", lineNum)
				}
				fmt.Println(line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	if count {
		if filename != "" {
			fmt.Printf("%s:%d\n", filename, matches)
		} else {
			fmt.Println(matches)
		}
	}

	return nil
}
