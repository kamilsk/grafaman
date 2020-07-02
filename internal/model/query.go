package model

import (
	"reflect"
	"unsafe"
)

type Query string

type Queries []Query

func (queries *Queries) Convert(src []string) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	*queries = *(*[]Query)(unsafe.Pointer(header))
}

func (queries Queries) Len() int           { return len(queries) }
func (queries Queries) Less(i, j int) bool { return queries[i] < queries[j] }
func (queries Queries) Swap(i, j int)      { queries[i], queries[j] = queries[j], queries[i] }
