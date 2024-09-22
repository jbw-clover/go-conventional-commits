package conventionalcommits_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/jbw-clover/go-conventional-commits/pkg/conventionalcommits"
	"github.com/stretchr/testify/require"
)

func TestConventionalCommit_String(t *testing.T) {
	emptyCommit := &conventionalcommits.ConventionalCommit{Issues: []string{}}

	tests := []struct {
		name               string
		conventionalCommit *conventionalcommits.ConventionalCommit
		prepare            func()
		expected           string
	}{
		{
			name:               "default",
			conventionalCommit: emptyCommit,
			prepare:            func() {},
			expected:           `{"issues":[],"category":"","scope":"","description":"","body":"","footer":null,"major":false,"minor":false,"patch":false}`,
		},
		{
			name:               "marshal error",
			conventionalCommit: emptyCommit,
			prepare: func() {
				conventionalcommits.Marshal = func(v interface{}) ([]byte, error) { return nil, errors.New("dummy error") }
			},
			expected: `dummy error`,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				conventionalcommits.Marshal = json.Marshal
			}()
			tc.prepare()
			actual := tc.conventionalCommit.String()
			require.Equal(t, tc.expected, actual)
		})
	}
}
