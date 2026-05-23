package tags

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a key-value metadata label attached to a service.
type Tag struct {
	Key   string
	Value string
}

// TagSet holds a collection of tags for a named service.
type TagSet struct {
	Service string
	Tags    map[string]string
}

// NewTagSet creates an empty TagSet for the given service.
func NewTagSet(service string) *TagSet {
	return &TagSet{
		Service: service,
		Tags:    make(map[string]string),
	}
}

// Set adds or updates a tag on the TagSet.
func (ts *TagSet) Set(key, value string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("tag key must not be empty")
	}
	ts.Tags[key] = value
	return nil
}

// Get retrieves a tag value by key. Returns the value and whether it was found.
func (ts *TagSet) Get(key string) (string, bool) {
	v, ok := ts.Tags[key]
	return v, ok
}

// Delete removes a tag by key.
func (ts *TagSet) Delete(key string) {
	delete(ts.Tags, key)
}

// Keys returns a sorted list of tag keys.
func (ts *TagSet) Keys() []string {
	keys := make([]string, 0, len(ts.Tags))
	for k := range ts.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Matches reports whether the TagSet contains all tags in the given filter map.
func (ts *TagSet) Matches(filter map[string]string) bool {
	for k, v := range filter {
		if ts.Tags[k] != v {
			return false
		}
	}
	return true
}

// String returns a human-readable representation of the TagSet.
func (ts *TagSet) String() string {
	parts := make([]string, 0, len(ts.Tags))
	for _, k := range ts.Keys() {
		parts = append(parts, fmt.Sprintf("%s=%s", k, ts.Tags[k]))
	}
	return fmt.Sprintf("%s[%s]", ts.Service, strings.Join(parts, ","))
}
