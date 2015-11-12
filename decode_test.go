package gocsv

import "testing"

func Test_maybeMissingStructFields(t *testing.T) {
	structTags := []fieldInfo{
		{Key: "foo"},
		{Key: "bar"},
		{Key: "baz"},
	}
	badHeaders := []string{"hi", "mom", "bacon"}
	goodHeaders := []string{"foo", "bar", "baz"}

	// no tags to match, expect no error
	if err := maybeMissingStructFields([]fieldInfo{}, goodHeaders); err != nil {
		t.Fatal(err)
	}

	// bad headers, expect an error
	if err := maybeMissingStructFields(structTags, badHeaders); err == nil {
		t.Fatal("expected an error, but no error found")
	}

	// good headers, expect no error
	if err := maybeMissingStructFields(structTags, goodHeaders); err != nil {
		t.Fatal(err)
	}

	// extra headers, but all structtags match; expect no error
	moarHeaders := append(goodHeaders, "qux", "quux", "corge", "grault")
	if err := maybeMissingStructFields(structTags, moarHeaders); err != nil {
		t.Fatal(err)
	}

	// not all structTags match, but there's plenty o' headers; expect
	// error
	mismatchedHeaders := []string{"foo", "qux", "quux", "corgi"}
	if err := maybeMissingStructFields(structTags, mismatchedHeaders); err == nil {
		t.Fatal("expected an error, but no error found")
	}
}
