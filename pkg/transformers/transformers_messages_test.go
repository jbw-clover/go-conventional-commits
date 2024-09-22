package transformers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsuyoshiwada/go-gitlog"

	"github.com/joselitofilho/go-conventional-commits/pkg/transformers"
)

func TestTransforms_Messages(t *testing.T) {
	commits := []*gitlog.Commit{{
		Hash: &gitlog.Hash{
			Long:  "f4f7deca6d08fd34919211d00daac1763fd20cbb",
			Short: "f4f7dec",
		},
		Subject: "feat: added a new feature",
		Body:    "Refs #GCC-123",
	}}
	messages := transformers.TransformMessages(commits, "https://url/commits/")
	require.Equal(t, []string{"feat: added a new feature #https://url/commits/f4f7dec\n\nRefs #GCC-123"}, messages)
}
