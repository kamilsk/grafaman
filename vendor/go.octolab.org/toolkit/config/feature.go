package config

import (
	"fmt"
	"strings"
)

// Feature describe a feature.
type Feature struct {
	Name    string
	Enabled bool
}

// String returns a string representation of the feature.
func (feature Feature) String() string {
	return fmt.Sprintf("%s=%v", feature.Name, feature.Enabled)
}

// Features defines a list of features.
type Features []Feature

// String returns a string representation of the feature list.
func (features Features) String() string {
	if len(features) == 0 {
		return "-"
	}
	list := make([]string, 0, len(features))
	for _, feature := range features {
		list = append(list, feature.String())
	}
	return strings.Join(list, ", ")
}
