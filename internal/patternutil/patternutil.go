package patternutil

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/ErdemOzgen/blackdagger/internal/logger"
)

const rePrefix = "re:"

// MatchOption represents an option for pattern matching
type MatchOption func(*matchOptions)

type matchOptions struct {
	exactMatch bool
	logger     logger.Logger
}

// WithExactMatch configures the matcher to use exact string matching for literal patterns
func WithExactMatch() MatchOption {
	return func(o *matchOptions) {
		o.exactMatch = true
	}
}

// MatchPattern matches content against patterns using either literal or regex matching.
// For files or large content, use MatchPatternScanner instead.
func MatchPattern(content string, patterns []string, opts ...MatchOption) bool {
	scanner := bufio.NewScanner(strings.NewReader(content))
	return MatchPatternScanner(scanner, patterns, opts...)
}

func MatchPatternScanner(scanner *bufio.Scanner, patterns []string, opts ...MatchOption) bool {
	if len(patterns) == 0 {
		return false
	}

	// Apply options
	options := &matchOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var regexps []*regexp.Regexp
	var literalPatterns []string

	// Compile regex patterns first
	for _, pattern := range patterns {
		switch {
		case strings.HasPrefix(pattern, rePrefix):
			re, err := regexp.Compile(strings.TrimPrefix(pattern, rePrefix))
			if err != nil {
				options.logger.Error("invalid regexp pattern", "pattern", pattern, "err", err)
				continue
			}
			regexps = append(regexps, re)
		case strings.HasPrefix(pattern, rePrefix):
			re, err := regexp.Compile(strings.TrimPrefix(pattern, rePrefix))
			if err != nil {
				options.logger.Error("invalid regexp pattern", "pattern", pattern, "err", err)
				continue
			}
			regexps = append(regexps, re)
		default:
			literalPatterns = append(literalPatterns, pattern)
		}
	}

	// Special case: if scanner is empty and we're looking for empty string
	if !scanner.Scan() {
		// Check for empty string patterns
		for _, p := range literalPatterns {
			if p == "" {
				return true
			}
		}
		// Check regex patterns against empty string
		for _, re := range regexps {
			if re.MatchString("") {
				return true
			}
		}
		return false
	}

	// Process first line (already read by scanner.Scan() above)
	line := scanner.Text()
	if matchLine(line, literalPatterns, regexps, options) {
		return true
	}

	// Process remaining lines
	for scanner.Scan() {
		if matchLine(scanner.Text(), literalPatterns, regexps, options) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		options.logger.Error("scanner error", "err", err)
	}

	return false
}

// matchLine checks if a single line matches any of the patterns
func matchLine(line string, literalPatterns []string, regexps []*regexp.Regexp, opts *matchOptions) bool {
	// Check literal patterns
	for _, p := range literalPatterns {
		if opts.exactMatch {
			if line == p {
				return true
			}
		} else {
			if strings.Contains(line, p) {
				return true
			}
		}
	}

	// Check regex patterns
	for _, re := range regexps {
		if re.MatchString(line) {
			return true
		}
	}

	return false
}
