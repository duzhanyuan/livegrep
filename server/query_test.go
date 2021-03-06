package server

import (
	"reflect"
	"testing"

	pb "github.com/livegrep/livegrep/src/proto/go_proto"
)

func TestParseQuery(t *testing.T) {
	cases := []struct {
		in  string
		out pb.Query
	}{
		{
			"hello",
			pb.Query{Line: "hello", FoldCase: true},
		},
		{
			"a b c",
			pb.Query{Line: "a b c", FoldCase: true},
		},
		{
			"line file:.rb",
			pb.Query{
				Line:     "line",
				File:     ".rb",
				FoldCase: true,
			},
		},
		{
			" a  ",
			pb.Query{Line: "a", FoldCase: true},
		},
		{
			"( a  )",
			pb.Query{Line: "( a  )", FoldCase: true},
		},
		{
			"Aa",
			pb.Query{Line: "Aa", FoldCase: false},
		},
		{
			"case:abc",
			pb.Query{Line: "abc", FoldCase: false},
		},
		{
			"case:abc file:^kernel/",
			pb.Query{Line: "abc", FoldCase: false, File: "^kernel/"},
		},
		{
			"case:abc file:( )",
			pb.Query{Line: "abc", FoldCase: false, File: "( )"},
		},
		{
			"a file:b c",
			pb.Query{Line: "a c", FoldCase: true, File: "b"},
		},
		{
			"a file:((abc()())()) c",
			pb.Query{Line: "a c", FoldCase: true, File: "((abc()())())"},
		},
		{
			"(  () (   ",
			pb.Query{Line: "(  () (", FoldCase: true},
		},
		{
			`a file:\(`,
			pb.Query{Line: "a", File: `\(`, FoldCase: true},
		},
		{
			`a file:(\()`,
			pb.Query{Line: "a", File: `(\()`, FoldCase: true},
		},
		{
			`(`,
			pb.Query{Line: "(", FoldCase: true},
		},
		{
			`(file:)`,
			pb.Query{Line: "(file:)", FoldCase: true},
		},
		{
			`re tags:kind:function`,
			pb.Query{Line: "re", FoldCase: true, Tags: "kind:function"},
		},
		{
			`-file:Godep re`,
			pb.Query{Line: "re", NotFile: "Godep", FoldCase: true},
		},
		{
			`-file:. -repo:Godep re`,
			pb.Query{Line: "re", NotFile: ".", NotRepo: "Godep", FoldCase: true},
		},
		{
			`-tags:kind:class re`,
			pb.Query{Line: "re", NotTags: "kind:class", FoldCase: true},
		},
		{
			`case:foo:`,
			pb.Query{Line: "foo:", FoldCase: false},
		},
		{
			`lit:.`,
			pb.Query{Line: `\.`, FoldCase: false},
		},
	}

	for _, tc := range cases {
		parsed, err := ParseQuery(tc.in)
		if !reflect.DeepEqual(tc.out, parsed) {
			t.Errorf("error parsing %q: expected %#v got %#v",
				tc.in, tc.out, parsed)
		}
		if err != nil {
			t.Errorf("parse(%v) error=%v", tc.in, err)
		}
	}

	_, err := ParseQuery(`hello case:foo`)
	if err == nil {
		t.Errorf("parse multiple regexes, no error")
	}
}
