package transformers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jbw-clover/go-conventional-commits/pkg/changelogs"
	"github.com/jbw-clover/go-conventional-commits/pkg/common"
	"github.com/jbw-clover/go-conventional-commits/pkg/conventionalcommits"
	"github.com/tsuyoshiwada/go-gitlog"
)

var (
	baseFormatRegex       = regexp.MustCompile(`(?is)^(?:(?P<category>[^\(!:]+)(?:\((?P<scope>[^\)]+)\))?(?P<breaking>!)?: (?P<description>[^\n\r]+))(?P<remainder>.*)`)
	bodyFooterFormatRegex = regexp.MustCompile(`(?isU)^(?:(?P<body>.*))?(?P<footer>(?-U:(?:[\w\-]+(?:: | #).*|(?i:BREAKING CHANGE:.*))+))`)
	footerFormatRegex     = regexp.MustCompile(`(?s)^(?P<footer>(?i:(?:[\w\-]+(?:: | #).*|(?i:BREAKING CHANGE:.*))+))`)
)

// TransformConventionalCommit takes a commits message and parses it into usable blocks
func TransformConventionalCommit(message string, issueExtractor func(string) ([]string, string)) (commit *conventionalcommits.ConventionalCommit) {
	parts := strings.SplitN(message, "\n", 2)
	parts = append(parts, "")

	issues, msgRemainder := issueExtractor(parts[0])
	msgBody := parts[1]

	match := baseFormatRegex.FindStringSubmatch(strings.Join([]string{msgRemainder, msgBody}, "\n"))
	if len(match) == 0 {
		return &conventionalcommits.ConventionalCommit{
			Issues:      issues,
			Category:    "chore",
			Major:       strings.Contains(msgBody, "BREAKING CHANGE"),
			Description: msgRemainder,
			Body:        strings.TrimSpace(msgBody),
		}
	}

	result := make(map[string]string)
	regExMapper(match, baseFormatRegex, result)

	// split the remainder into body & footer
	match = bodyFooterFormatRegex.FindStringSubmatch(result["remainder"])
	if len(match) > 0 {
		regExMapper(match, bodyFooterFormatRegex, result)
	} else {
		result["body"] = result["remainder"]
	}

	for _, category := range common.MajorCategories {
		if result["category"] == category {
			result["breaking"] = "!"
			break
		}
	}

	var footers []string
	for _, v := range strings.Split(result["footer"], "\n") {
		// v = strings.TrimSpace(v)
		if !footerFormatRegex.MatchString(v) && len(footers) > 0 {
			footers[len(footers)-1] += fmt.Sprintf("\n%s", v)
			continue
		}
		footers = append(footers, v)
	}
	for i := range footers {
		footers[i] = strings.TrimSpace(footers[i])
		if footers[i] == "" { // Remove the element at index i from footers.
			copy(footers[i:], footers[i+1:])   // Shift a[i+1:] left one index.
			footers[len(footers)-1] = ""       // Erase last element (write zero value).
			footers = footers[:len(footers)-1] // Truncate slice.
		}
	}
	if len(footers) == 0 {
		footers = nil
	}

	commit = &conventionalcommits.ConventionalCommit{
		Issues:      issues,
		Category:    result["category"],
		Scope:       result["scope"],
		Major:       result["breaking"] == "!" || strings.Contains(result["footer"], "BREAKING CHANGE"),
		Description: result["description"],
		Body:        result["body"],
		Footer:      footers,
	}

	if commit.Major {
		return commit
	}

	for _, category := range common.MinorCategories {
		if result["category"] == category {
			commit.Minor = true
			return commit
		}
	}

	for _, category := range common.PatchCategories {
		if result["category"] == category {
			commit.Patch = true
			return commit
		}
	}

	return commit
}

func NullIssuesParser(message string) ([]string, string) {
	return make([]string, 0), message
}

func TransformConventionalCommits(messages []string) (commits conventionalcommits.ConventionalCommits) {
	for _, message := range messages {
		commits = append(commits, TransformConventionalCommit(message, NullIssuesParser))
	}
	return
}

// TransformChangeLog takes a commits message and parses it into change log blocks
func TransformChangeLog(message, projectLink string) *changelogs.ChangeLog {
	commit := TransformConventionalCommit(message, NullIssuesParser)

	desc := commit.Description
	ref := ""
	footerTitle := ""
	link := ""

	for _, footer := range commit.Footer {
		fr := footerByKey(footer, "Refs")
		if fr != "" {
			ref = fr
		}

		ft := footerByKey(footer, "Title")
		if ft != "" {
			footerTitle = ft
		}
	}

	if ref == "" {
		descParts := strings.Split(desc, " #")
		desc = descParts[0]

		if len(descParts) > 1 {
			ref = "#" + descParts[1]
		}
	} else {
		link = ref
		if projectLink != "" {
			link = fmt.Sprintf("[%s](%s%s)", ref, projectLink, ref)
		}
	}

	title := desc
	if footerTitle != "" {
		title = footerTitle
	}

	category := common.Changes

	switch {
	case strings.Contains(commit.Category, "fix"):
		category = common.Fixes
	case strings.Contains(commit.Category, "feat"):
		category = common.Features
	}

	return &changelogs.ChangeLog{
		Category: category,
		Refs:     ref,
		Title:    title,
		Link:     link,
	}
}

func TransformChangeLogs(messages []string, projectLink string) changelogs.ChangeLogs {
	parsedChangelogs := changelogs.ChangeLogs{}

	for _, message := range messages {
		changeLog := TransformChangeLog(message, projectLink)
		if changeLog != nil {
			parsedChangelogs[changeLog.Refs] = changeLog
		}
	}

	return parsedChangelogs
}

// TransformChangeLog takes a commits message and parses it into a slice of string
func TransformMessages(commits []*gitlog.Commit, commitsURL string) []string {
	messages := make([]string, 0, len(commits))

	for _, commit := range commits {
		hash := ""
		if commit.Hash != nil && len(commit.Hash.Short) > 0 {
			hash = " #" + commitsURL + commit.Hash.Short
		}

		message := commit.Subject + hash + "\n\n" + commit.Body
		messages = append(messages, message)
	}

	return messages
}

func regExMapper(match []string, expectedFormatRegex *regexp.Regexp, result map[string]string) {
	for i, name := range expectedFormatRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = strings.TrimSpace(match[i])
		}
	}
}

// TODO: Move to other package.
func footerByKey(footer, key string) string {
	result := ""
	footerLower := strings.ToLower(footer)
	keyLower := strings.ToLower(key)
	if strings.Contains(footerLower, fmt.Sprintf("%s #", keyLower)) {
		result = strings.Split(footer, "#")[1]
	}
	if strings.Contains(footerLower, fmt.Sprintf("%s: ", keyLower)) {
		result = strings.Split(footer, ": ")[1]
	}
	return result
}
