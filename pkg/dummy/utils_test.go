package dummy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDummyArrays(t *testing.T) {
	dummies := createDummyArrays()
	require.Len(t, dummies, 3)
}
