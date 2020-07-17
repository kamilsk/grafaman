package model

import (
	"reflect"
	"sort"
	"unsafe"

	"github.com/gobwas/glob"
)

// A Matcher checks that a metric name satisfies a Graphite query or not.
type Matcher = glob.Glob

// A Query represents a Graphite query.
type Query string

// MustCompile converts a Graphite query into a Matcher.
// If it cannot, then panic will occur.
func (query Query) MustCompile() Matcher {
	return glob.MustCompile(string(query))
}

// Metrics represent a slice of Graphite queries.
type Queries []Query

// Convert is an efficient converter from a slice of strings
// to a slice of Graphite queries.
func (queries *Queries) Convert(src []string) Queries {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	if queries == nil {
		return *(*[]Query)(unsafe.Pointer(header))
	}
	*queries = *(*[]Query)(unsafe.Pointer(header))
	return *queries
}

// MustMatchers converts a slice of Graphite queries
// to a slice of Matchers.
// If it cannot, then panic will occur.
func (queries Queries) MustMatchers() []Matcher {
	out := make([]Matcher, 0, len(queries))
	for _, query := range queries {
		out = append(out, query.MustCompile())
	}
	return out
}

// Sort orders the Graphite queries in-place by ascending.
func (queries Queries) Sort() Queries {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&queries))
	sort.Strings(*(*[]string)(unsafe.Pointer(header)))
	return queries
}
