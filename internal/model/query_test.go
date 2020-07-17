package model_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/model"
)

func TestQueries(t *testing.T) {
	raw := []string{"b.*", "c.*", "a.*"}

	var queries Queries
	require.NotPanics(t, func() { assert.Len(t, (*Queries)(nil).Convert(raw), len(raw)) })
	require.NotPanics(t, func() { assert.Len(t, queries.Convert(raw), len(raw)) })

	assert.False(t, sort.StringsAreSorted(raw))
	assert.Len(t, queries.Sort(), len(raw))
	assert.True(t, sort.StringsAreSorted(raw))

	assert.NotPanics(t, func() { assert.Len(t, queries.MustMatchers(), len(raw)) })
}
