package internal

import (
	"fmt"
	"reflect"
)

const (
	TagName    = "name"
	TagUsage   = "usage"
	TagDefault = "default"

	TagNameIgnored = "-"
)

func NewTag(t reflect.StructTag, prefix string) *Tag {
	return &Tag{
		tag:    t,
		prefix: prefix,
	}
}

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

func (t Tag) String() string {
	return fmt.Sprintf("tag=%s prefix=%s", t.tag, t.prefix)
}
