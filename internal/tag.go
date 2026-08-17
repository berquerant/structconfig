package internal

import (
	"fmt"
	"reflect"
)

const (
	TagName    = "name"
	TagUsage   = "usage"
	TagDefault = "default"
	TagShort   = "short"
	TagCount   = "count"
	TagSplit   = "split"
	TagSep     = "sep"

	TagNameIgnored = "-"
)

// NewTag returns a new [Tag].
// [Tag] retrieves values of all [TagName], [TagUsage], [TagDefault].
// prefix adds a prefix to the tag names.
func NewTag(t reflect.StructTag, prefix string) *Tag {
	return &Tag{
		tag:    t,
		prefix: prefix,
	}
}

// Tag is a tag of struct.
type Tag struct {
	tag    reflect.StructTag
	prefix string
}

func (t Tag) Name() (string, bool) {
	if v := t.tag.Get(t.prefix + TagName); v != "" && v != TagNameIgnored {
		return v, true
	}
	return "", false
}

func (t Tag) Usage() string {
	return t.tag.Get(t.prefix + TagUsage)
}

func (t Tag) Default() (string, bool) {
	return t.tag.Lookup(t.prefix + TagDefault)
}

func (t Tag) Short() (string, bool) {
	return t.tag.Lookup(t.prefix + TagShort)
}

func (t Tag) Count() bool {
	v, ok := t.tag.Lookup(t.prefix + TagCount)
	if !ok {
		return false
	}
	return v == "true" || v == ""
}

func (t Tag) Split() bool {
	v, ok := t.tag.Lookup(t.prefix + TagSplit)
	if !ok {
		return false
	}
	return v == "true" || v == ""
}

func (t Tag) Sep() string {
	if v, ok := t.tag.Lookup(t.prefix + TagSep); ok {
		return v
	}
	return ","
}

func (t Tag) String() string {
	return fmt.Sprintf("tag=%s prefix=%s", t.tag, t.prefix)
}
