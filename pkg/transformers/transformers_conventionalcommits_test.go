package transformers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joselitofilho/go-conventional-commits/pkg/conventionalcommits"
	"github.com/joselitofilho/go-conventional-commits/pkg/transformers"
)

func TestTransforms_ConventionalCommit(t *testing.T) {
	message := "feat: added a new feature"
	convetionalCommit := transformers.TransformConventionalCommit(message)
	require.Equal(t, "feat", convetionalCommit.Category)
}

func TestTransforms_ConventionalCommit_WithPatchChange(t *testing.T) {
	message := "fix: fixed the problem"
	convetionalCommit := transformers.TransformConventionalCommit(message)
	require.True(t, convetionalCommit.Patch)
}

func TestTransforms_ConventionalCommit_WithMinorChange(t *testing.T) {
	message := "feat: added a new feature"
	convetionalCommit := transformers.TransformConventionalCommit(message)
	require.True(t, convetionalCommit.Minor)
}

func TestTransforms_ConventionalCommit_WithMajorChange(t *testing.T) {
	message := "feat!: added a new feature"
	convetionalCommit := transformers.TransformConventionalCommit(message)
	require.True(t, convetionalCommit.Major)
}

func TestTransforms_ConventionalCommit_WithFooter(t *testing.T) {
	message := `feat: added a new feature

Refs #GCC-123
`
	convetionalCommit := transformers.TransformConventionalCommit(message)
	require.Equal(t, []string{"Refs #GCC-123"}, convetionalCommit.Footer)
}

func TestTransforms_ConventionalCommit_WithBody(t *testing.T) {
	message := `feat: added a new feature

Description of the new feature
more details

Refs #GCC-123
`

	convetionalCommit := transformers.TransformConventionalCommit(message)

	expected := `Description of the new feature
more details`
	require.Equal(t, expected, convetionalCommit.Body)
}

func TestTransforms_ConventionalCommits(t *testing.T) {
	messages := []string{`feat: added a new feature

Description of the new feature
more details

Refs #GCC-123
`,
	}

	commits := transformers.TransformConventionalCommits(messages)

	expected := conventionalcommits.ConventionalCommits{{
		Category:    "feat",
		Description: "added a new feature",
		Body:        "Description of the new feature\nmore details",
		Footer:      []string{"Refs #GCC-123"},
		Minor:       true,
	}}
	require.Equal(t, expected, commits)
}
