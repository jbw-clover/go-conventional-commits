package conventionalcommits

import (
	"encoding/json"
	"fmt"
)

var Marshal = json.Marshal

// ConventionalCommits a slice of parsed conventional commit messages
type ConventionalCommits []*ConventionalCommit

// ConventionalCommit a parsed conventional commit message
type ConventionalCommit struct {
	Issues      []string `json:"issues,emitempty"`
	Category    string   `json:"category"`
	Scope       string   `json:"scope"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	Footer      []string `json:"footer"`
	Major       bool     `json:"major"`
	Minor       bool     `json:"minor"`
	Patch       bool     `json:"patch"`
}

func (cc *ConventionalCommit) String() string {
	data, err := Marshal(cc)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(data)
}

func (cc ConventionalCommits) IsMajor() bool {
	for _, c := range cc {
		if c.Major {
			return true
		}
	}
	return false
}

func (cc ConventionalCommits) IsMinor() bool {
	for _, c := range cc {
		if c.Minor {
			return true
		}
	}
	return false
}
